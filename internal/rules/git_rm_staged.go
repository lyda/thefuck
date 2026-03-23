package rules

import (
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "git_rm_staged",
		Match: func(cmd types.Command) bool {
			return strings.HasPrefix(cmd.Script, "git ") &&
				strings.Contains(cmd.Script, " rm ") &&
				strings.Contains(cmd.Output, "error: the following file has changes staged in the index") &&
				strings.Contains(cmd.Output, "use --cached to keep the file, or -f to force removal")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			parts := cmd.ScriptParts()
			idx := -1
			for i, p := range parts {
				if p == "rm" {
					idx = i
					break
				}
			}
			if idx < 0 {
				return nil
			}
			// Build --cached variant
			cached := make([]string, len(parts)+1)
			copy(cached, parts[:idx+1])
			cached[idx+1] = "--cached"
			copy(cached[idx+2:], parts[idx+1:])

			// Build -f variant
			forced := make([]string, len(cached))
			copy(forced, cached)
			forced[idx+1] = "-f"

			return multi([]string{
				strings.Join(cached, " "),
				strings.Join(forced, " "),
			})
		},
	})
}
