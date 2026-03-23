package rules

import (
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "git_pull_clone",
		Match: func(cmd types.Command) bool {
			if !strings.HasPrefix(cmd.Script, "git ") {
				return false
			}
			return strings.Contains(cmd.Output, "fatal: Not a git repository") &&
				strings.Contains(cmd.Output, "Stopping at filesystem boundary (GIT_DISCOVERY_ACROSS_FILESYSTEM not set).")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			return single(replaceArgument(cmd.Script, "pull", "clone"))
		},
	})
}
