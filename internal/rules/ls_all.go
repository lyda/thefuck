package rules

import (
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "ls_all",
		Match: func(cmd types.Command) bool {
			return cmd.ScriptParts()[0] == "ls" &&
				strings.TrimSpace(cmd.Output) == ""
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			parts := cmd.ScriptParts()
			newParts := append([]string{"ls", "-A"}, parts[1:]...)
			return single(strings.Join(newParts, " "))
		},
	})
}
