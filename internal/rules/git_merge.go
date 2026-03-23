package rules

import (
	"regexp"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

var mergeBranchRe = regexp.MustCompile(`merge: (.+) - not something we can merge`)
var mergeRemoteBranchRe = regexp.MustCompile(`Did you mean this\?\n\t([^\n]+)`)

func init() {
	register(Rule{
		Name: "git_merge",
		Match: func(cmd types.Command) bool {
			if !strings.HasPrefix(cmd.Script, "git ") {
				return false
			}
			return strings.Contains(cmd.Script, "merge") &&
				strings.Contains(cmd.Output, " - not something we can merge") &&
				strings.Contains(cmd.Output, "Did you mean this?")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			bm := mergeBranchRe.FindStringSubmatch(cmd.Output)
			if len(bm) < 2 {
				return nil
			}
			rm := mergeRemoteBranchRe.FindStringSubmatch(cmd.Output)
			if len(rm) < 2 {
				return nil
			}
			return single(replaceArgument(cmd.Script, bm[1], rm[1]))
		},
	})
}
