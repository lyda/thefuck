package rules

import (
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "git_rebase_no_changes",
		Match: func(cmd types.Command) bool {
			if !strings.HasPrefix(cmd.Script, "git ") {
				return false
			}
			hasRebase := false
			hasContinue := false
			for _, p := range cmd.ScriptParts() {
				if p == "rebase" {
					hasRebase = true
				}
				if p == "--continue" {
					hasContinue = true
				}
			}
			return hasRebase && hasContinue &&
				strings.Contains(cmd.Output, "No changes - did you forget to use 'git add'?")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			return single("git rebase --skip")
		},
	})
}
