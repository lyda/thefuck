package rules

import (
	"regexp"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

var yarnAliasRe = regexp.MustCompile("Did you mean [`\"](?:yarn )?([^`\"]*)[`\"]")

func init() {
	register(Rule{
		Name: "yarn_alias",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			return len(parts) >= 1 && parts[0] == "yarn" &&
				strings.Contains(cmd.Output, "Did you mean")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			parts := cmd.ScriptParts()
			if len(parts) < 2 {
				return nil
			}
			broken := parts[1]
			m := yarnAliasRe.FindStringSubmatch(cmd.Output)
			if len(m) < 2 {
				return nil
			}
			return single(replaceArgument(cmd.Script, broken, m[1]))
		},
	})
}
