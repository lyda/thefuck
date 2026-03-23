package rules

import (
	"regexp"
	"slices"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

var bisectBrokenRe = regexp.MustCompile(`git bisect ([^ $]*).*`)
var bisectUsageRe = regexp.MustCompile(`usage: git bisect \[([^\]]+)\]`)

func init() {
	register(Rule{
		Name: "git_bisect_usage",
		Match: func(cmd types.Command) bool {
			if !strings.HasPrefix(cmd.Script, "git ") {
				return false
			}
			parts := cmd.ScriptParts()
			hasBisect := slices.Contains(parts, "bisect")
			return hasBisect && strings.Contains(cmd.Output, "usage: git bisect")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			bm := bisectBrokenRe.FindStringSubmatch(cmd.Script)
			if len(bm) < 2 {
				return nil
			}
			broken := bm[1]
			um := bisectUsageRe.FindStringSubmatch(cmd.Output)
			if len(um) < 2 {
				return nil
			}
			options := strings.Split(um[1], "|")
			scripts := make([]string, 0, len(options))
			for _, opt := range options {
				opt = strings.TrimSpace(opt)
				scripts = append(scripts, replaceArgument(cmd.Script, broken, opt))
			}
			return multi(scripts)
		},
	})
}
