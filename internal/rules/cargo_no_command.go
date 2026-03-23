package rules

import (
	"regexp"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

var cargoDidYouMeanRe = regexp.MustCompile("Did you mean `([^`]*)`")

func init() {
	register(Rule{
		Name: "cargo_no_command",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			if len(parts) < 1 || parts[0] != "cargo" {
				return false
			}
			return strings.Contains(strings.ToLower(cmd.Output), "no such subcommand") &&
				strings.Contains(cmd.Output, "Did you mean")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			parts := cmd.ScriptParts()
			if len(parts) < 2 {
				return nil
			}
			m := cargoDidYouMeanRe.FindStringSubmatch(cmd.Output)
			if len(m) < 2 {
				return nil
			}
			return single(replaceArgument(cmd.Script, parts[1], m[1]))
		},
	})
}
