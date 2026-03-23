package rules

import (
	"regexp"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

var pipUnknownCmdBrokenRe = regexp.MustCompile(`ERROR: unknown command "([^"]+)"`)
var pipUnknownCmdFixRe = regexp.MustCompile(`maybe you meant "([^"]+)"`)

func init() {
	register(Rule{
		Name: "pip_unknown_command",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			if len(parts) == 0 {
				return false
			}
			app := parts[0]
			if app != "pip" && app != "pip2" && app != "pip3" {
				return false
			}
			return strings.Contains(cmd.Script, "pip") &&
				strings.Contains(cmd.Output, "unknown command") &&
				strings.Contains(cmd.Output, "maybe you meant")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			mBroken := pipUnknownCmdBrokenRe.FindStringSubmatch(cmd.Output)
			mFix := pipUnknownCmdFixRe.FindStringSubmatch(cmd.Output)
			if len(mBroken) < 2 || len(mFix) < 2 {
				return nil
			}
			return single(replaceArgument(cmd.Script, mBroken[1], mFix[1]))
		},
	})
}
