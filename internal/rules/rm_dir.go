package rules

import (
	"regexp"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

var rmRe = regexp.MustCompile(`\brm (.*)`)

func init() {
	register(Rule{
		Name: "rm_dir",
		Match: func(cmd types.Command) bool {
			return strings.Contains(cmd.Script, "rm") &&
				strings.Contains(strings.ToLower(cmd.Output), "is a directory")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			return single(rmRe.ReplaceAllString(cmd.Script, "rm -rf $1"))
		},
	})
}
