package rules

import (
	"os/exec"
	"regexp"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

var gitCheckoutPathspecRe = regexp.MustCompile(`error: pathspec '([^']*)' did not match any file\(s\) known to git`)

func gitCheckoutBranches() []string {
	out, err := exec.Command("git", "branch", "-a", "--no-color", "--no-column").Output()
	if err != nil {
		return nil
	}
	var branches []string
	for line := range strings.SplitSeq(string(out), "\n") {
		if strings.Contains(line, "->") {
			continue
		}
		line = strings.TrimPrefix(line, "*")
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "remotes/") {
			parts := strings.SplitN(line, "/", 3)
			if len(parts) == 3 {
				line = parts[2]
			}
		}
		if line != "" {
			branches = append(branches, line)
		}
	}
	return branches
}

func init() {
	register(Rule{
		Name: "git_checkout",
		Match: func(cmd types.Command) bool {
			return strings.HasPrefix(cmd.Script, "git ") &&
				strings.Contains(cmd.Output, "did not match any file(s) known to git") &&
				!strings.Contains(cmd.Output, "Did you forget to 'git add'?")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			m := gitCheckoutPathspecRe.FindStringSubmatch(cmd.Output)
			if len(m) < 2 {
				return nil
			}
			missingFile := m[1]
			branches := gitCheckoutBranches()
			closest := getCloseMatches(missingFile, branches, 0.6)

			var scripts []string
			if len(closest) > 0 {
				scripts = append(scripts, replaceArgument(cmd.Script, missingFile, closest[0]))
			}

			parts := cmd.ScriptParts()
			if len(parts) > 1 && parts[1] == "checkout" {
				scripts = append(scripts, replaceArgument(cmd.Script, "checkout", "checkout -b"))
			}

			if len(scripts) == 0 {
				scripts = append(scripts, shellAnd("git branch "+missingFile, cmd.Script))
			}
			return multi(scripts)
		},
	})
}
