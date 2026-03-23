package rules

import (
	"regexp"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

var tsuruBrokenCmdRe = regexp.MustCompile(`tsuru: "([^"]*)" is not a tsuru command`)

func init() {
	register(Rule{
		Name: "tsuru_not_command",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			if len(parts) == 0 || parts[0] != "tsuru" {
				return false
			}
			return strings.Contains(cmd.Output, ` is not a tsuru command. See "tsuru help".`) &&
				strings.Contains(cmd.Output, "\nDid you mean?\n\t")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			m := tsuruBrokenCmdRe.FindStringSubmatch(cmd.Output)
			if m == nil {
				return nil
			}
			brokenCmd := m[1]
			matches := getAllMatchedCommands(cmd.Output, []string{"Did you mean?"})
			if len(matches) == 0 {
				return nil
			}
			var scripts []string
			for _, match := range matches {
				scripts = append(scripts, replaceArgument(cmd.Script, brokenCmd, match))
			}
			return multi(scripts)
		},
	})
}
