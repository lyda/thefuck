package rules

import (
	"os"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "cat_dir",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			if len(parts) < 2 || parts[0] != "cat" {
				return false
			}
			if !strings.HasPrefix(cmd.Output, "cat: ") {
				return false
			}
			info, err := os.Stat(parts[1])
			return err == nil && info.IsDir()
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			return single(strings.Replace(cmd.Script, "cat", "ls", 1))
		},
	})
}
