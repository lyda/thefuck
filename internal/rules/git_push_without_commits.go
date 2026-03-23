package rules

import (
	"regexp"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

var gitPushWithoutCommitsRe = regexp.MustCompile(`src refspec \w+ does not match any`)

func init() {
	register(Rule{
		Name: "git_push_without_commits",
		Match: func(cmd types.Command) bool {
			return strings.HasPrefix(cmd.Script, "git ") &&
				gitPushWithoutCommitsRe.MatchString(cmd.Output)
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			return single(shellAnd(`git commit -m "Initial commit"`, cmd.Script))
		},
	})
}
