package rules

import (
	"regexp"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

var lfsUnknownCmdRe = regexp.MustCompile(`Error: unknown command "([^"]*)" for "git-lfs"`)

func init() {
	register(Rule{
		Name: "git_lfs_mistype",
		Match: func(cmd types.Command) bool {
			if !strings.HasPrefix(cmd.Script, "git ") {
				return false
			}
			return strings.Contains(cmd.Script, "lfs") &&
				strings.Contains(cmd.Output, "Did you mean this?")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			m := lfsUnknownCmdRe.FindStringSubmatch(cmd.Output)
			if len(m) < 2 {
				return nil
			}
			brokenCmd := m[1]
			matched := getAllMatchedCommands(cmd.Output, []string{"Did you mean", " for usage."})
			scripts := make([]string, 0, len(matched))
			for _, fix := range matched {
				scripts = append(scripts, replaceArgument(cmd.Script, brokenCmd, fix))
			}
			if len(scripts) == 0 {
				return nil
			}
			return multi(scripts)
		},
	})
}
