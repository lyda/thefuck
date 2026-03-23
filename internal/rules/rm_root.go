package rules

import (
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "rm_root",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			hasRm := false
			hasSlash := false
			for _, p := range parts {
				if p == "rm" {
					hasRm = true
				}
				if p == "/" {
					hasSlash = true
				}
			}
			return hasRm && hasSlash &&
				!strings.Contains(cmd.Script, "--no-preserve-root") &&
				strings.Contains(cmd.Output, "--no-preserve-root")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			return single(cmd.Script + " --no-preserve-root")
		},
	})
}
