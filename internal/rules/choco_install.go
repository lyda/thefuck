package rules

import (
	"slices"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "choco_install",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			hasCinst := slices.Contains(parts, "cinst")
			return (strings.HasPrefix(cmd.Script, "choco install") || hasCinst) &&
				strings.Contains(cmd.Output, "Installing the following packages")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			parts := cmd.ScriptParts()
			for _, part := range parts {
				if part == "choco" || part == "cinst" || part == "install" {
					continue
				}
				if strings.HasPrefix(part, "-") {
					continue
				}
				if strings.Contains(part, "=") || strings.Contains(part, "/") {
					continue
				}
				return single(strings.Replace(cmd.Script, part, part+".install", 1))
			}
			return nil
		},
	})
}
