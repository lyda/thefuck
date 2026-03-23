package rules

import (
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "git_remote_seturl_add",
		Match: func(cmd types.Command) bool {
			if !strings.HasPrefix(cmd.Script, "git ") {
				return false
			}
			return strings.Contains(cmd.Script, "set-url") &&
				strings.Contains(cmd.Output, "fatal: No such remote")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			return single(replaceArgument(cmd.Script, "set-url", "add"))
		},
	})
}
