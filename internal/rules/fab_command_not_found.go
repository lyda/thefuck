package rules

import (
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

// getBetween collects lines between a start marker and an optional end marker,
// returning the first whitespace-delimited token of each non-empty line.
func getBetween(content, start, end string) []string {
	var result []string
	shouldCollect := false
	for line := range strings.SplitSeq(content, "\n") {
		if strings.Contains(line, start) {
			shouldCollect = true
			continue
		}
		if end != "" && strings.Contains(line, end) {
			return result
		}
		if shouldCollect && line != "" {
			token := strings.TrimSpace(line)
			if idx := strings.Index(token, " "); idx >= 0 {
				token = token[:idx]
			}
			result = append(result, token)
		}
	}
	return result
}

func init() {
	register(Rule{
		Name: "fab_command_not_found",
		Match: func(cmd types.Command) bool {
			return cmd.ScriptParts()[0] == "fab" &&
				strings.Contains(cmd.Output, "Warning: Command(s) not found:")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			notFoundCmds := getBetween(cmd.Output, "Warning: Command(s) not found:", "Available commands:")
			possibleCmds := getBetween(cmd.Output, "Available commands:", "")

			script := cmd.Script
			for _, notFound := range notFoundCmds {
				matches := getCloseMatches(notFound, possibleCmds, 0.6)
				if len(matches) == 0 {
					continue
				}
				fix := matches[0]
				script = strings.Replace(script, " "+notFound, " "+fix, 1)
			}
			return single(script)
		},
	})
}
