package rules

import (
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "git_pull",
		Match: func(cmd types.Command) bool {
			return strings.HasPrefix(cmd.Script, "git ") &&
				strings.Contains(cmd.Script, "pull") &&
				strings.Contains(cmd.Output, "set-upstream")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			lines := strings.Split(cmd.Output, "\n")
			// Use the third-to-last non-empty line that contains the suggested command
			if len(lines) < 3 {
				return nil
			}
			line := strings.TrimSpace(lines[len(lines)-3])
			parts := strings.Split(line, " ")
			branch := parts[len(parts)-1]
			setUpstream := strings.NewReplacer("<remote>", "origin", "<branch>", branch).Replace(line)
			return single(shellAnd(setUpstream, cmd.Script))
		},
	})
}
