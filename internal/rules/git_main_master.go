package rules

import (
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name:     "git_main_master",
		Priority: 1200,
		Match: func(cmd types.Command) bool {
			if !strings.HasPrefix(cmd.Script, "git ") {
				return false
			}
			return strings.Contains(cmd.Output, "'master'") ||
				strings.Contains(cmd.Output, "'main'")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			if strings.Contains(cmd.Output, "'master'") {
				return single(strings.ReplaceAll(cmd.Script, "master", "main"))
			}
			return single(strings.ReplaceAll(cmd.Script, "main", "master"))
		},
	})
}
