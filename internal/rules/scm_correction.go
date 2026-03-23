package rules

import (
	"os"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

var scmWrongPatterns = map[string]string{
	"git": "fatal: Not a git repository",
	"hg":  "abort: no repository found",
}

var scmPathMap = map[string]string{
	".git": "git",
	".hg":  "hg",
}

func scmGetActual() string {
	for path, scm := range scmPathMap {
		info, err := os.Stat(path)
		if err == nil && info.IsDir() {
			return scm
		}
	}
	return ""
}

func init() {
	register(Rule{
		Name: "scm_correction",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			if len(parts) == 0 {
				return false
			}
			scm := parts[0]
			pattern, ok := scmWrongPatterns[scm]
			if !ok {
				return false
			}
			return strings.Contains(cmd.Output, pattern) && scmGetActual() != ""
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			actual := scmGetActual()
			if actual == "" {
				return nil
			}
			parts := cmd.ScriptParts()
			newParts := make([]string, len(parts))
			copy(newParts, parts)
			newParts[0] = actual
			return single(strings.Join(newParts, " "))
		},
	})
}
