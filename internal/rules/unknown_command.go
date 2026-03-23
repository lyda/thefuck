package rules

import (
	"regexp"

	"github.com/lyda/thefuck/internal/types"
)

var unknownCmdBrokenRe = regexp.MustCompile(`([^:]*): Unknown command.*`)
var unknownCmdSuggestionRe = regexp.MustCompile(`Did you mean ([^?]*)?`)

func init() {
	register(Rule{
		Name: "unknown_command",
		Match: func(cmd types.Command) bool {
			return unknownCmdBrokenRe.MatchString(cmd.Output) &&
				unknownCmdSuggestionRe.MatchString(cmd.Output)
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			brokenM := unknownCmdBrokenRe.FindStringSubmatch(cmd.Output)
			if brokenM == nil {
				return nil
			}
			brokenCmd := brokenM[1]
			suggestMatches := unknownCmdSuggestionRe.FindAllStringSubmatch(cmd.Output, -1)
			if len(suggestMatches) == 0 {
				return nil
			}
			var scripts []string
			for _, m := range suggestMatches {
				scripts = append(scripts, replaceArgument(cmd.Script, brokenCmd, m[1]))
			}
			return multi(scripts)
		},
	})
}
