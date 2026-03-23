package rules

import (
	"regexp"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

var cpRe = regexp.MustCompile(`^cp `)

func init() {
	register(Rule{
		Name: "cp_omitting_directory",
		Match: func(cmd types.Command) bool {
			lower := strings.ToLower(cmd.Output)
			return strings.HasPrefix(cmd.Script, "cp") &&
				(strings.Contains(lower, "omitting directory") ||
					strings.Contains(lower, "is a directory"))
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			return single(cpRe.ReplaceAllString(cmd.Script, "cp -a "))
		},
	})
}
