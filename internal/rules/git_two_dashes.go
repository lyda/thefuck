package rules

import (
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "git_two_dashes",
		Match: func(cmd types.Command) bool {
			return strings.HasPrefix(cmd.Script, "git ") &&
				strings.Contains(cmd.Output, "error: did you mean `") &&
				strings.Contains(cmd.Output, "` (with two dashes ?)")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			// Output looks like: error: did you mean `--foo` (with two dashes ?)
			parts := strings.Split(cmd.Output, "`")
			if len(parts) < 2 {
				return nil
			}
			to := parts[1] // e.g. "--foo"
			from := to[1:] // e.g. "-foo"
			return single(replaceArgument(cmd.Script, from, to))
		},
	})
}
