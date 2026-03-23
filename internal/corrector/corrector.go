package corrector

import (
	"sort"

	"github.com/lyda/thefuck/internal/rules"
	"github.com/lyda/thefuck/internal/types"
)

// GetCorrectedCommands applies all registered rules to cmd and returns
// deduplicated, priority-sorted corrections (lowest priority value first).
func GetCorrectedCommands(cmd types.Command) []types.CorrectedCommand {
	var results []types.CorrectedCommand
	seen := make(map[string]bool)

	for _, rule := range rules.All {
		if rule.Match(cmd) {
			for i, c := range rule.GetNewCommand(cmd) {
				// Rule priority scales each correction; +1 so first suggestion
				// (i=0) gets exactly rule.Priority, second gets 2×, etc.
				c.Priority = rule.Priority * (i + 1)
				if !seen[c.Script] {
					seen[c.Script] = true
					results = append(results, c)
				}
			}
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Priority < results[j].Priority
	})
	return results
}
