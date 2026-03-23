package rules

import (
	"regexp"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

var gitFlagAfterFilenameRe1 = regexp.MustCompile(`fatal: bad flag '(.*?)' used after filename`)
var gitFlagAfterFilenameRe2 = regexp.MustCompile(`fatal: option '(.*?)' must come before non-option arguments`)

func gitFlagAfterFilenameBadFlag(output string) string {
	if m := gitFlagAfterFilenameRe1.FindStringSubmatch(output); len(m) >= 2 {
		return m[1]
	}
	if m := gitFlagAfterFilenameRe2.FindStringSubmatch(output); len(m) >= 2 {
		return m[1]
	}
	return ""
}

func init() {
	register(Rule{
		Name: "git_flag_after_filename",
		Match: func(cmd types.Command) bool {
			return strings.HasPrefix(cmd.Script, "git ") &&
				gitFlagAfterFilenameBadFlag(cmd.Output) != ""
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			badFlag := gitFlagAfterFilenameBadFlag(cmd.Output)
			if badFlag == "" {
				return nil
			}

			parts := cmd.ScriptParts()

			// Find the index of the bad flag
			badFlagIndex := -1
			for i, p := range parts {
				if p == badFlag {
					badFlagIndex = i
					break
				}
			}
			if badFlagIndex < 0 {
				return nil
			}

			// Find the last non-flag argument before the bad flag (the filename)
			filenameIndex := -1
			for i := badFlagIndex - 1; i >= 0; i-- {
				if !strings.HasPrefix(parts[i], "-") {
					filenameIndex = i
					break
				}
			}
			if filenameIndex < 0 {
				return nil
			}

			// Swap bad flag and filename
			newParts := make([]string, len(parts))
			copy(newParts, parts)
			newParts[badFlagIndex], newParts[filenameIndex] = newParts[filenameIndex], newParts[badFlagIndex]

			return single(strings.Join(newParts, " "))
		},
	})
}
