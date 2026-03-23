package rules

import (
	"os/exec"
	"regexp"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

var npmMissingScriptRe = regexp.MustCompile(`.*missing script: (.*)\n`)

// getNpmScripts runs "npm run-script" and returns the list of defined scripts.
func getNpmScripts() []string {
	out, err := exec.Command("npm", "run-script").Output()
	if err != nil {
		return nil
	}
	var scripts []string
	yielding := false
	for line := range strings.SplitSeq(string(out), "\n") {
		if strings.Contains(line, "available via `npm run-script`:") {
			yielding = true
			continue
		}
		if yielding {
			// Script lines are indented with exactly two spaces followed by a non-space
			if len(line) >= 2 && line[0] == ' ' && line[1] == ' ' && len(line) > 2 && line[2] != ' ' {
				parts := strings.Fields(line)
				if len(parts) > 0 {
					scripts = append(scripts, parts[0])
				}
			}
		}
	}
	return scripts
}

func init() {
	register(Rule{
		Name: "npm_missing_script",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			if len(parts) == 0 || parts[0] != "npm" {
				return false
			}
			hasRun := false
			for _, p := range parts {
				if strings.HasPrefix(p, "ru") {
					hasRun = true
					break
				}
			}
			return hasRun && strings.Contains(cmd.Output, "npm ERR! missing script: ")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			m := npmMissingScriptRe.FindStringSubmatch(cmd.Output)
			if m == nil {
				return nil
			}
			misspelled := m[1]
			scripts := getNpmScripts()
			if len(scripts) == 0 {
				return nil
			}
			var cmds []string
			for _, s := range scripts {
				cmds = append(cmds, replaceArgument(cmd.Script, misspelled, s))
			}
			return multi(cmds)
		},
	})
}
