package rules

import (
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "unsudo",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			if len(parts) == 0 || parts[0] != "sudo" {
				return false
			}
			return strings.Contains(strings.ToLower(cmd.Output), "you cannot perform this operation as root")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			parts := cmd.ScriptParts()
			if len(parts) <= 1 {
				return nil
			}
			return single(strings.Join(parts[1:], " "))
		},
	})
}
