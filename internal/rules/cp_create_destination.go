package rules

import (
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "cp_create_destination",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			if len(parts) < 1 {
				return false
			}
			if parts[0] != "cp" && parts[0] != "mv" {
				return false
			}
			return strings.Contains(cmd.Output, "No such file or directory") ||
				(strings.HasPrefix(cmd.Output, "cp: directory") &&
					strings.HasSuffix(strings.TrimRight(cmd.Output, "\n"), "does not exist"))
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			parts := cmd.ScriptParts()
			if len(parts) == 0 {
				return nil
			}
			dest := parts[len(parts)-1]
			return single(shellAnd("mkdir -p "+dest, cmd.Script))
		},
	})
}
