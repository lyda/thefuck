package rules

import (
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "vagrant_up",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			if len(parts) == 0 || parts[0] != "vagrant" {
				return false
			}
			return strings.Contains(strings.ToLower(cmd.Output), "run `vagrant up`")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			parts := cmd.ScriptParts()
			startAll := shellAnd("vagrant up", cmd.Script)
			if len(parts) >= 3 {
				machine := parts[2]
				return multi([]string{
					shellAnd("vagrant up "+machine, cmd.Script),
					startAll,
				})
			}
			return single(startAll)
		},
	})
}
