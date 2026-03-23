package rules

import (
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "git_remote_delete",
		Match: func(cmd types.Command) bool {
			return strings.HasPrefix(cmd.Script, "git ") &&
				strings.Contains(cmd.Script, "remote delete")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			return single(strings.Replace(cmd.Script, "delete", "remove", 1))
		},
	})
}
