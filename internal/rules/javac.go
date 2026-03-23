package rules

import (
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "javac",
		Match: func(cmd types.Command) bool {
			return cmd.ScriptParts()[0] == "javac" &&
				!strings.HasSuffix(cmd.Script, ".java")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			return single(cmd.Script + ".java")
		},
	})
}
