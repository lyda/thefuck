package rules

import (
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "ag_literal",
		Match: func(cmd types.Command) bool {
			return cmd.ScriptParts()[0] == "ag" &&
				strings.HasSuffix(cmd.Output, "run ag with -Q\n")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			return single(strings.Replace(cmd.Script, "ag", "ag -Q", 1))
		},
	})
}
