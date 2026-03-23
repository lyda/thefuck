package rules

import (
	"os"
	"regexp"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

var pathFromHistoryPatterns = []*regexp.Regexp{
	regexp.MustCompile(`no such file or directory: (.*)$`),
	regexp.MustCompile(`cannot access '(.*)': No such file or directory`),
	regexp.MustCompile(`: (.*): No such file or directory`),
	regexp.MustCompile(`can't cd to (.*)$`),
}

func getPathFromHistoryDestination(cmd types.Command) string {
	for _, re := range pathFromHistoryPatterns {
		m := re.FindStringSubmatch(cmd.Output)
		if len(m) > 1 {
			dest := m[1]
			// Check if this destination appears in the script parts
			for _, part := range cmd.ScriptParts() {
				if part == dest {
					return dest
				}
			}
		}
	}
	return ""
}

func getAllAbsolutePathsFromHistory(cmd types.Command) []string {
	history := readHistory()
	counts := make(map[string]int)
	var order []string

	for _, line := range history {
		if line == cmd.Script {
			continue
		}
		parts := strings.Fields(line)
		for _, param := range parts[1:] {
			if strings.HasPrefix(param, "/") || strings.HasPrefix(param, "~") {
				p := strings.TrimRight(param, "/")
				if _, seen := counts[p]; !seen {
					order = append(order, p)
				}
				counts[p]++
			}
		}
	}

	// Return ordered by most common (stable sort by insertion order for ties)
	// Simple approach: just return in frequency order
	type kv struct {
		key   string
		count int
	}
	kvs := make([]kv, 0, len(order))
	for _, k := range order {
		kvs = append(kvs, kv{k, counts[k]})
	}
	// Stable sort by count descending
	for i := 1; i < len(kvs); i++ {
		for j := i; j > 0 && kvs[j].count > kvs[j-1].count; j-- {
			kvs[j], kvs[j-1] = kvs[j-1], kvs[j]
		}
	}
	result := make([]string, len(kvs))
	for i, kv := range kvs {
		result[i] = kv.key
	}
	return result
}

func expandHome(path string) string {
	if strings.HasPrefix(path, "~/") {
		home := os.Getenv("HOME")
		return home + path[1:]
	}
	if path == "~" {
		return os.Getenv("HOME")
	}
	return path
}

func init() {
	register(Rule{
		Name:     "path_from_history",
		Priority: 800,
		Match: func(cmd types.Command) bool {
			return getPathFromHistoryDestination(cmd) != ""
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			destination := getPathFromHistoryDestination(cmd)
			if destination == "" {
				return nil
			}
			paths := getAllAbsolutePathsFromHistory(cmd)
			var results []string
			for _, p := range paths {
				if strings.HasSuffix(p, destination) {
					expanded := expandHome(p)
					if info, err := os.Stat(expanded); err == nil && info != nil {
						results = append(results, replaceArgument(cmd.Script, destination, p))
					}
				}
			}
			return multi(results)
		},
	})
}
