package rules

import (
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "quotation_marks",
		Match: func(cmd types.Command) bool {
			return strings.Contains(cmd.Script, "'") &&
				strings.Contains(cmd.Script, "\"")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			return single(strings.ReplaceAll(cmd.Script, "'", "\""))
		},
	})
}
