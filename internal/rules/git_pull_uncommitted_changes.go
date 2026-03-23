package rules

import (
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "git_pull_uncommitted_changes",
		Match: func(cmd types.Command) bool {
			if !strings.HasPrefix(cmd.Script, "git ") {
				return false
			}
			return strings.Contains(cmd.Script, "pull") &&
				(strings.Contains(cmd.Output, "You have unstaged changes") ||
					strings.Contains(cmd.Output, "contains uncommitted changes"))
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			return single(shellAnd("git stash", "git pull", "git stash pop"))
		},
	})
}
