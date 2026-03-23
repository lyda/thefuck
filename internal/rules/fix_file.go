package rules

import (
	"fmt"
	"os"
	"regexp"

	"github.com/lyda/thefuck/internal/types"
)

// rawPatterns are the error-location patterns from the Python fix_file rule.
// Named groups: file, line, col (col is optional).
var fixFileRawPatterns = []string{
	// js, node:
	`^    at (?P<file>[^:\n]+):(?P<line>[0-9]+):(?P<col>[0-9]+)`,
	// cargo:
	`^   (?P<file>[^:\n]+):(?P<line>[0-9]+):(?P<col>[0-9]+)`,
	// python, thefuck:
	`^  File "(?P<file>[^:\n]+)", line (?P<line>[0-9]+)`,
	// awk:
	`^awk: (?P<file>[^:\n]+):(?P<line>[0-9]+):`,
	// git:
	`^fatal: bad config file line (?P<line>[0-9]+) in (?P<file>[^:\n]+)`,
	// llc:
	`^llc: (?P<file>[^:\n]+):(?P<line>[0-9]+):(?P<col>[0-9]+):`,
	// lua:
	`^lua: (?P<file>[^:\n]+):(?P<line>[0-9]+):`,
	// fish:
	`^(?P<file>[^:\n]+) \(line (?P<line>[0-9]+)\):`,
	// bash, sh, ssh:
	`^(?P<file>[^:\n]+): line (?P<line>[0-9]+): `,
	// cargo, clang, gcc, go, pep8, rustc:
	`^(?P<file>[^:\n]+):(?P<line>[0-9]+):(?P<col>[0-9]+)`,
	// ghc, make, ruby, zsh:
	`^(?P<file>[^:\n]+):(?P<line>[0-9]+):`,
	// perl:
	`at (?P<file>[^:\n]+) line (?P<line>[0-9]+)`,
}

var fixFilePatterns []*regexp.Regexp

func init() {
	for _, raw := range fixFileRawPatterns {
		fixFilePatterns = append(fixFilePatterns, regexp.MustCompile(`(?m)`+raw))
	}

	register(Rule{
		Name: "fix_file",
		Match: func(cmd types.Command) bool {
			if os.Getenv("EDITOR") == "" {
				return false
			}
			return searchFixFilePattern(cmd.Output) != nil
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			editor := os.Getenv("EDITOR")
			if editor == "" {
				return nil
			}
			m := searchFixFilePattern(cmd.Output)
			if m == nil {
				return nil
			}
			editorCall := fmt.Sprintf("%s %s +%s", editor, m.file, m.line)
			return single(shellAnd(editorCall, cmd.Script))
		},
	})
}

type fixFileResult struct {
	re   *regexp.Regexp
	file string
	line string
}

// searchFixFilePattern searches the output for a known error pattern whose
// file group points to an existing file. Returns nil if none found.
func searchFixFilePattern(output string) *fixFileResult {
	for _, re := range fixFilePatterns {
		m := re.FindStringSubmatch(output)
		if m == nil {
			continue
		}
		file := groupByName(re, m, "file")
		line := groupByName(re, m, "line")
		if file == "" || line == "" {
			continue
		}
		if _, err := os.Stat(file); err == nil {
			return &fixFileResult{re: re, file: file, line: line}
		}
	}
	return nil
}

// groupByName returns the value of a named capture group from a submatch.
func groupByName(re *regexp.Regexp, match []string, name string) string {
	for i, gname := range re.SubexpNames() {
		if gname == name && i < len(match) {
			return match[i]
		}
	}
	return ""
}
