package rules

import (
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "python_execute",
		Match: func(cmd types.Command) bool {
			return cmd.ScriptParts()[0] == "python" &&
				!strings.HasSuffix(cmd.Script, ".py")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			return single(cmd.Script + ".py")
		},
	})
}
