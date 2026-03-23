package rules

import (
	"os/exec"
	"regexp"
	"strings"
	"sync"

	"github.com/lyda/thefuck/internal/types"
)

var (
	gruntTasksOnce sync.Once
	gruntTasks     []string
)

var gruntTaskRe = regexp.MustCompile(`Warning: Task "([^"]*)" not found\.`)

// getGruntTasks runs `grunt --help` and returns the list of available tasks.
func getGruntTasks() []string {
	gruntTasksOnce.Do(func() {
		out, err := exec.Command("grunt", "--help").Output() // #nosec G204
		if err != nil {
			return
		}
		shouldYield := false
		for _, line := range strings.Split(string(out), "\n") {
			trimmed := strings.TrimSpace(line)
			if strings.Contains(trimmed, "Available tasks") {
				shouldYield = true
				continue
			}
			if shouldYield && trimmed == "" {
				break
			}
			if shouldYield && strings.Contains(trimmed, "  ") {
				gruntTasks = append(gruntTasks, strings.Fields(trimmed)[0])
			}
		}
	})
	return gruntTasks
}

func init() {
	register(Rule{
		Name: "grunt_task_not_found",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			if len(parts) == 0 || parts[0] != "grunt" {
				return false
			}
			if _, err := exec.LookPath("grunt"); err != nil {
				return false
			}
			return gruntTaskRe.MatchString(cmd.Output)
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			m := gruntTaskRe.FindStringSubmatch(cmd.Output)
			if len(m) < 2 {
				return nil
			}
			// Strip subtask suffix (e.g. "task:subtask" -> "task")
			misspelled := strings.SplitN(m[1], ":", 2)[0]
			tasks := getGruntTasks()
			closest := getCloseMatches(misspelled, tasks, 0.6)
			if len(closest) == 0 {
				return nil
			}
			fixed := closest[0]
			scripts := make([]string, 0, len(closest))
			for _, c := range closest {
				scripts = append(scripts, strings.Replace(cmd.Script, " "+misspelled, " "+c, 1))
				_ = fixed
			}
			return multi(scripts)
		},
	})
}
