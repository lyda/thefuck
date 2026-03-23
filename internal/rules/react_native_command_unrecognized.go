package rules

import (
	"os/exec"
	"regexp"
	"strings"
	"sync"

	"github.com/lyda/thefuck/internal/types"
)

var (
	rnCommandsOnce sync.Once
	rnCommands     []string
)

var rnUnrecognizedRe = regexp.MustCompile(`Unrecognized command '(.*)'`)

// getReactNativeCommands runs `react-native --help` and parses available commands.
func getReactNativeCommands() []string {
	rnCommandsOnce.Do(func() {
		out, err := exec.Command("react-native", "--help").Output() // #nosec G204
		if err != nil {
			return
		}
		shouldYield := false
		for _, line := range strings.Split(string(out), "\n") {
			trimmed := strings.TrimSpace(line)
			if trimmed == "" {
				continue
			}
			if strings.Contains(trimmed, "Commands:") {
				shouldYield = true
				continue
			}
			if shouldYield {
				rnCommands = append(rnCommands, strings.Fields(trimmed)[0])
			}
		}
	})
	return rnCommands
}

func init() {
	register(Rule{
		Name: "react_native_command_unrecognized",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			if len(parts) == 0 || parts[0] != "react-native" {
				return false
			}
			if _, err := exec.LookPath("react-native"); err != nil {
				return false
			}
			return rnUnrecognizedRe.MatchString(cmd.Output)
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			m := rnUnrecognizedRe.FindStringSubmatch(cmd.Output)
			if len(m) < 2 {
				return nil
			}
			misspelled := m[1]
			commands := getReactNativeCommands()
			closest := getCloseMatches(misspelled, commands, 0.6)
			scripts := make([]string, 0, len(closest))
			for _, c := range closest {
				scripts = append(scripts, replaceArgument(cmd.Script, misspelled, c))
			}
			if len(scripts) == 0 {
				return nil
			}
			return multi(scripts)
		},
	})
}
