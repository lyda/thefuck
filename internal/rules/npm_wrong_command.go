package rules

import (
	"os/exec"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func npmGetWrongCommand(parts []string) string {
	for _, p := range parts[1:] {
		if !strings.HasPrefix(p, "-") {
			return p
		}
	}
	return ""
}

func npmGetAvailableCommands(output string) []string {
	var commands []string
	listing := false
	for _, line := range strings.Split(output, "\n") {
		if strings.HasPrefix(line, "where <command> is one of:") {
			listing = true
			continue
		}
		if listing {
			if strings.TrimSpace(line) == "" {
				break
			}
			for _, cmd := range strings.Split(line, ",") {
				if s := strings.TrimSpace(cmd); s != "" {
					commands = append(commands, s)
				}
			}
		}
	}
	return commands
}

func init() {
	register(Rule{
		Name: "npm_wrong_command",
		Match: func(cmd types.Command) bool {
			if _, err := exec.LookPath("npm"); err != nil {
				return false
			}
			parts := cmd.ScriptParts()
			if len(parts) == 0 || parts[0] != "npm" {
				return false
			}
			return strings.Contains(cmd.Output, "where <command> is one of:") &&
				npmGetWrongCommand(parts) != ""
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			parts := cmd.ScriptParts()
			if len(parts) < 2 {
				return nil
			}
			wrongCmd := npmGetWrongCommand(parts)
			if wrongCmd == "" {
				return nil
			}
			available := npmGetAvailableCommands(cmd.Output)
			closest := getCloseMatches(wrongCmd, available, 0.6)
			if len(closest) == 0 {
				return nil
			}
			return single(replaceArgument(cmd.Script, wrongCmd, closest[0]))
		},
	})
}
