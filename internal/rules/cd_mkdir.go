package rules

import (
	"regexp"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

var cdRe = regexp.MustCompile(`^cd (.*)`)

func init() {
	register(Rule{
		Name: "cd_mkdir",
		Match: func(cmd types.Command) bool {
			if !strings.HasPrefix(cmd.Script, "cd ") {
				return false
			}
			lower := strings.ToLower(cmd.Output)
			return strings.Contains(lower, "no such file or directory") ||
				strings.Contains(lower, "cd: can't cd to") ||
				strings.Contains(lower, "does not exist")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			m := cdRe.FindStringSubmatch(cmd.Script)
			if len(m) < 2 {
				return nil
			}
			dir := m[1]
			return single(shellAnd("mkdir -p "+dir, "cd "+dir))
		},
	})
}
