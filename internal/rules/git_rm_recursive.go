package rules

import (
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "git_rm_recursive",
		Match: func(cmd types.Command) bool {
			return strings.HasPrefix(cmd.Script, "git ") &&
				strings.Contains(cmd.Script, " rm ") &&
				strings.Contains(cmd.Output, "fatal: not removing '") &&
				strings.Contains(cmd.Output, "' recursively without -r")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			parts := cmd.ScriptParts()
			idx := -1
			for i, p := range parts {
				if p == "rm" {
					idx = i
					break
				}
			}
			if idx < 0 {
				return nil
			}
			newParts := make([]string, len(parts)+1)
			copy(newParts, parts[:idx+1])
			newParts[idx+1] = "-r"
			copy(newParts[idx+2:], parts[idx+1:])
			return single(strings.Join(newParts, " "))
		},
	})
}
