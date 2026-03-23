package rules

import (
	"os/exec"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

// which returns true if the command exists in PATH.
func which(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func init() {
	register(Rule{
		Name:     "no_command",
		Priority: 3000,
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			if len(parts) == 0 {
				return false
			}
			if which(parts[0]) {
				return false
			}
			if !strings.Contains(cmd.Output, "not found") &&
				!strings.Contains(cmd.Output, "is not recognized as") {
				return false
			}
			matches := getCloseMatches(parts[0], getAllExecutables(), 0.6)
			return len(matches) > 0
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			parts := cmd.ScriptParts()
			if len(parts) == 0 {
				return nil
			}
			oldCmd := parts[0]

			// First try history: find the closest command used before
			history := readHistory()
			usedExecutables := make([]string, 0, len(history))
			for _, line := range history {
				if line == cmd.Script {
					continue
				}
				lineParts := strings.Fields(line)
				if len(lineParts) > 0 {
					usedExecutables = append(usedExecutables, lineParts[0])
				}
			}

			var newCmds []string
			historyMatches := getCloseMatches(oldCmd, usedExecutables, 0.6)
			if len(historyMatches) > 0 {
				newCmds = append(newCmds, historyMatches[0])
			}

			// Then add matches from all executables
			execMatches := getCloseMatches(oldCmd, getAllExecutables(), 0.6)
			inNew := make(map[string]bool)
			for _, c := range newCmds {
				inNew[c] = true
			}
			for _, c := range execMatches {
				if !inNew[c] {
					newCmds = append(newCmds, c)
				}
			}

			if len(newCmds) == 0 {
				return nil
			}

			results := make([]string, len(newCmds))
			for i, c := range newCmds {
				results[i] = strings.Replace(cmd.Script, oldCmd, c, 1)
			}
			return multi(results)
		},
	})
}
