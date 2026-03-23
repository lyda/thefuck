package rules

import (
	"regexp"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

var mkdirPRe = regexp.MustCompile(`\bmkdir (.*)`)

func init() {
	register(Rule{
		Name: "mkdir_p",
		Match: func(cmd types.Command) bool {
			return strings.Contains(cmd.Script, "mkdir") &&
				strings.Contains(cmd.Output, "No such file or directory")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			return single(mkdirPRe.ReplaceAllString(cmd.Script, "mkdir -p $1"))
		},
	})
}
