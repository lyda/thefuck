package rules

import (
	"os"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "chmod_x",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			if len(parts) == 0 || !strings.HasPrefix(parts[0], "./") {
				return false
			}
			if !strings.Contains(strings.ToLower(cmd.Output), "permission denied") {
				return false
			}
			info, err := os.Stat(parts[0])
			if err != nil {
				return false
			}
			return info.Mode()&0o111 == 0
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			parts := cmd.ScriptParts()
			// Strip leading "./" from filename for chmod
			file := parts[0][2:]
			return single(shellAnd("chmod +x "+file, cmd.Script))
		},
	})
}
