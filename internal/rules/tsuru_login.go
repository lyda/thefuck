package rules

import (
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "tsuru_login",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			return len(parts) > 0 && parts[0] == "tsuru" &&
				strings.Contains(cmd.Output, "not authenticated") &&
				strings.Contains(cmd.Output, "session has expired")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			return single(shellAnd("tsuru login", cmd.Script))
		},
	})
}
