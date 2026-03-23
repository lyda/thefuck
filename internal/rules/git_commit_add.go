package rules

import (
	"slices"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "git_commit_add",
		Match: func(cmd types.Command) bool {
			if !strings.HasPrefix(cmd.Script, "git ") {
				return false
			}
			parts := cmd.ScriptParts()
			hasCommit := slices.Contains(parts, "commit")
			return hasCommit && strings.Contains(cmd.Output, "no changes added to commit")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			scripts := []string{
				replaceArgument(cmd.Script, "commit", "commit -a"),
				replaceArgument(cmd.Script, "commit", "commit -p"),
			}
			return multi(scripts)
		},
	})
}
