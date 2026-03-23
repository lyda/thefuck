package rules

import (
	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "cargo",
		Match: func(cmd types.Command) bool {
			return cmd.Script == "cargo"
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			return single("cargo build")
		},
	})
}
