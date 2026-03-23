package rules

import (
	"regexp"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

var pushBranchSuggestionRe = regexp.MustCompile(`(?m)^ +(git push [^\s]+ [^\s]+)`)

func init() {
	register(Rule{
		Name: "git_push_different_branch_names",
		Match: func(cmd types.Command) bool {
			if !strings.HasPrefix(cmd.Script, "git ") {
				return false
			}
			return strings.Contains(cmd.Script, "push") &&
				strings.Contains(cmd.Output, "The upstream branch of your current branch does not match")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			m := pushBranchSuggestionRe.FindStringSubmatch(cmd.Output)
			if len(m) < 2 {
				return nil
			}
			return single(m[1])
		},
	})
}
