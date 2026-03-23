package rules

import (
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "python_command",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			if len(parts) == 0 || !strings.HasSuffix(parts[0], ".py") {
				return false
			}
			lower := cmd.Output
			return strings.Contains(lower, "Permission denied") ||
				strings.Contains(lower, "command not found")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			return single("python " + cmd.Script)
		},
	})
}
