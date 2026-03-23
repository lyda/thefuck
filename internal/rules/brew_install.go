package rules

import (
	"os/exec"
	"regexp"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

var brewFormulaRe = regexp.MustCompile(`(?i)did you mean (.+)\?`)

func init() {
	register(Rule{
		Name: "brew_install",
		Match: func(cmd types.Command) bool {
			if _, err := exec.LookPath("brew"); err != nil {
				return false
			}
			return strings.Contains(cmd.Script, "brew") &&
				strings.Contains(cmd.Script, "install") &&
				strings.Contains(cmd.Output, "No available formula") &&
				(strings.Contains(cmd.Output, "Did you mean") ||
					strings.Contains(cmd.Output, "Searching formulae") ||
					strings.Contains(cmd.Output, "Searching taps"))
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			m := brewFormulaRe.FindStringSubmatch(cmd.Output)
			if len(m) < 2 {
				return nil
			}
			// Suggestions can be comma- or " or "-separated
			raw := m[1]
			raw = strings.ReplaceAll(raw, " or ", ", ")
			candidates := strings.Split(raw, ", ")
			scripts := make([]string, 0, len(candidates))
			for _, c := range candidates {
				c = strings.TrimSpace(strings.Trim(c, "\"'"))
				if c != "" {
					scripts = append(scripts, "brew install "+c)
				}
			}
			if len(scripts) == 0 {
				return nil
			}
			return multi(scripts)
		},
	})
}
