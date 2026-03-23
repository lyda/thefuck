package rules

import (
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "brew_update_formula",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			if len(parts) < 2 || parts[0] != "brew" {
				return false
			}
			return strings.Contains(cmd.Script, "update") &&
				strings.Contains(cmd.Output, "Error: This command updates brew itself") &&
				strings.Contains(cmd.Output, "Use `brew upgrade")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			return single(strings.Replace(cmd.Script, "update", "upgrade", 1))
		},
	})
}
