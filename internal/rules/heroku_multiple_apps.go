package rules

import (
	"regexp"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

var herokuMultipleAppsRe = regexp.MustCompile(`([^ ]*) \([^)]*\)`)

func init() {
	register(Rule{
		Name: "heroku_multiple_apps",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			return len(parts) > 0 && parts[0] == "heroku" &&
				strings.Contains(cmd.Output, "https://devcenter.heroku.com/articles/multiple-environments")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			matches := herokuMultipleAppsRe.FindAllStringSubmatch(cmd.Output, -1)
			scripts := make([]string, 0, len(matches))
			for _, m := range matches {
				scripts = append(scripts, cmd.Script+" --app "+m[1])
			}
			if len(scripts) == 0 {
				return nil
			}
			return multi(scripts)
		},
	})
}
