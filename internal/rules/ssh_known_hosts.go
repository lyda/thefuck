package rules

import (
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

var (
	sshKnownHostsWarningRe = regexp.MustCompile(
		`WARNING: REMOTE HOST IDENTIFICATION HAS CHANGED!|` +
			`WARNING: POSSIBLE DNS SPOOFING DETECTED!|` +
			`Warning: the \S+ host key for '[^']+' differs from the key for the IP address '[^']+'`)

	sshOffendingRe = regexp.MustCompile(
		`(?:Offending (?:key for IP|\S+ key)|Matching host key) in ([^:]+):(\d+)`)
)

func init() {
	register(Rule{
		Name: "ssh_known_hosts",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			if len(parts) == 0 {
				return false
			}
			app := parts[0]
			if app != "ssh" && app != "scp" {
				return false
			}
			return sshKnownHostsWarningRe.MatchString(cmd.Output)
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			// Remove the offending key line(s) from known_hosts, then re-run.
			offending := sshOffendingRe.FindAllStringSubmatch(cmd.Output, -1)
			for _, m := range offending {
				if len(m) < 3 {
					continue
				}
				filePath := strings.TrimSpace(m[1])
				lineNo, err := strconv.Atoi(m[2])
				if err != nil || lineNo < 1 {
					continue
				}

				data, err := os.ReadFile(filePath) // #nosec G304 -- path is from ssh output
				if err != nil {
					continue
				}
				lines := strings.Split(string(data), "\n")
				idx := lineNo - 1
				if idx >= len(lines) {
					continue
				}
				lines = append(lines[:idx], lines[idx+1:]...)
				_ = os.WriteFile(filePath, []byte(strings.Join(lines, "\n")), 0600) // #nosec G703 -- filePath comes from SSH's own error output, not user input
			}
			return single(cmd.Script)
		},
	})
}
