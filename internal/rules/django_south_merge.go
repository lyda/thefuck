package rules

import (
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "django_south_merge",
		Match: func(cmd types.Command) bool {
			return strings.Contains(cmd.Script, "manage.py") &&
				strings.Contains(cmd.Script, "migrate") &&
				strings.Contains(cmd.Output, "--merge: will just attempt the migration")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			return single(cmd.Script + " --merge")
		},
	})
}
