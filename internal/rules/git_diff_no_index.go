package rules

import (
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "git_diff_no_index",
		Match: func(cmd types.Command) bool {
			if !strings.HasPrefix(cmd.Script, "git ") {
				return false
			}
			if !strings.Contains(cmd.Script, "diff") {
				return false
			}
			if strings.Contains(cmd.Script, "--no-index") {
				return false
			}
			// Count non-flag args after "diff"
			parts := cmd.ScriptParts()
			files := 0
			inDiff := false
			for _, p := range parts {
				if p == "diff" {
					inDiff = true
					continue
				}
				if inDiff && !strings.HasPrefix(p, "-") {
					files++
				}
			}
			return files == 2
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			return single(replaceArgument(cmd.Script, "diff", "diff --no-index"))
		},
	})
}
