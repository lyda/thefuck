package rules

import (
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "terraform_init",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			if len(parts) == 0 || parts[0] != "terraform" {
				return false
			}
			lower := strings.ToLower(cmd.Output)
			return strings.Contains(lower, "this module is not yet installed") ||
				strings.Contains(lower, "initialization required")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			return single(shellAnd("terraform init", cmd.Script))
		},
	})
}
