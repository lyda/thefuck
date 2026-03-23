package rules

import (
	"regexp"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

var condaSuggestionRe = regexp.MustCompile(`'conda ([^']*)'`)

func init() {
	register(Rule{
		Name: "conda_mistype",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			if len(parts) < 1 || parts[0] != "conda" {
				return false
			}
			return strings.Contains(cmd.Output, "Did you mean 'conda")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			matches := condaSuggestionRe.FindAllStringSubmatch(cmd.Output, -1)
			if len(matches) < 2 {
				return nil
			}
			brokenCmd := matches[0][1]
			correctCmd := matches[1][1]
			return single(replaceArgument(cmd.Script, brokenCmd, correctCmd))
		},
	})
}
