package rules

import (
	"fmt"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "git_help_aliased",
		Match: func(cmd types.Command) bool {
			if !strings.HasPrefix(cmd.Script, "git ") {
				return false
			}
			return strings.Contains(cmd.Script, "help") &&
				strings.Contains(cmd.Output, " is aliased to ")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			// output format: "... `alias' is aliased to `real_cmd ...'"
			// split on backtick: [before, alias, " is aliased to ", real_cmd_part, ...]
			parts := strings.SplitN(cmd.Output, "`", 3)
			if len(parts) < 3 {
				return nil
			}
			remainder := parts[2]
			// remainder starts with "real_cmd ..." followed by "'"
			aliased := strings.SplitN(remainder, "'", 2)[0]
			aliased = strings.SplitN(aliased, " ", 2)[0]
			return single(fmt.Sprintf("git help %s", aliased))
		},
	})
}
