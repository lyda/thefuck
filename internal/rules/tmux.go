package rules

import (
	"regexp"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

var tmuxAmbiguousRe = regexp.MustCompile(`ambiguous command: (.*), could be: (.*)`)

func init() {
	register(Rule{
		Name: "tmux",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			if len(parts) == 0 || parts[0] != "tmux" {
				return false
			}
			return strings.Contains(cmd.Output, "ambiguous command:") &&
				strings.Contains(cmd.Output, "could be:")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			m := tmuxAmbiguousRe.FindStringSubmatch(cmd.Output)
			if m == nil {
				return nil
			}
			oldCmd := m[1]
			suggestionParts := strings.Split(m[2], ",")
			var suggestions []string
			for _, s := range suggestionParts {
				suggestions = append(suggestions, strings.TrimSpace(s))
			}
			var scripts []string
			for _, s := range suggestions {
				scripts = append(scripts, replaceArgument(cmd.Script, oldCmd, s))
			}
			return multi(scripts)
		},
	})
}
