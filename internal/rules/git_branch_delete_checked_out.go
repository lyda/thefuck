package rules

import (
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "git_branch_delete_checked_out",
		Match: func(cmd types.Command) bool {
			return strings.HasPrefix(cmd.Script, "git ") &&
				(strings.Contains(cmd.Script, "branch -d") || strings.Contains(cmd.Script, "branch -D")) &&
				strings.Contains(cmd.Output, "error: Cannot delete branch '") &&
				strings.Contains(cmd.Output, "' checked out at '")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			fixed := replaceArgument(cmd.Script, "-d", "-D")
			return single(shellAnd("git checkout master", fixed))
		},
	})
}
