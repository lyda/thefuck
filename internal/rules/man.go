package rules

import (
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "man",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			return len(parts) >= 1 && parts[0] == "man"
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			parts := cmd.ScriptParts()
			if len(parts) == 0 {
				return nil
			}

			// If already using section 3, suggest section 2 and vice versa.
			if strings.Contains(cmd.Script, "3") {
				return single(strings.Replace(cmd.Script, "3", "2", 1))
			}
			if strings.Contains(cmd.Script, "2") {
				return single(strings.Replace(cmd.Script, "2", "3", 1))
			}

			lastArg := parts[len(parts)-1]
			helpCommand := lastArg + " --help"

			// If there's no manual entry, suggest --help only.
			if strings.TrimSpace(cmd.Output) == "No manual entry for "+lastArg {
				return single(helpCommand)
			}

			// Build "man 3 <cmd>" and "man 2 <cmd>" by inserting the section
			// number after "man" (mirrors Python's list insert at index 1).
			// parts = ["man", ..., lastArg]
			// insert " 3 " at index 1 → "man" + " 3 " + rest joined
			beforeSection := parts[0]                    // "man"
			afterSection := strings.Join(parts[1:], " ") // everything after "man"

			cmd3 := beforeSection + " 3 " + afterSection
			cmd2 := beforeSection + " 2 " + afterSection

			return multi([]string{cmd3, cmd2, helpCommand})
		},
	})
}
