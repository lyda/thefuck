package rules

import (
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "git_push_pull",
		Match: func(cmd types.Command) bool {
			return strings.HasPrefix(cmd.Script, "git ") &&
				strings.Contains(cmd.Script, "push") &&
				strings.Contains(cmd.Output, "! [rejected]") &&
				strings.Contains(cmd.Output, "failed to push some refs to") &&
				(strings.Contains(cmd.Output, "Updates were rejected because the tip of your current branch is behind") ||
					strings.Contains(cmd.Output, "Updates were rejected because the remote contains work that you do"))
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			pullCmd := replaceArgument(cmd.Script, "push", "pull")
			return single(shellAnd(pullCmd, cmd.Script))
		},
	})
}
