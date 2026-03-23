package rules

import (
	"regexp"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

var herokuNotCmdRe = regexp.MustCompile(`Run heroku _ to run ([^.]*)`)

func init() {
	register(Rule{
		Name: "heroku_not_command",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			return len(parts) > 0 && parts[0] == "heroku" &&
				strings.Contains(cmd.Output, "Run heroku _ to run")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			m := herokuNotCmdRe.FindStringSubmatch(cmd.Output)
			if len(m) < 2 {
				return nil
			}
			return single(m[1])
		},
	})
}
