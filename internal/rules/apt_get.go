package rules

import (
	"os/exec"
	"regexp"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

var aptCmdRe = regexp.MustCompile(`(?:command '([^']+)'|'([^']+)') not found`)

func aptGetPkgName(cmd types.Command) string {
	m := aptCmdRe.FindStringSubmatch(cmd.Output)
	if len(m) >= 2 {
		for _, g := range m[1:] {
			if g != "" {
				return g
			}
		}
	}
	// Fallback: use first script part
	parts := cmd.ScriptParts()
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}

func init() {
	register(Rule{
		Name: "apt_get",
		Match: func(cmd types.Command) bool {
			if _, err := exec.LookPath("apt-get"); err != nil {
				return false
			}
			if !strings.Contains(strings.ToLower(cmd.Output), "not found") {
				return false
			}
			pkg := aptGetPkgName(cmd)
			if pkg == "" {
				return false
			}
			// Check whether the package exists in apt
			return exec.Command("apt-cache", "show", pkg).Run() == nil // #nosec G204 -- pkg is extracted from command output, not raw user input
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			pkg := aptGetPkgName(cmd)
			if pkg == "" {
				return nil
			}
			return single(shellAnd("sudo apt-get install "+pkg, cmd.Script))
		},
	})
}
