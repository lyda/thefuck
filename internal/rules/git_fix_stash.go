package rules

import (
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

var stashCommands = []string{
	"apply",
	"branch",
	"clear",
	"drop",
	"list",
	"pop",
	"save",
	"show",
}

func init() {
	register(Rule{
		Name: "git_fix_stash",
		Match: func(cmd types.Command) bool {
			if !strings.HasPrefix(cmd.Script, "git ") {
				return false
			}
			parts := cmd.ScriptParts()
			if len(parts) <= 1 {
				return false
			}
			return parts[1] == "stash" && strings.Contains(cmd.Output, "usage:")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			parts := cmd.ScriptParts()
			if len(parts) < 3 {
				return nil
			}
			stashCmd := parts[2]
			matches := getCloseMatches(stashCmd, stashCommands, 0.6)
			if len(matches) > 0 {
				return single(replaceArgument(cmd.Script, stashCmd, matches[0]))
			}
			// Insert "save" before the mistyped subcommand position
			newParts := make([]string, 0, len(parts)+1)
			newParts = append(newParts, parts[:2]...)
			newParts = append(newParts, "save")
			newParts = append(newParts, parts[2:]...)
			return single(strings.Join(newParts, " "))
		},
	})
}
