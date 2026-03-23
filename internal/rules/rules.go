package rules

import (
	"regexp"
	"sort"
	"strings"

	"github.com/lyda/thefuck/internal/shells"
	"github.com/lyda/thefuck/internal/types"
)

// Rule defines a correction rule.
type Rule struct {
	Name          string
	Priority      int
	Match         func(types.Command) bool
	GetNewCommand func(types.Command) []types.CorrectedCommand
}

// All contains all registered rules, populated by each rule file's init().
var All []Rule

// register appends a rule to All, filling in DefaultPriority if unset.
func register(r Rule) {
	if r.Priority == 0 {
		r.Priority = types.DefaultPriority
	}
	All = append(All, r)
}

// shellAnd joins commands with the detected shell's AND operator.
func shellAnd(commands ...string) string {
	return shells.Detect().And(commands...)
}

// single returns a one-element correction slice.
// Priority is set by the corrector based on Rule.Priority.
func single(script string) []types.CorrectedCommand {
	return []types.CorrectedCommand{{Script: script}}
}

// multi returns corrections for multiple scripts.
// Priorities are set by the corrector based on Rule.Priority and index.
func multi(scripts []string) []types.CorrectedCommand {
	out := make([]types.CorrectedCommand, len(scripts))
	for i, s := range scripts {
		out[i] = types.CorrectedCommand{Script: s}
	}
	return out
}

// replaceArgument replaces the first occurrence of `from` as a shell word
// with `to` in script. Mirrors Python's utils.replace_argument.
func replaceArgument(script, from, to string) string {
	re := regexp.MustCompile(` ` + regexp.QuoteMeta(from) + `$`)
	if re.MatchString(script) {
		return re.ReplaceAllString(script, " "+to)
	}
	return strings.Replace(script, " "+from+" ", " "+to+" ", 1)
}

// getAllMatchedCommands extracts suggestion lines that appear after any of the
// given separator strings in output. Mirrors Python's utils.get_all_matched_commands.
func getAllMatchedCommands(output string, separators []string) []string {
	var results []string
	yielding := false
	for line := range strings.SplitSeq(output, "\n") {
		isSep := false
		for _, sep := range separators {
			if strings.Contains(line, sep) {
				isSep = true
				yielding = true
				break
			}
		}
		if yielding && !isSep {
			if t := strings.TrimSpace(line); t != "" {
				results = append(results, t)
			}
		}
	}
	return results
}

// levenshtein computes the edit distance between two strings.
func levenshtein(a, b string) int {
	ra, rb := []rune(a), []rune(b)
	la, lb := len(ra), len(rb)
	if la == 0 {
		return lb
	}
	if lb == 0 {
		return la
	}
	prev := make([]int, lb+1)
	curr := make([]int, lb+1)
	for j := 0; j <= lb; j++ {
		prev[j] = j
	}
	for i := 1; i <= la; i++ {
		curr[0] = i
		for j := 1; j <= lb; j++ {
			cost := 1
			if ra[i-1] == rb[j-1] {
				cost = 0
			}
			d := prev[j] + 1
			if v := curr[j-1] + 1; v < d {
				d = v
			}
			if v := prev[j-1] + cost; v < d {
				d = v
			}
			curr[j] = d
		}
		prev, curr = curr, prev
	}
	return prev[lb]
}

// getCloseMatches returns possibilities with similarity >= cutoff to word,
// sorted by similarity descending. Similarity = 1 - dist/maxLen.
func getCloseMatches(word string, possibilities []string, cutoff float64) []string {
	type candidate struct {
		s   string
		sim float64
	}
	var matches []candidate
	for _, p := range possibilities {
		wl := len([]rune(word))
		pl := len([]rune(p))
		maxLen := max(pl, wl)
		if maxLen == 0 {
			continue
		}
		dist := levenshtein(word, p)
		sim := 1.0 - float64(dist)/float64(maxLen)
		if sim >= cutoff {
			matches = append(matches, candidate{p, sim})
		}
	}
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].sim > matches[j].sim
	})
	result := make([]string, len(matches))
	for i, m := range matches {
		result[i] = m.s
	}
	return result
}
