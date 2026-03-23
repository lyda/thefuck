package rules

import (
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "go_run",
		Match: func(cmd types.Command) bool {
			return cmd.ScriptParts()[0] == "go" &&
				strings.HasPrefix(cmd.Script, "go run ") &&
				!strings.HasSuffix(cmd.Script, ".go")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			return single(cmd.Script + ".go")
		},
	})
}
