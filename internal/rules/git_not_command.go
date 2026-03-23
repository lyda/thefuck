package rules

import (
	"regexp"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

var gitNotCmdRe = regexp.MustCompile(`git: '([^']*)' is not a git command`)

func init() {
	register(Rule{
		Name: "git_not_command",
		Match: func(cmd types.Command) bool {
			return strings.HasPrefix(cmd.Script, "git ") &&
				strings.Contains(cmd.Output, " is not a git command. See 'git --help'.") &&
				(strings.Contains(cmd.Output, "The most similar command") ||
					strings.Contains(cmd.Output, "Did you mean"))
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			m := gitNotCmdRe.FindStringSubmatch(cmd.Output)
			if len(m) < 2 {
				return nil
			}
			brokenCmd := m[1]
			matched := getAllMatchedCommands(cmd.Output, []string{"The most similar command", "Did you mean"})
			scripts := make([]string, 0, len(matched))
			for _, suggestion := range matched {
				scripts = append(scripts, replaceArgument(cmd.Script, brokenCmd, suggestion))
			}
			if len(scripts) == 0 {
				return nil
			}
			return multi(scripts)
		},
	})
}
