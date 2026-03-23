package rules

import (
	"os"
	"regexp"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

var gitAddPathspecRe = regexp.MustCompile(`error: pathspec '([^']*)' did not match any file\(s\) known to git\.`)

func gitAddMissingFile(cmd types.Command) string {
	m := gitAddPathspecRe.FindStringSubmatch(cmd.Output)
	if len(m) < 2 {
		return ""
	}
	path := m[1]
	if _, err := os.Stat(path); err != nil {
		return ""
	}
	return path
}

func init() {
	register(Rule{
		Name: "git_add",
		Match: func(cmd types.Command) bool {
			return strings.HasPrefix(cmd.Script, "git ") &&
				strings.Contains(cmd.Output, "did not match any file(s) known to git.") &&
				gitAddMissingFile(cmd) != ""
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			f := gitAddMissingFile(cmd)
			if f == "" {
				return nil
			}
			return single(shellAnd("git add -- "+f, cmd.Script))
		},
	})
}
