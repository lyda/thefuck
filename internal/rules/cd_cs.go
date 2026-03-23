package rules

import (
	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name:     "cd_cs",
		Priority: 900,
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			return len(parts) > 0 && parts[0] == "cs"
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			return single("cd" + cmd.Script[2:])
		},
	})
}
