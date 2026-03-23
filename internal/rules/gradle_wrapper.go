package rules

import (
	"os"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "gradle_wrapper",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			if len(parts) == 0 || parts[0] != "gradle" {
				return false
			}
			_, err := os.Stat("gradlew")
			return strings.Contains(cmd.Output, "not found") && err == nil
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			parts := cmd.ScriptParts()
			if len(parts) <= 1 {
				return single("./gradlew")
			}
			return single("./gradlew " + strings.Join(parts[1:], " "))
		},
	})
}
