package rules

import (
	"os"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func grepGetActualFile(parts []string) string {
	for _, part := range parts[1:] {
		if _, err := os.Stat(part); err == nil {
			return part
		}
	}
	return ""
}

func init() {
	register(Rule{
		Name: "grep_arguments_order",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			if len(parts) == 0 {
				return false
			}
			app := parts[0]
			if app != "grep" && app != "egrep" {
				return false
			}
			return strings.Contains(cmd.Output, ": No such file or directory") &&
				grepGetActualFile(parts) != ""
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			parts := cmd.ScriptParts()
			actualFile := grepGetActualFile(parts)
			if actualFile == "" {
				return nil
			}
			// Move the file to the end
			newParts := make([]string, 0, len(parts))
			for _, p := range parts {
				if p != actualFile {
					newParts = append(newParts, p)
				}
			}
			newParts = append(newParts, actualFile)
			return single(strings.Join(newParts, " "))
		},
	})
}
