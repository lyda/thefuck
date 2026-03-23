package rules

import (
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "git_stash_pop",
		Match: func(cmd types.Command) bool {
			return strings.HasPrefix(cmd.Script, "git ") &&
				strings.Contains(cmd.Script, "stash") &&
				strings.Contains(cmd.Script, "pop") &&
				strings.Contains(cmd.Output, "Your local changes to the following files would be overwritten by merge")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			return single(shellAnd("git add --update", "git stash pop", "git reset ."))
		},
	})
}
