package rules

import (
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "brew_uninstall",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			if len(parts) < 2 || parts[0] != "brew" {
				return false
			}
			return (parts[1] == "uninstall" || parts[1] == "rm" || parts[1] == "remove") &&
				strings.Contains(cmd.Output, "brew uninstall --force")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			parts := cmd.ScriptParts()
			// Build: brew uninstall --force <rest>
			result := make([]string, 0, len(parts)+1)
			result = append(result, "brew", "uninstall", "--force")
			result = append(result, parts[2:]...)
			return single(strings.Join(result, " "))
		},
	})
}
