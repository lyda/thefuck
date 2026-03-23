package rules

import (
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "git_clone_git_clone",
		Match: func(cmd types.Command) bool {
			return strings.HasPrefix(cmd.Script, "git ") &&
				strings.Contains(cmd.Script, " git clone ") &&
				strings.Contains(cmd.Output, "fatal: Too many arguments.")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			return single(strings.Replace(cmd.Script, " git clone ", " ", 1))
		},
	})
}
