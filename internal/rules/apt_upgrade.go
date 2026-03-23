package rules

import (
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "apt_upgrade",
		Match: func(cmd types.Command) bool {
			if cmd.Script != "apt list --upgradable" {
				return false
			}
			lines := strings.Split(strings.TrimSpace(cmd.Output), "\n")
			return len(lines) > 1
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			return single("apt upgrade")
		},
	})
}
