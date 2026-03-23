package rules

import (
	"regexp"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

var shellPromptLiteralRe = regexp.MustCompile(`^[\s]*\$ [\S]+`)

func init() {
	register(Rule{
		Name: "remove_shell_prompt_literal",
		Match: func(cmd types.Command) bool {
			return strings.Contains(cmd.Output, "$: command not found") &&
				shellPromptLiteralRe.MatchString(cmd.Script)
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			return single(strings.TrimLeft(cmd.Script, "$ "))
		},
	})
}
