package rules

import (
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "ln_no_hard_link",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			return len(parts) > 0 && parts[0] == "ln" &&
				strings.HasSuffix(cmd.Output, "hard link not allowed for directory")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			// Replace leading "ln " with "ln -s "
			if strings.HasPrefix(cmd.Script, "ln ") {
				return single("ln -s " + cmd.Script[3:])
			}
			return nil
		},
	})
}
