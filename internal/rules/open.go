package rules

import (
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

var openAppNames = []string{"open", "xdg-open", "gnome-open", "kde-open"}

func isOpenApp(cmd types.Command) bool {
	parts := cmd.ScriptParts()
	if len(parts) == 0 {
		return false
	}
	for _, name := range openAppNames {
		if parts[0] == name {
			return true
		}
	}
	return false
}

func isArgURL(cmd types.Command) bool {
	return strings.Contains(cmd.Script, ".com") ||
		strings.Contains(cmd.Script, ".edu") ||
		strings.Contains(cmd.Script, ".info") ||
		strings.Contains(cmd.Script, ".io") ||
		strings.Contains(cmd.Script, ".ly") ||
		strings.Contains(cmd.Script, ".me") ||
		strings.Contains(cmd.Script, ".net") ||
		strings.Contains(cmd.Script, ".org") ||
		strings.Contains(cmd.Script, ".se") ||
		strings.Contains(cmd.Script, "www.")
}

func init() {
	register(Rule{
		Name: "open",
		Match: func(cmd types.Command) bool {
			if !isOpenApp(cmd) {
				return false
			}
			output := strings.TrimSpace(cmd.Output)
			return isArgURL(cmd) ||
				(strings.HasPrefix(output, "The file ") && strings.HasSuffix(output, " does not exist."))
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			output := strings.TrimSpace(cmd.Output)
			if isArgURL(cmd) {
				// Prepend http:// to the argument by replacing "open " with "open http://"
				appName := cmd.ScriptParts()[0]
				return single(strings.Replace(cmd.Script, appName+" ", appName+" http://", 1))
			}
			if strings.HasPrefix(output, "The file ") && strings.HasSuffix(output, " does not exist.") {
				// Suggest creating the file/directory first
				idx := strings.Index(cmd.Script, " ")
				if idx < 0 {
					return nil
				}
				arg := cmd.Script[idx+1:]
				return multi([]string{
					shellAnd("touch "+arg, cmd.Script),
					shellAnd("mkdir "+arg, cmd.Script),
				})
			}
			return nil
		},
	})
}
