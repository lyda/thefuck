package rules

import (
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "remove_trailing_cedilla",
		Match: func(cmd types.Command) bool {
			return strings.HasSuffix(cmd.Script, "ç")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			// "ç" is a multi-byte UTF-8 rune; trim it by rune boundary.
			r := []rune(cmd.Script)
			return single(string(r[:len(r)-1]))
		},
	})
}
