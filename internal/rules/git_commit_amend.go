package rules

import (
	"slices"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "git_commit_amend",
		Match: func(cmd types.Command) bool {
			if !hasGitPrefix(cmd) {
				return false
			}
			return slices.Contains(cmd.ScriptParts(), "commit")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			return single("git commit --amend")
		},
	})
}

// hasGitPrefix is a small helper used by several git rules.
func hasGitPrefix(cmd types.Command) bool {
	parts := cmd.ScriptParts()
	return len(parts) > 0 && parts[0] == "git"
}
