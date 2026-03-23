package rules

import (
	"os"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func lnGetDestination(parts []string) string {
	for _, part := range parts {
		if part == "ln" || part == "-s" || part == "--symbolic" {
			continue
		}
		if _, err := os.Stat(part); err == nil {
			return part
		}
	}
	return ""
}

func init() {
	register(Rule{
		Name: "ln_s_order",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			if len(parts) == 0 || parts[0] != "ln" {
				return false
			}
			hasSymbolic := false
			for _, p := range parts {
				if p == "-s" || p == "--symbolic" {
					hasSymbolic = true
					break
				}
			}
			return hasSymbolic &&
				strings.Contains(cmd.Output, "File exists") &&
				lnGetDestination(parts) != ""
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			parts := cmd.ScriptParts()
			destination := lnGetDestination(parts)
			if destination == "" {
				return nil
			}
			newParts := make([]string, 0, len(parts))
			for _, p := range parts {
				if p != destination {
					newParts = append(newParts, p)
				}
			}
			newParts = append(newParts, destination)
			return single(strings.Join(newParts, " "))
		},
	})
}
