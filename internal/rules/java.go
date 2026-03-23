package rules

import (
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "java",
		Match: func(cmd types.Command) bool {
			return cmd.ScriptParts()[0] == "java" &&
				strings.HasSuffix(cmd.Script, ".java")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			return single(cmd.Script[:len(cmd.Script)-5])
		},
	})
}
