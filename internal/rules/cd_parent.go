package rules

import "github.com/lyda/thefuck/internal/types"

func init() {
	register(Rule{
		Name:     "cd_parent",
		Priority: 100,
		Match: func(cmd types.Command) bool {
			return cmd.Script == "cd.."
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			return single("cd ..")
		},
	})
}
