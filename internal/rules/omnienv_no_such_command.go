package rules

import (
	"os/exec"
	"regexp"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

var omnienvNoSuchCmdRe = regexp.MustCompile(`env: no such command ['` + "`" + `]([^']*)['` + "`" + `]`)

var omnienvSupportedApps = []string{"goenv", "nodenv", "pyenv", "rbenv"}

// omnienvCommonTypos maps broken commands to known correct alternatives.
var omnienvCommonTypos = map[string][]string{
	"list":   {"versions", "install --list"},
	"remove": {"uninstall"},
}

func getOmnienvAppCommands(app string) []string {
	out, err := exec.Command(app, "commands").Output() // #nosec G204
	if err != nil {
		return nil
	}
	var cmds []string
	for _, line := range strings.Split(string(out), "\n") {
		if s := strings.TrimSpace(line); s != "" {
			cmds = append(cmds, s)
		}
	}
	return cmds
}

func init() {
	register(Rule{
		Name: "omnienv_no_such_command",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			if len(parts) == 0 {
				return false
			}
			app := parts[0]
			isSupported := false
			for _, a := range omnienvSupportedApps {
				if app == a {
					isSupported = true
					break
				}
			}
			if !isSupported {
				return false
			}
			if _, err := exec.LookPath(app); err != nil {
				return false
			}
			return strings.Contains(cmd.Output, "env: no such command ")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			m := omnienvNoSuchCmdRe.FindStringSubmatch(cmd.Output)
			if len(m) < 2 {
				return nil
			}
			broken := m[1]
			parts := cmd.ScriptParts()
			if len(parts) == 0 {
				return nil
			}
			app := parts[0]

			var results []string

			// Add common typo fixes first
			for _, alt := range omnienvCommonTypos[broken] {
				results = append(results, replaceArgument(cmd.Script, broken, alt))
			}

			// Add close matches from `app commands`
			appCmds := getOmnienvAppCommands(app)
			closest := getCloseMatches(broken, appCmds, 0.6)
			for _, c := range closest {
				results = append(results, replaceArgument(cmd.Script, broken, c))
			}

			if len(results) == 0 {
				return nil
			}
			return multi(results)
		},
	})
}
