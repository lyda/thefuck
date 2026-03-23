package rules

import (
	"regexp"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

var noSuchFilePatterns = []*regexp.Regexp{
	regexp.MustCompile(`mv: cannot move '[^']*' to '([^']*)': No such file or directory`),
	regexp.MustCompile(`mv: cannot move '[^']*' to '([^']*)': Not a directory`),
	regexp.MustCompile(`cp: cannot create regular file '([^']*)': No such file or directory`),
	regexp.MustCompile(`cp: cannot create regular file '([^']*)': Not a directory`),
}

func init() {
	register(Rule{
		Name: "no_such_file",
		Match: func(cmd types.Command) bool {
			for _, re := range noSuchFilePatterns {
				if re.MatchString(cmd.Output) {
					return true
				}
			}
			return false
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			for _, re := range noSuchFilePatterns {
				m := re.FindStringSubmatch(cmd.Output)
				if len(m) >= 2 {
					file := m[1]
					idx := strings.LastIndex(file, "/")
					if idx < 0 {
						continue
					}
					dir := file[:idx]
					return single(shellAnd("mkdir -p "+dir, cmd.Script))
				}
			}
			return nil
		},
	})
}
