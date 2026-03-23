package rules

import (
	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name:     "history",
		Priority: 9999,
		Match: func(cmd types.Command) bool {
			history := readHistory()
			// Filter out the current script from history
			filtered := make([]string, 0, len(history))
			for _, h := range history {
				if h != cmd.Script {
					filtered = append(filtered, h)
				}
			}
			matches := getCloseMatches(cmd.Script, filtered, 0.6)
			return len(matches) > 0
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			history := readHistory()
			filtered := make([]string, 0, len(history))
			for _, h := range history {
				if h != cmd.Script {
					filtered = append(filtered, h)
				}
			}
			matches := getCloseMatches(cmd.Script, filtered, 0.6)
			if len(matches) == 0 {
				return nil
			}
			return single(matches[0])
		},
	})
}
