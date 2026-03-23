package rules

import (
	"regexp"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

var mvnFailedLifecycleRe = regexp.MustCompile(`\[ERROR\] Unknown lifecycle phase "(.+)"`)
var mvnAvailableLifecyclesRe = regexp.MustCompile(`Available lifecycle phases are: (.+) -> \[Help 1\]`)

func init() {
	register(Rule{
		Name: "mvn_unknown_lifecycle_phase",
		Match: func(cmd types.Command) bool {
			if len(cmd.ScriptParts()) == 0 || cmd.ScriptParts()[0] != "mvn" {
				return false
			}
			return mvnFailedLifecycleRe.MatchString(cmd.Output) &&
				mvnAvailableLifecyclesRe.MatchString(cmd.Output)
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			failedM := mvnFailedLifecycleRe.FindStringSubmatch(cmd.Output)
			availableM := mvnAvailableLifecyclesRe.FindStringSubmatch(cmd.Output)
			if failedM == nil || availableM == nil {
				return nil
			}
			failedPhase := failedM[1]
			phases := strings.Split(availableM[1], ", ")
			matches := getCloseMatches(failedPhase, phases, 0.6)
			if len(matches) == 0 {
				return nil
			}
			var scripts []string
			for _, m := range matches {
				scripts = append(scripts, replaceArgument(cmd.Script, failedPhase, m))
			}
			return multi(scripts)
		},
	})
}
