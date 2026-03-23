package rules

import (
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "systemctl",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			if len(parts) == 0 || parts[0] != "systemctl" {
				return false
			}
			// Catches "Unknown operation 'service'." when systemctl args are misordered
			// The Python code checks: len(cmd) - cmd.index('systemctl') == 3
			// meaning exactly 2 args after 'systemctl' (total 3 parts from systemctl onward)
			idx := -1
			for i, p := range parts {
				if p == "systemctl" {
					idx = i
					break
				}
			}
			return strings.Contains(cmd.Output, "Unknown operation '") &&
				len(parts)-idx == 3
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			parts := cmd.ScriptParts()
			newParts := make([]string, len(parts))
			copy(newParts, parts)
			last := len(newParts) - 1
			newParts[last], newParts[last-1] = newParts[last-1], newParts[last]
			return single(strings.Join(newParts, " "))
		},
	})
}
