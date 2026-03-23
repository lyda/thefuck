package rules

import (
	"slices"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

var gitHookBypassCommands = []string{"am", "commit", "push"}

func init() {
	register(Rule{
		Name:     "git_hook_bypass",
		Priority: 1100,
		Match: func(cmd types.Command) bool {
			if !strings.HasPrefix(cmd.Script, "git ") {
				return false
			}
			parts := cmd.ScriptParts()
			for _, hooked := range gitHookBypassCommands {
				if slices.Contains(parts, hooked) {
					return true
				}
			}
			return false
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			parts := cmd.ScriptParts()
			for _, hooked := range gitHookBypassCommands {
				if slices.Contains(parts, hooked) {
					return single(replaceArgument(cmd.Script, hooked, hooked+" --no-verify"))
				}
			}
			return nil
		},
	})
}
