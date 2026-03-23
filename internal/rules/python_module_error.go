package rules

import (
	"regexp"

	"github.com/lyda/thefuck/internal/types"
)

var missingModuleRe = regexp.MustCompile(`ModuleNotFoundError: No module named '([^']+)'`)

func init() {
	register(Rule{
		Name: "python_module_error",
		Match: func(cmd types.Command) bool {
			return missingModuleRe.MatchString(cmd.Output)
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			m := missingModuleRe.FindStringSubmatch(cmd.Output)
			if m == nil {
				return nil
			}
			return single(shellAnd("pip install "+m[1], cmd.Script))
		},
	})
}
