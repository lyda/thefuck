package rules

import (
	"os/exec"
	"strings"
	"sync"

	"github.com/lyda/thefuck/internal/types"
)

var (
	aptOpsOnce sync.Once
	// aptOpsCache maps binary name ("apt", "apt-get", "apt-cache") to its
	// list of valid operations, populated lazily on first use.
	aptOpsCache = map[string][]string{}
)

// getAptOperations returns the valid operations for the given apt binary by
// parsing its --help output.  Results are cached for the lifetime of the
// process.
func getAptOperations(app string) []string {
	// We use a single Once to avoid spawning multiple subprocesses in
	// parallel; the cache is populated for all three binaries on first call
	// to any of them, since we only need one subprocess per binary.
	// Actually, call per-app lazily with a simple mutex-protected map.
	aptOpsOnce.Do(func() {
		for _, binary := range []string{"apt", "apt-get", "apt-cache"} {
			if _, err := exec.LookPath(binary); err != nil {
				continue
			}
			out, err := exec.Command(binary, "--help").CombinedOutput() // #nosec G204
			if err != nil && len(out) == 0 {
				continue
			}
			aptOpsCache[binary] = parseAptOperations(binary, string(out))
		}
	})
	return aptOpsCache[app]
}

func parseAptOperations(app, helpText string) []string {
	lines := strings.Split(helpText, "\n")
	var ops []string
	collecting := false

	// "apt" uses "Basic commands:" or "Most used commands:"
	// "apt-get" / "apt-cache" use "Commands:" or "Most used commands:"
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if collecting {
			if trimmed == "" {
				if app != "apt" {
					// apt-get / apt-cache: stop at first blank line
					break
				}
				continue
			}
			fields := strings.Fields(trimmed)
			if len(fields) > 0 {
				ops = append(ops, fields[0])
			}
		} else {
			switch {
			case strings.HasPrefix(trimmed, "Basic commands:"),
				strings.HasPrefix(trimmed, "Most used commands:"),
				strings.HasPrefix(trimmed, "Commands:"):
				collecting = true
			}
		}
	}
	return ops
}

func init() {
	register(Rule{
		Name: "apt_invalid_operation",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			if len(parts) == 0 {
				return false
			}
			app := parts[0]
			// Strip leading "sudo" if present.
			if app == "sudo" && len(parts) > 1 {
				app = parts[1]
			}
			if app != "apt" && app != "apt-get" && app != "apt-cache" {
				return false
			}
			return strings.Contains(cmd.Output, "E: Invalid operation")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			// The invalid operation is the last word of the error line.
			var invalidOp string
			for _, line := range strings.Split(cmd.Output, "\n") {
				if strings.Contains(line, "E: Invalid operation") {
					fields := strings.Fields(line)
					if len(fields) > 0 {
						invalidOp = fields[len(fields)-1]
					}
					break
				}
			}
			if invalidOp == "" {
				return nil
			}

			// Common alias: uninstall -> remove
			if invalidOp == "uninstall" {
				return single(strings.Replace(cmd.Script, "uninstall", "remove", 1))
			}

			parts := cmd.ScriptParts()
			app := parts[0]
			if app == "sudo" && len(parts) > 1 {
				app = parts[1]
			}

			ops := getAptOperations(app)
			closest := getCloseMatches(invalidOp, ops, 0.6)
			var result []string
			for _, c := range closest {
				result = append(result, strings.Replace(cmd.Script, invalidOp, c, 1))
			}
			if len(result) == 0 {
				return nil
			}
			return multi(result)
		},
	})
}
