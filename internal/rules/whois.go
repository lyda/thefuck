package rules

import (
	"net/url"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "whois",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			return len(parts) >= 2 && parts[0] == "whois"
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			parts := cmd.ScriptParts()
			if len(parts) < 2 {
				return nil
			}
			arg := parts[1]

			if strings.Contains(cmd.Script, "/") {
				// Strip scheme and path — keep only the hostname.
				u, err := url.Parse(arg)
				if err != nil {
					return nil
				}
				host := u.Host
				if host == "" {
					// url.Parse puts scheme-less URLs in Path
					host = u.Path
				}
				return single("whois " + host)
			} else if strings.Contains(cmd.Script, ".") {
				// Remove successive left-most subdomains.
				u, err := url.Parse(arg)
				if err != nil {
					return nil
				}
				path := u.Path
				if path == "" {
					path = arg
				}
				dotParts := strings.Split(path, ".")
				var results []string
				for n := 1; n < len(dotParts); n++ {
					results = append(results, "whois "+strings.Join(dotParts[n:], "."))
				}
				return multi(results)
			}
			return nil
		},
	})
}
