package rules

import (
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "git_push_force",
		Match: func(cmd types.Command) bool {
			if !strings.HasPrefix(cmd.Script, "git ") {
				return false
			}
			return strings.Contains(cmd.Script, "push") &&
				strings.Contains(cmd.Output, "! [rejected]") &&
				strings.Contains(cmd.Output, "failed to push some refs to") &&
				strings.Contains(cmd.Output, "Updates were rejected because the tip of your current branch is behind")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			return single(replaceArgument(cmd.Script, "push", "push --force-with-lease"))
		},
	})
}
