package rules

import (
	"slices"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "brew_cask_dependency",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			hasInstall := slices.Contains(parts, "install")
			return cmd.ScriptParts()[0] == "brew" &&
				hasInstall &&
				strings.Contains(cmd.Output, "brew cask install")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			var caskLines []string
			for line := range strings.SplitSeq(cmd.Output, "\n") {
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, "brew cask install") {
					caskLines = append(caskLines, line)
				}
			}
			if len(caskLines) == 0 {
				return nil
			}
			var brewCaskScript string
			if len(caskLines) > 1 {
				brewCaskScript = shellAnd(caskLines...)
			} else {
				brewCaskScript = caskLines[0]
			}
			return single(shellAnd(brewCaskScript, cmd.Script))
		},
	})
}
