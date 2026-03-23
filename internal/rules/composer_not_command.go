package rules

import (
	"regexp"
	"slices"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

var composerBrokenCmdRe = regexp.MustCompile(`Command "([^"]*)" is not defined`)
var composerDidYouMeanRe = regexp.MustCompile(`Did you mean this\?[^\n]*\n\s*([^\n]*)`)
var composerDidYouMeanOneRe = regexp.MustCompile(`Did you mean one of these\?[^\n]*\n\s*([^\n]*)`)

func init() {
	register(Rule{
		Name: "composer_not_command",
		Match: func(cmd types.Command) bool {
			if cmd.ScriptParts()[0] != "composer" {
				return false
			}
			outputLower := strings.ToLower(cmd.Output)
			if strings.Contains(outputLower, "did you mean this?") ||
				strings.Contains(outputLower, "did you mean one of these?") {
				return true
			}
			hasInstall := slices.Contains(cmd.ScriptParts(), "install")
			return hasInstall && strings.Contains(outputLower, "composer require")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			hasInstall := slices.Contains(cmd.ScriptParts(), "install")
			if hasInstall && strings.Contains(strings.ToLower(cmd.Output), "composer require") {
				return single(replaceArgument(cmd.Script, "install", "require"))
			}

			m := composerBrokenCmdRe.FindStringSubmatch(cmd.Output)
			if len(m) < 2 {
				return nil
			}
			brokenCmd := m[1]

			var newCmd string
			if mm := composerDidYouMeanRe.FindStringSubmatch(cmd.Output); len(mm) >= 2 {
				newCmd = strings.TrimSpace(mm[1])
			} else if mm := composerDidYouMeanOneRe.FindStringSubmatch(cmd.Output); len(mm) >= 2 {
				newCmd = strings.TrimSpace(mm[1])
			} else {
				return nil
			}
			return single(replaceArgument(cmd.Script, brokenCmd, newCmd))
		},
	})
}
