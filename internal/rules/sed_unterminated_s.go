package rules

import (
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "sed_unterminated_s",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			if len(parts) == 0 || parts[0] != "sed" {
				return false
			}
			return strings.Contains(cmd.Output, "unterminated `s' command")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			parts := cmd.ScriptParts()
			fixed := make([]string, len(parts))
			copy(fixed, parts)
			for i, part := range fixed {
				if (strings.HasPrefix(part, "s/") || strings.HasPrefix(part, "-es/")) &&
					!strings.HasSuffix(part, "/") {
					fixed[i] = part + "/"
				}
			}
			// Quote parts that contain spaces or special chars (simple quoting)
			quoted := make([]string, len(fixed))
			for i, p := range fixed {
				if strings.ContainsAny(p, " \t'\"\\") {
					quoted[i] = "'" + strings.ReplaceAll(p, "'", "'\\''") + "'"
				} else {
					quoted[i] = p
				}
			}
			return single(strings.Join(quoted, " "))
		},
	})
}
