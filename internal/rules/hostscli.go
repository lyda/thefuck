package rules

import (
	"regexp"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

var hostscliNoCommandRe = regexp.MustCompile(`Error: No such command "([^"]*)"`)

const hostscliNoWebsite = "hostscli.errors.WebsiteImportError"

func init() {
	register(Rule{
		Name: "hostscli",
		Match: func(cmd types.Command) bool {
			if cmd.ScriptParts()[0] != "hostscli" {
				return false
			}
			return strings.Contains(cmd.Output, "Error: No such command") ||
				strings.Contains(cmd.Output, hostscliNoWebsite)
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			if strings.Contains(cmd.Output, hostscliNoWebsite) {
				return single("hostscli websites")
			}

			m := hostscliNoCommandRe.FindStringSubmatch(cmd.Output)
			if len(m) < 2 {
				return nil
			}
			misspelled := m[1]
			knownCmds := []string{"block", "unblock", "websites", "block_all", "unblock_all"}
			matches := getCloseMatches(misspelled, knownCmds, 0.1)
			scripts := make([]string, 0, len(matches))
			for _, match := range matches {
				scripts = append(scripts, replaceArgument(cmd.Script, misspelled, match))
			}
			if len(scripts) == 0 {
				return nil
			}
			return multi(scripts)
		},
	})
}
