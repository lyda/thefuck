package rules

import (
	"os/exec"
	"regexp"
	"strings"
	"sync"

	"github.com/lyda/thefuck/internal/types"
)

var (
	gulpTasksOnce sync.Once
	gulpTasks     []string
)

var gulpTaskRe = regexp.MustCompile(`Task '(\w+)' is not in your gulpfile`)

// getGulpTasks runs `gulp --tasks-simple` and returns the list of available tasks.
func getGulpTasks() []string {
	gulpTasksOnce.Do(func() {
		out, err := exec.Command("gulp", "--tasks-simple").Output() // #nosec G204
		if err != nil {
			return
		}
		for _, line := range strings.Split(string(out), "\n") {
			if t := strings.TrimSpace(line); t != "" {
				gulpTasks = append(gulpTasks, t)
			}
		}
	})
	return gulpTasks
}

func init() {
	register(Rule{
		Name: "gulp_not_task",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			if len(parts) == 0 || parts[0] != "gulp" {
				return false
			}
			if _, err := exec.LookPath("gulp"); err != nil {
				return false
			}
			return strings.Contains(cmd.Output, "is not in your gulpfile") ||
				strings.Contains(cmd.Output, "Task function must be specified")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			m := gulpTaskRe.FindStringSubmatch(cmd.Output)
			if len(m) < 2 {
				return nil
			}
			wrongTask := m[1]
			tasks := getGulpTasks()
			closest := getCloseMatches(wrongTask, tasks, 0.6)
			scripts := make([]string, 0, len(closest))
			for _, t := range closest {
				scripts = append(scripts, replaceArgument(cmd.Script, wrongTask, t))
			}
			if len(scripts) == 0 {
				return nil
			}
			return multi(scripts)
		},
	})
}
