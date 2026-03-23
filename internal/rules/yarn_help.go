package rules

import (
	"os/exec"
	"regexp"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

var yarnHelpURLRe = regexp.MustCompile(`Visit ([^ ]*) for documentation about this command\.`)

// openCommand returns the platform command to open a URL in a browser.
func openCommand(url string) string {
	if _, err := exec.LookPath("xdg-open"); err == nil {
		return "xdg-open " + url
	}
	return "open " + url
}

func init() {
	register(Rule{
		Name: "yarn_help",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			if len(parts) < 2 || parts[0] != "yarn" {
				return false
			}
			return parts[1] == "help" &&
				strings.Contains(cmd.Output, "for documentation about this command.")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			m := yarnHelpURLRe.FindStringSubmatch(cmd.Output)
			if m == nil {
				return nil
			}
			return single(openCommand(m[1]))
		},
	})
}
