package rules

import (
	"os/exec"
	"strings"
	"sync"

	"github.com/lyda/thefuck/internal/types"
)

var (
	yumOperationsOnce sync.Once
	yumOperations     []string
)

func getYumOperations() []string {
	yumOperationsOnce.Do(func() {
		out, _ := exec.Command("yum").CombinedOutput() // #nosec G204
		lines := strings.Split(string(out), "\n")
		collecting := false
		for _, line := range lines {
			if strings.HasPrefix(line, "List of Commands:") {
				collecting = true
				continue
			}
			if !collecting {
				continue
			}
			trimmed := strings.TrimSpace(line)
			if trimmed == "" {
				if len(yumOperations) > 0 {
					break
				}
				continue
			}
			// Each line is "  command    description" — skip the first two
			// blank lines after the header by only collecting non-empty
			fields := strings.Fields(trimmed)
			if len(fields) > 0 {
				yumOperations = append(yumOperations, fields[0])
			}
		}
	})
	return yumOperations
}

func init() {
	register(Rule{
		Name: "yum_invalid_operation",
		Match: func(cmd types.Command) bool {
			if _, err := exec.LookPath("yum"); err != nil {
				return false
			}
			return strings.Contains(cmd.Output, "No such command:")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			parts := cmd.ScriptParts()
			if len(parts) < 2 {
				return nil
			}
			invalidOp := parts[1]
			if invalidOp == "uninstall" {
				return single(strings.Replace(cmd.Script, "uninstall", "remove", 1))
			}
			closest := getCloseMatches(invalidOp, getYumOperations(), 0.6)
			if len(closest) == 0 {
				return nil
			}
			var results []string
			for _, c := range closest {
				results = append(results, replaceArgument(cmd.Script, invalidOp, c))
			}
			return multi(results)
		},
	})
}
