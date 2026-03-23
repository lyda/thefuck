package rules

import (
	"regexp"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

var terraformMistakeRe = regexp.MustCompile(`Terraform has no command named "([^"]+)"`)
var terraformFixRe = regexp.MustCompile(`Did you mean "([^"]+)"\?`)

func init() {
	register(Rule{
		Name: "terraform_no_command",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			if len(parts) == 0 || parts[0] != "terraform" {
				return false
			}
			return terraformMistakeRe.MatchString(cmd.Output) &&
				terraformFixRe.MatchString(cmd.Output)
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			mMistake := terraformMistakeRe.FindStringSubmatch(cmd.Output)
			mFix := terraformFixRe.FindStringSubmatch(cmd.Output)
			if len(mMistake) < 2 || len(mFix) < 2 {
				return nil
			}
			return single(strings.Replace(cmd.Script, mMistake[1], mFix[1], 1))
		},
	})
}
