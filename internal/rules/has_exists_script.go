package rules

import (
	"os"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "has_exists_script",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			if len(parts) == 0 {
				return false
			}
			_, err := os.Stat(parts[0])
			return err == nil && strings.Contains(cmd.Output, "command not found")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			return single("./" + cmd.Script)
		},
	})
}
