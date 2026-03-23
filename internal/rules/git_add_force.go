package rules

import (
	"slices"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "git_add_force",
		Match: func(cmd types.Command) bool {
			if !strings.HasPrefix(cmd.Script, "git ") {
				return false
			}
			parts := cmd.ScriptParts()
			hasAdd := slices.Contains(parts, "add")
			return hasAdd && strings.Contains(cmd.Output, "Use -f if you really want to add them.")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			return single(replaceArgument(cmd.Script, "add", "add --force"))
		},
	})
}
