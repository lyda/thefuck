package rules

import (
	"slices"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "npm_run_script",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			if len(parts) < 2 || parts[0] != "npm" {
				return false
			}
			if !strings.Contains(cmd.Output, "Usage: npm <command>") {
				return false
			}
			// Must NOT already have a "run"-prefixed part
			for _, p := range parts {
				if strings.HasPrefix(p, "ru") {
					return false
				}
			}
			// The subcommand (parts[1]) must be a known npm script
			scripts := getNpmScripts()
			return slices.Contains(scripts, parts[1])
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			parts := cmd.ScriptParts()
			// Insert "run-script" after "npm"
			newParts := make([]string, 0, len(parts)+1)
			newParts = append(newParts, parts[0], "run-script")
			newParts = append(newParts, parts[1:]...)
			return single(strings.Join(newParts, " "))
		},
	})
}
