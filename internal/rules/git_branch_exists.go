package rules

import (
	"regexp"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

var gitBranchExistsRe = regexp.MustCompile(`fatal: A branch named '([^']*)' already exists`)

func init() {
	register(Rule{
		Name: "git_branch_exists",
		Match: func(cmd types.Command) bool {
			return strings.HasPrefix(cmd.Script, "git ") &&
				strings.Contains(cmd.Output, "fatal: A branch named '") &&
				strings.Contains(cmd.Output, "' already exists")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			m := gitBranchExistsRe.FindStringSubmatch(cmd.Output)
			if len(m) < 2 {
				return nil
			}
			branch := m[1]
			return multi([]string{
				shellAnd("git branch -d "+branch, cmd.Script),
				shellAnd("git branch -D "+branch, cmd.Script),
			})
		},
	})
}
