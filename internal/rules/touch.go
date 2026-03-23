package rules

import (
	"regexp"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

var touchPathRe = regexp.MustCompile(`touch: (?:cannot touch ')?(.+)/.+'?:`)

func init() {
	register(Rule{
		Name: "touch",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			if len(parts) == 0 || parts[0] != "touch" {
				return false
			}
			return strings.Contains(cmd.Output, "No such file or directory")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			m := touchPathRe.FindStringSubmatch(cmd.Output)
			if m == nil {
				return nil
			}
			return single(shellAnd("mkdir -p "+m[1], cmd.Script))
		},
	})
}
