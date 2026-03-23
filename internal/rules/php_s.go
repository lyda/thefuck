package rules

import (
	"slices"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "php_s",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			if len(parts) < 2 || parts[0] != "php" {
				return false
			}
			hasS := slices.Contains(parts, "-s")
			return hasS && parts[len(parts)-1] != "-s"
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			return single(replaceArgument(cmd.Script, "-s", "-S"))
		},
	})
}
