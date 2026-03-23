package rules

import (
	"os/exec"
	"strings"
	"sync"

	"github.com/lyda/thefuck/internal/types"
)

var (
	goCommandsOnce sync.Once
	goCommands     []string
)

// getGoCommands returns the list of valid `go` subcommands by parsing
// the usage output of `go` (which prints to stderr).
func getGoCommands() []string {
	goCommandsOnce.Do(func() {
		out, _ := exec.Command("go").CombinedOutput() // #nosec G204
		lines := strings.Split(string(out), "\n")

		// Find "The commands are:" then skip one blank line, collect until blank.
		collecting := false
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if trimmed == "The commands are:" {
				collecting = true
				continue
			}
			if !collecting {
				continue
			}
			if trimmed == "" {
				if len(goCommands) > 0 {
					break
				}
				continue
			}
			// Each line is "  subcommand    description"
			goCommands = append(goCommands, strings.Fields(trimmed)[0])
		}
	})
	return goCommands
}

func init() {
	register(Rule{
		Name: "go_unknown_command",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			return len(parts) >= 2 && parts[0] == "go" &&
				strings.Contains(cmd.Output, "unknown command")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			parts := cmd.ScriptParts()
			if len(parts) < 2 {
				return nil
			}
			closest := getCloseMatches(parts[1], getGoCommands(), 0.6)
			if len(closest) == 0 {
				return nil
			}
			return single(replaceArgument(cmd.Script, parts[1], closest[0]))
		},
	})
}
