package rules

import (
	"regexp"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

var pacmanInvalidOptionRe = regexp.MustCompile(` -[dfqrstuv]`)

func init() {
	register(Rule{
		Name: "pacman_invalid_option",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			if len(parts) == 0 || parts[0] != "pacman" {
				return false
			}
			if !strings.HasPrefix(cmd.Output, "error: invalid option '-") {
				return false
			}
			for _, option := range []string{"s", "u", "r", "q", "f", "d", "v", "t"} {
				if strings.Contains(cmd.Script, " -"+option) {
					return true
				}
			}
			return false
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			m := pacmanInvalidOptionRe.FindString(cmd.Script)
			if m == "" {
				return nil
			}
			return single(strings.Replace(cmd.Script, m, strings.ToUpper(m), 1))
		},
	})
}
