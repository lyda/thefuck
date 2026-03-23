package rules

import (
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "fix_alt_space",
		Match: func(cmd types.Command) bool {
			return strings.ContainsRune(cmd.Script, '\u00a0') &&
				strings.Contains(strings.ToLower(cmd.Output), "command not found")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			return single(strings.ReplaceAll(cmd.Script, "\u00a0", " "))
		},
	})
}
