package rules

import (
	"slices"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "git_commit_reset",
		Match: func(cmd types.Command) bool {
			if !hasGitPrefix(cmd) {
				return false
			}
			return slices.Contains(cmd.ScriptParts(), "commit")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			return single("git reset HEAD~")
		},
	})
}
