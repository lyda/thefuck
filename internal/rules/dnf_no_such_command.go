package rules

import (
	"os/exec"
	"regexp"
	"strings"
	"sync"

	"github.com/lyda/thefuck/internal/types"
)

var dnfNoSuchCmdRe = regexp.MustCompile(`No such command: (.*)\.`)

var (
	dnfOperationsOnce sync.Once
	dnfOperations     []string
)

func getDNFOperations() []string {
	dnfOperationsOnce.Do(func() {
		out, _ := exec.Command("dnf", "--help").Output() // #nosec G204
		opRe := regexp.MustCompile(`(?m)^([a-z-]+) +`)
		for _, m := range opRe.FindAllStringSubmatch(string(out), -1) {
			dnfOperations = append(dnfOperations, m[1])
		}
	})
	return dnfOperations
}

func init() {
	register(Rule{
		Name: "dnf_no_such_command",
		Match: func(cmd types.Command) bool {
			if _, err := exec.LookPath("dnf"); err != nil {
				return false
			}
			return strings.Contains(strings.ToLower(cmd.Output), "no such command")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			m := dnfNoSuchCmdRe.FindStringSubmatch(cmd.Output)
			if len(m) < 2 {
				return nil
			}
			misspelled := m[1]
			closest := getCloseMatches(misspelled, getDNFOperations(), 0.6)
			if len(closest) == 0 {
				return nil
			}
			var results []string
			for _, c := range closest {
				results = append(results, replaceArgument(cmd.Script, misspelled, c))
			}
			return multi(results)
		},
	})
}
