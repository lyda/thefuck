package rules

import (
	"net/url"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "git_clone_missing",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			// Must be a single token (the URL itself)
			if len(parts) != 1 {
				return false
			}
			// Must have produced the expected error
			if !strings.Contains(cmd.Output, "No such file or directory") &&
				!strings.Contains(cmd.Output, "not found") &&
				!strings.Contains(cmd.Output, "is not recognised as") {
				return false
			}

			script := cmd.Script
			u, err := url.Parse(script)
			if err != nil {
				return false
			}
			// Default to ssh when no scheme is present
			if u.Scheme == "" {
				u.Scheme = "ssh"
			}

			switch u.Scheme {
			case "http", "https":
				return u.Host != ""
			case "ssh":
				return strings.Contains(script, "@") && strings.Contains(script, ":")
			}
			return false
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			return single("git clone " + cmd.Script)
		},
	})
}
