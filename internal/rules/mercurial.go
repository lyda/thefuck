package rules

import (
	"regexp"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

var mercurialDidYouMeanRe = regexp.MustCompile(`\n\(did you mean one of ([^\?]+)\?\)`)
var mercurialListRe = regexp.MustCompile(`\n    ([^$]+)$`)

func mercurialExtractPossibilities(output string) []string {
	if m := mercurialDidYouMeanRe.FindStringSubmatch(output); len(m) >= 2 {
		return strings.Split(m[1], ", ")
	}
	if m := mercurialListRe.FindStringSubmatch(output); len(m) >= 2 {
		return strings.Split(m[1], " ")
	}
	return nil
}

func init() {
	register(Rule{
		Name: "mercurial",
		Match: func(cmd types.Command) bool {
			if cmd.ScriptParts()[0] != "hg" {
				return false
			}
			return (strings.Contains(cmd.Output, "hg: unknown command") &&
				strings.Contains(cmd.Output, "(did you mean one of ")) ||
				(strings.Contains(cmd.Output, "hg: command '") &&
					strings.Contains(cmd.Output, "' is ambiguous:"))
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			parts := cmd.ScriptParts()
			if len(parts) < 2 {
				return nil
			}
			possibilities := mercurialExtractPossibilities(cmd.Output)
			if len(possibilities) == 0 {
				return nil
			}
			matches := getCloseMatches(parts[1], possibilities, 0.6)
			if len(matches) == 0 {
				// Fall back to first possibility
				matches = possibilities[:1]
			}
			newParts := make([]string, len(parts))
			copy(newParts, parts)
			newParts[1] = matches[0]
			return single(strings.Join(newParts, " "))
		},
	})
}
