package rules

import (
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "man_no_space",
		Match: func(cmd types.Command) bool {
			return strings.HasPrefix(cmd.Script, "man") &&
				len(cmd.Script) > 3 &&
				cmd.Script[3] != ' ' &&
				strings.Contains(strings.ToLower(cmd.Output), "command not found")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			return single("man " + cmd.Script[3:])
		},
	})
}
