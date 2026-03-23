package rules

import (
	"regexp"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

var awsInvalidChoiceRe = regexp.MustCompile(`Invalid choice: '([^']*)', maybe you meant:`)
var awsOptionsRe = regexp.MustCompile(`(?m)^\s*\*\s(.*)`)

func init() {
	register(Rule{
		Name: "aws_cli",
		Match: func(cmd types.Command) bool {
			return len(cmd.ScriptParts()) > 0 &&
				cmd.ScriptParts()[0] == "aws" &&
				strings.Contains(cmd.Output, "usage:") &&
				strings.Contains(cmd.Output, "maybe you meant:")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			mistakeM := awsInvalidChoiceRe.FindStringSubmatch(cmd.Output)
			if len(mistakeM) < 2 {
				return nil
			}
			mistake := mistakeM[1]
			optionMatches := awsOptionsRe.FindAllStringSubmatch(cmd.Output, -1)
			scripts := make([]string, 0, len(optionMatches))
			for _, m := range optionMatches {
				scripts = append(scripts, replaceArgument(cmd.Script, mistake, m[1]))
			}
			return multi(scripts)
		},
	})
}
