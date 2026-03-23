package rules

import (
	"os/exec"
	"regexp"
	"strings"
	"sync"

	"github.com/lyda/thefuck/internal/types"
)

var yarnCmdNotFoundRe = regexp.MustCompile(`error Command "([^"]+)" not found\.`)

// npmToYarnCommands maps npm command names to their yarn equivalents.
var npmToYarnCommands = map[string]string{
	"require": "add",
}

var (
	yarnTasksOnce sync.Once
	yarnTasks     []string
)

func getYarnTasks() []string {
	yarnTasksOnce.Do(func() {
		out, _ := exec.Command("yarn", "--help").Output() // #nosec G204
		shouldYield := false
		for _, line := range strings.Split(string(out), "\n") {
			trimmed := strings.TrimSpace(line)
			if strings.Contains(trimmed, "Commands:") {
				shouldYield = true
				continue
			}
			if shouldYield && strings.Contains(trimmed, "- ") {
				fields := strings.Fields(trimmed)
				if len(fields) > 0 {
					yarnTasks = append(yarnTasks, fields[len(fields)-1])
				}
			}
		}
	})
	return yarnTasks
}

func init() {
	register(Rule{
		Name: "yarn_command_not_found",
		Match: func(cmd types.Command) bool {
			if _, err := exec.LookPath("yarn"); err != nil {
				return false
			}
			return yarnCmdNotFoundRe.MatchString(cmd.Output)
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			m := yarnCmdNotFoundRe.FindStringSubmatch(cmd.Output)
			if len(m) < 2 {
				return nil
			}
			misspelled := m[1]

			// Check known npm->yarn mappings first
			if yarnCmd, ok := npmToYarnCommands[misspelled]; ok {
				return single(replaceArgument(cmd.Script, misspelled, yarnCmd))
			}

			tasks := getYarnTasks()
			closest := getCloseMatches(misspelled, tasks, 0.6)
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
