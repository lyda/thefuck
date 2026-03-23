package rules

import (
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "git_rebase_merge_dir",
		Match: func(cmd types.Command) bool {
			return strings.Contains(cmd.Script, "rebase") &&
				strings.Contains(cmd.Output, "It seems that there is already a rebase-merge directory") &&
				strings.Contains(cmd.Output, "I wonder if you are in the middle of another rebase")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			commandList := []string{
				"git rebase --continue",
				"git rebase --abort",
				"git rebase --skip",
			}

			// The rm command is on line [-4] of the output (4th from last).
			lines := strings.Split(cmd.Output, "\n")
			if len(lines) >= 4 {
				rmCmd := strings.TrimSpace(lines[len(lines)-4])
				if rmCmd != "" {
					commandList = append(commandList, rmCmd)
				}
			}

			// Return all commands sorted by similarity to the original script
			// (cutoff 0 means include all).
			matches := getCloseMatches(cmd.Script, commandList, 0)
			if len(matches) == 0 {
				return multi(commandList)
			}
			return multi(matches)
		},
	})
}
