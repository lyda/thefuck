package rules

import (
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name:     "dry",
		Priority: 900,
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			return len(parts) >= 2 && parts[0] == parts[1]
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			parts := cmd.ScriptParts()
			return single(strings.Join(parts[1:], " "))
		},
	})
}
