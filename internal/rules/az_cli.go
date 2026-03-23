package rules

import (
	"regexp"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

// Python: "(?=az)(?:.*): '(.*)' is not in the '.*' command group."
// The positive lookahead (?=az) just asserts the string starts with "az"; we
// check ScriptParts()[0] == "az" in Match instead and simplify the regex.
var azInvalidChoiceRe = regexp.MustCompile(`(?:.*): '([^']*)' is not in the '[^']*' command group\.`)

// Python: "^The most similar choice to '.*' is:\n\s*(.*)$" (MULTILINE)
var azOptionsRe = regexp.MustCompile(`(?m)^The most similar choice to '[^']*' is:\n\s*(.*)$`)

func init() {
	register(Rule{
		Name: "az_cli",
		Match: func(cmd types.Command) bool {
			return len(cmd.ScriptParts()) > 0 &&
				cmd.ScriptParts()[0] == "az" &&
				strings.Contains(cmd.Output, "is not in the") &&
				strings.Contains(cmd.Output, "command group")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			mistakeM := azInvalidChoiceRe.FindStringSubmatch(cmd.Output)
			if len(mistakeM) < 2 {
				return nil
			}
			mistake := mistakeM[1]
			optionMatches := azOptionsRe.FindAllStringSubmatch(cmd.Output, -1)
			scripts := make([]string, 0, len(optionMatches))
			for _, m := range optionMatches {
				scripts = append(scripts, replaceArgument(cmd.Script, mistake, strings.TrimSpace(m[1])))
			}
			return multi(scripts)
		},
	})
}
