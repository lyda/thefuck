package rules

import (
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "grep_recursive",
		Match: func(cmd types.Command) bool {
			return strings.HasPrefix(cmd.Script, "grep") &&
				strings.Contains(strings.ToLower(cmd.Output), "is a directory")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			return single("grep -r " + cmd.Script[5:])
		},
	})
}
