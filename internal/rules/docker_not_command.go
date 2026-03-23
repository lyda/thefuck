package rules

import (
	"os/exec"
	"regexp"
	"strings"
	"sync"

	"github.com/lyda/thefuck/internal/types"
)

var (
	dockerCommandsOnce sync.Once
	dockerCommands     []string
)

var dockerNotCmdRe = regexp.MustCompile(`docker: '(\w+)' is not a docker command\.`)

// parseDockerSection parses lines from docker help output starting after a
// line that begins with startsWith, collecting the first word of each
// non-empty line until an empty line is encountered.
func parseDockerSection(lines []string, startsWith string) []string {
	var result []string
	found := false
	for _, line := range lines {
		if !found {
			if strings.HasPrefix(line, startsWith) {
				found = true
			}
			continue
		}
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			break
		}
		result = append(result, strings.Fields(trimmed)[0])
	}
	return result
}

// getDockerCommands returns all valid docker commands by parsing `docker help`.
func getDockerCommands() []string {
	dockerCommandsOnce.Do(func() {
		out, _ := exec.Command("docker").CombinedOutput() // #nosec G204
		lines := strings.Split(string(out), "\n")

		var mgmt []string
		for _, line := range lines {
			if strings.HasPrefix(line, "Management Commands:") {
				mgmt = parseDockerSection(lines, "Management Commands:")
				break
			}
		}
		regular := parseDockerSection(lines, "Commands:")
		dockerCommands = append(mgmt, regular...)
	})
	return dockerCommands
}

func init() {
	register(Rule{
		Name: "docker_not_command",
		Match: func(cmd types.Command) bool {
			if _, err := exec.LookPath("docker"); err != nil {
				return false
			}
			return strings.Contains(cmd.Output, "is not a docker command") ||
				strings.Contains(cmd.Output, "Usage:\tdocker")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			// Handle sub-command usage error (e.g. `docker run` with wrong sub-sub-command)
			if strings.Contains(cmd.Output, "Usage:") {
				parts := cmd.ScriptParts()
				if len(parts) > 2 {
					subCmds := parseDockerSection(strings.Split(cmd.Output, "\n"), "Commands:")
					closest := getCloseMatches(parts[2], subCmds, 0.6)
					scripts := make([]string, 0, len(closest))
					for _, c := range closest {
						scripts = append(scripts, replaceArgument(cmd.Script, parts[2], c))
					}
					if len(scripts) > 0 {
						return multi(scripts)
					}
					return nil
				}
			}

			m := dockerNotCmdRe.FindStringSubmatch(cmd.Output)
			if len(m) < 2 {
				return nil
			}
			wrongCmd := m[1]
			closest := getCloseMatches(wrongCmd, getDockerCommands(), 0.6)
			scripts := make([]string, 0, len(closest))
			for _, c := range closest {
				scripts = append(scripts, replaceArgument(cmd.Script, wrongCmd, c))
			}
			if len(scripts) == 0 {
				return nil
			}
			return multi(scripts)
		},
	})
}
