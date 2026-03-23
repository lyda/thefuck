package rules

import (
	"regexp"

	"github.com/lyda/thefuck/internal/types"
)

var nixEnvRe = regexp.MustCompile(`nix-env -iA ([^\s]*)`)

func init() {
	register(Rule{
		Name: "nixos_cmd_not_found",
		Match: func(cmd types.Command) bool {
			return nixEnvRe.MatchString(cmd.Output)
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			m := nixEnvRe.FindStringSubmatch(cmd.Output)
			if m == nil {
				return nil
			}
			name := m[1]
			return single(shellAnd("nix-env -iA "+name, cmd.Script))
		},
	})
}
