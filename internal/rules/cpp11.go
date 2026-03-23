package rules

import (
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "cpp11",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			if len(parts) == 0 {
				return false
			}
			bin := parts[0]
			if bin != "g++" && bin != "clang++" {
				return false
			}
			return strings.Contains(cmd.Output, "This file requires compiler and library support for the ISO C++ 2011 standard.") ||
				strings.Contains(cmd.Output, "-Wc++11-extensions")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			return single(cmd.Script + " -std=c++11")
		},
	})
}
