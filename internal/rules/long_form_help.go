package rules

import (
	"regexp"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

var longFormHelpRe = regexp.MustCompile(`(?i)(?:Run|Try) '([^']+)'(?: or '[^']+')? for (?:details|more information)\.`)

func init() {
	register(Rule{
		Name: "long_form_help",
		Match: func(cmd types.Command) bool {
			if longFormHelpRe.MatchString(cmd.Output) {
				return true
			}
			return strings.Contains(cmd.Output, "--help")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			if m := longFormHelpRe.FindStringSubmatch(cmd.Output); len(m) >= 2 {
				return single(m[1])
			}
			return single(replaceArgument(cmd.Script, "-h", "--help"))
		},
	})
}
