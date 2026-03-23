package rules

import (
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "git_diff_staged",
		Match: func(cmd types.Command) bool {
			if !strings.HasPrefix(cmd.Script, "git ") {
				return false
			}
			return strings.Contains(cmd.Script, "diff") &&
				!strings.Contains(cmd.Script, "--staged")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			return single(replaceArgument(cmd.Script, "diff", "diff --staged"))
		},
	})
}
