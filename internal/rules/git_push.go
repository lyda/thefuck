package rules

import (
	"regexp"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

var gitPushSuggestionRe = regexp.MustCompile(`git push (.*)`)

func gitPushUpstreamIndex(parts []string) int {
	for i, p := range parts {
		if p == "--set-upstream" || p == "-u" {
			return i
		}
	}
	return -1
}

func init() {
	register(Rule{
		Name: "git_push",
		Match: func(cmd types.Command) bool {
			return strings.HasPrefix(cmd.Script, "git ") &&
				strings.Contains(cmd.Script, "push") &&
				strings.Contains(cmd.Output, "git push --set-upstream")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			parts := cmd.ScriptParts()

			upstreamIdx := gitPushUpstreamIndex(parts)
			if upstreamIdx >= 0 {
				// Remove --set-upstream/-u
				parts = append(parts[:upstreamIdx], parts[upstreamIdx+1:]...)
				// Remove its argument if present
				if len(parts) > upstreamIdx {
					parts = append(parts[:upstreamIdx], parts[upstreamIdx+1:]...)
				}
			} else {
				// Remove trailing non-flag positional args (repository, refspec)
				pushIdx := -1
				for i, p := range parts {
					if p == "push" {
						pushIdx = i
						break
					}
				}
				if pushIdx >= 0 {
					for len(parts) > pushIdx+1 && !strings.HasPrefix(parts[len(parts)-1], "-") {
						parts = parts[:len(parts)-1]
					}
				}
			}

			// Extract git's suggestion (last match in output)
			matches := gitPushSuggestionRe.FindAllStringSubmatch(cmd.Output, -1)
			if len(matches) == 0 {
				return nil
			}
			arguments := strings.TrimSpace(strings.ReplaceAll(matches[len(matches)-1][1], "'", `\'`))
			newScript := replaceArgument(strings.Join(parts, " "), "push", "push "+arguments)
			return single(newScript)
		},
	})
}
