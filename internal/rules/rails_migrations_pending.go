package rules

import (
	"regexp"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

var railsMigrationSuggestionRe = regexp.MustCompile(`To resolve this issue, run:\s+(.*?)\n`)

func init() {
	register(Rule{
		Name: "rails_migrations_pending",
		Match: func(cmd types.Command) bool {
			return strings.Contains(cmd.Output, "Migrations are pending. To resolve this issue, run:")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			m := railsMigrationSuggestionRe.FindStringSubmatch(cmd.Output)
			if m == nil {
				return nil
			}
			return single(shellAnd(m[1], cmd.Script))
		},
	})
}
