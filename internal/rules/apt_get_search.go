package rules

import (
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "apt_get_search",
		Match: func(cmd types.Command) bool {
			return cmd.ScriptParts()[0] == "apt-get" &&
				strings.HasPrefix(cmd.Script, "apt-get search")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			return single(strings.Replace(cmd.Script, "apt-get", "apt-cache", 1))
		},
	})
}
