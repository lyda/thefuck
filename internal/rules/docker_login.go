package rules

import (
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "docker_login",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			if len(parts) < 1 || parts[0] != "docker" {
				return false
			}
			return strings.Contains(cmd.Script, "docker") &&
				strings.Contains(cmd.Output, "access denied") &&
				strings.Contains(cmd.Output, "may require 'docker login'")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			return single(shellAnd("docker login", cmd.Script))
		},
	})
}
