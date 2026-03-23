package rules

import (
	"os/exec"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name:     "wrong_hyphen_before_subcommand",
		Priority: 4500,
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			if len(parts) == 0 {
				return false
			}
			firstPart := parts[0]
			if !strings.Contains(firstPart, "-") {
				return false
			}
			// firstPart itself must not be a valid executable
			if _, err := exec.LookPath(firstPart); err == nil {
				return false
			}
			// The part before the first hyphen must be a valid executable
			cmdBase := strings.SplitN(firstPart, "-", 2)[0]
			_, err := exec.LookPath(cmdBase)
			return err == nil
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			return single(strings.Replace(cmd.Script, "-", " ", 1))
		},
	})
}
