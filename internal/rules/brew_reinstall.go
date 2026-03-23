package rules

import (
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "brew_reinstall",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			if len(parts) < 1 || parts[0] != "brew" {
				return false
			}
			return strings.Contains(cmd.Script, "install") &&
				strings.Contains(cmd.Output, "already installed") &&
				strings.Contains(cmd.Output, "brew reinstall")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			return single(strings.Replace(cmd.Script, "install", "reinstall", 1))
		},
	})
}
