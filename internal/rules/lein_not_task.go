package rules

import (
	"regexp"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

var leinBrokenCmdRe = regexp.MustCompile(`'([^']*)' is not a task`)

func init() {
	register(Rule{
		Name: "lein_not_task",
		Match: func(cmd types.Command) bool {
			return cmd.ScriptParts()[0] == "lein" &&
				strings.Contains(cmd.Output, "is not a task. See 'lein help'") &&
				strings.Contains(cmd.Output, "Did you mean this?")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			m := leinBrokenCmdRe.FindStringSubmatch(cmd.Output)
			if len(m) < 2 {
				return nil
			}
			brokenCmd := m[1]
			newCmds := getAllMatchedCommands(cmd.Output, []string{"Did you mean this?"})
			scripts := make([]string, 0, len(newCmds))
			for _, suggestion := range newCmds {
				scripts = append(scripts, replaceArgument(cmd.Script, brokenCmd, suggestion))
			}
			if len(scripts) == 0 {
				return nil
			}
			return multi(scripts)
		},
	})
}
