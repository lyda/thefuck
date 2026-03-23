package rules

import (
	"os/exec"
	"regexp"
	"strings"
	"sync"

	"github.com/lyda/thefuck/internal/types"
)

var gemUnknownCmdRe = regexp.MustCompile(`Unknown command (.*)$`)

var (
	gemCommandsOnce sync.Once
	gemCommands     []string
)

func getGemCommands() []string {
	gemCommandsOnce.Do(func() {
		out, _ := exec.Command("gem", "help", "commands").Output() // #nosec G204
		for _, line := range strings.Split(string(out), "\n") {
			if strings.HasPrefix(line, "    ") {
				fields := strings.Fields(line)
				if len(fields) > 0 {
					gemCommands = append(gemCommands, fields[0])
				}
			}
		}
	})
	return gemCommands
}

func init() {
	register(Rule{
		Name: "gem_unknown_command",
		Match: func(cmd types.Command) bool {
			if _, err := exec.LookPath("gem"); err != nil {
				return false
			}
			return strings.Contains(cmd.Output, "ERROR:  While executing gem ... (Gem::CommandLineError)") &&
				strings.Contains(cmd.Output, "Unknown command")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			m := gemUnknownCmdRe.FindStringSubmatch(cmd.Output)
			if len(m) < 2 {
				return nil
			}
			unknown := strings.TrimSpace(m[1])
			closest := getCloseMatches(unknown, getGemCommands(), 0.6)
			if len(closest) == 0 {
				return nil
			}
			var results []string
			for _, c := range closest {
				results = append(results, replaceArgument(cmd.Script, unknown, c))
			}
			return multi(results)
		},
	})
}
