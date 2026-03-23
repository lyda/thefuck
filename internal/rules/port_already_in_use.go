package rules

import (
	"os/exec"
	"regexp"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

var portPatterns = []*regexp.Regexp{
	regexp.MustCompile(`bind on address \('.*', (?P<port>\d+)\)`),
	regexp.MustCompile(`Unable to bind [^ ]*:(?P<port>\d+)`),
	regexp.MustCompile(`can't listen on port (?P<port>\d+)`),
	regexp.MustCompile(`listen EADDRINUSE [^ ]*:(?P<port>\d+)`),
}

// getUsedPort extracts the port number from the command output.
func getUsedPort(output string) string {
	for _, re := range portPatterns {
		m := re.FindStringSubmatch(output)
		if m != nil {
			return groupByName(re, m, "port")
		}
	}
	return ""
}

// getPidByPort runs lsof to find the PID using the given port.
func getPidByPort(port string) string {
	out, err := exec.Command("lsof", "-i", ":"+port).Output() // #nosec G204
	if err != nil {
		return ""
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) < 2 {
		return ""
	}
	fields := strings.Fields(lines[1])
	if len(fields) < 2 {
		return ""
	}
	return fields[1]
}

func init() {
	register(Rule{
		Name: "port_already_in_use",
		Match: func(cmd types.Command) bool {
			port := getUsedPort(cmd.Output)
			if port == "" {
				return false
			}
			pid := getPidByPort(port)
			return pid != ""
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			port := getUsedPort(cmd.Output)
			if port == "" {
				return nil
			}
			pid := getPidByPort(port)
			if pid == "" {
				return nil
			}
			return single(shellAnd("kill "+pid, cmd.Script))
		},
	})
}
