package rules

import (
	"slices"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "git_tag_force",
		Match: func(cmd types.Command) bool {
			if !strings.HasPrefix(cmd.Script, "git ") {
				return false
			}
			parts := cmd.ScriptParts()
			hasTag := slices.Contains(parts, "tag")
			return hasTag && strings.Contains(cmd.Output, "already exists")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			return single(replaceArgument(cmd.Script, "tag", "tag --force"))
		},
	})
}
