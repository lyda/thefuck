package rules

import (
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "git_branch_list",
		Match: func(cmd types.Command) bool {
			if !strings.HasPrefix(cmd.Script, "git ") {
				return false
			}
			parts := cmd.ScriptParts()
			if len(parts) < 3 {
				return false
			}
			return parts[1] == "branch" && parts[2] == "list"
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			return single(shellAnd("git branch --delete list", "git branch"))
		},
	})
}
