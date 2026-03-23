package rules

import (
	"regexp"

	"github.com/lyda/thefuck/internal/types"
)

var yarnCommandReplacedRe = regexp.MustCompile(`Run "(.*)" instead`)

func init() {
	register(Rule{
		Name: "yarn_command_replaced",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			if len(parts) < 1 || parts[0] != "yarn" {
				return false
			}
			return yarnCommandReplacedRe.MatchString(cmd.Output)
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			m := yarnCommandReplacedRe.FindStringSubmatch(cmd.Output)
			if len(m) < 2 {
				return nil
			}
			return single(m[1])
		},
	})
}
