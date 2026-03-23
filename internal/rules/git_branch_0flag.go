package rules

import (
	"fmt"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

// first0Flag returns the first script part that is exactly 2 characters
// long and starts with "0" (e.g. "0d" instead of "-d").
func first0Flag(parts []string) string {
	for _, p := range parts {
		if len(p) == 2 && strings.HasPrefix(p, "0") {
			return p
		}
	}
	return ""
}

func init() {
	register(Rule{
		Name: "git_branch_0flag",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			if !strings.HasPrefix(cmd.Script, "git ") {
				return false
			}
			if len(parts) < 2 || parts[1] != "branch" {
				return false
			}
			return first0Flag(parts) != ""
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			parts := cmd.ScriptParts()
			branchName := first0Flag(parts)
			if branchName == "" {
				return nil
			}
			fixedFlag := strings.Replace(branchName, "0", "-", 1)
			fixedScript := strings.Replace(cmd.Script, branchName, fixedFlag, 1)
			if strings.Contains(cmd.Output, "A branch named '") &&
				strings.Contains(cmd.Output, "' already exists.") {
				deleteCmd := fmt.Sprintf("git branch -D %s", branchName)
				return single(shellAnd(deleteCmd, fixedScript))
			}
			return single(fixedScript)
		},
	})
}
