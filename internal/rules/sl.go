package rules

import "github.com/lyda/thefuck/internal/types"

func init() {
	register(Rule{
		Name:     "sl",
		Priority: 100,
		Match: func(cmd types.Command) bool {
			return cmd.Script == "sl"
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			return single("ls")
		},
	})
}
