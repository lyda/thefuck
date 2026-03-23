package rules

import (
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "mvn_no_command",
		Match: func(cmd types.Command) bool {
			return cmd.ScriptParts()[0] == "mvn" &&
				strings.Contains(cmd.Output, "No goals have been specified for this build")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			return multi([]string{
				cmd.Script + " clean package",
				cmd.Script + " clean install",
			})
		},
	})
}
