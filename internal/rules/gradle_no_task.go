package rules

import (
	"os/exec"
	"regexp"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

var gradleTaskRe = regexp.MustCompile(`Task '(.*)' (is ambiguous|not found)`)

// getGradleTasks runs `<gradle> tasks` and returns the list of available tasks.
func getGradleTasks(gradle string) []string {
	out, err := exec.Command(gradle, "tasks").Output() // #nosec G204
	if err != nil {
		return nil
	}
	var tasks []string
	shouldYield := false
	for _, line := range strings.Split(string(out), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "----") {
			shouldYield = true
			continue
		}
		if trimmed == "" {
			shouldYield = false
			continue
		}
		if shouldYield && !strings.HasPrefix(trimmed, "All tasks runnable from root project") {
			tasks = append(tasks, strings.Fields(trimmed)[0])
		}
	}
	return tasks
}

func init() {
	register(Rule{
		Name: "gradle_no_task",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			if len(parts) == 0 {
				return false
			}
			gradle := parts[0]
			if gradle != "gradle" && gradle != "gradlew" {
				return false
			}
			if _, err := exec.LookPath(gradle); err != nil {
				return false
			}
			return gradleTaskRe.MatchString(cmd.Output)
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			m := gradleTaskRe.FindStringSubmatch(cmd.Output)
			if len(m) < 2 {
				return nil
			}
			wrongTask := m[1]
			parts := cmd.ScriptParts()
			if len(parts) == 0 {
				return nil
			}
			allTasks := getGradleTasks(parts[0])
			closest := getCloseMatches(wrongTask, allTasks, 0.6)
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
