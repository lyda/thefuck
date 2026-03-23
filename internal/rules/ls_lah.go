package rules

import (
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "ls_lah",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			return len(parts) > 0 && parts[0] == "ls" &&
				!strings.Contains(cmd.Script, "ls -")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			parts := cmd.ScriptParts()
			newParts := make([]string, len(parts))
			copy(newParts, parts)
			newParts[0] = "ls -lah"
			return single(strings.Join(newParts, " "))
		},
	})
}
