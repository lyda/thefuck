package rules

import (
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "pip_install",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			if len(parts) == 0 {
				return false
			}
			prog := parts[0]
			if prog != "pip" && prog != "pip2" && prog != "pip3" {
				return false
			}
			return strings.Contains(cmd.Script, "pip install") &&
				strings.Contains(cmd.Output, "Permission denied")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			if !strings.Contains(cmd.Script, "--user") {
				// Attempt 1: add --user
				return single(strings.Replace(cmd.Script, " install ", " install --user ", 1))
			}
			// Attempt 2: sudo without --user
			return single("sudo " + strings.Replace(cmd.Script, " --user", "", 1))
		},
	})
}
