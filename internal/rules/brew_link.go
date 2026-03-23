package rules

import (
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "brew_link",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			if len(parts) < 2 || parts[0] != "brew" {
				return false
			}
			return (parts[1] == "ln" || parts[1] == "link") &&
				strings.Contains(cmd.Output, "brew link --overwrite --dry-run")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			parts := cmd.ScriptParts()
			// Build: brew link --overwrite --dry-run <rest>
			result := make([]string, 0, len(parts)+2)
			result = append(result, "brew", "link", "--overwrite", "--dry-run")
			result = append(result, parts[2:]...)
			return single(strings.Join(result, " "))
		},
	})
}
