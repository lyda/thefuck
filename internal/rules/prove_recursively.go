package rules

import (
	"os"
	"slices"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func proveIsRecursive(part string) bool {
	if part == "--recurse" {
		return true
	}
	if !strings.HasPrefix(part, "--") && strings.HasPrefix(part, "-") && strings.Contains(part, "r") {
		return true
	}
	return false
}

func proveIsDir(part string) bool {
	if strings.HasPrefix(part, "-") {
		return false
	}
	info, err := os.Stat(part)
	return err == nil && info.IsDir()
}

func init() {
	register(Rule{
		Name: "prove_recursively",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			if len(parts) == 0 || parts[0] != "prove" {
				return false
			}
			if !strings.Contains(cmd.Output, "NOTESTS") {
				return false
			}
			if slices.ContainsFunc(parts[1:], proveIsRecursive) {
				return false
			}
			return slices.ContainsFunc(parts[1:], proveIsDir)
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			parts := cmd.ScriptParts()
			newParts := make([]string, 0, len(parts)+1)
			newParts = append(newParts, parts[0])
			newParts = append(newParts, "-r")
			newParts = append(newParts, parts[1:]...)
			return single(strings.Join(newParts, " "))
		},
	})
}
