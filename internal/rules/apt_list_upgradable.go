package rules

import (
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "apt_list_upgradable",
		Match: func(cmd types.Command) bool {
			return cmd.ScriptParts()[0] == "apt" &&
				strings.Contains(cmd.Output, "apt list --upgradable")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			return single("apt list --upgradable")
		},
	})
}
