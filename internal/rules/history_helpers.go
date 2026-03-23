package rules

import (
	"bufio"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/lyda/thefuck/internal/shells"
)

// readHistory returns a deduplicated list of history commands, most recent first.
// It reads from HISTFILE env var, or detects the shell and uses a default location.
// For fish shell it runs `fish -c 'builtin history'`.
func readHistory() []string {
	histFile := os.Getenv("HISTFILE")

	if histFile == "" {
		sh := shells.Detect()
		home := os.Getenv("HOME")
		switch sh.(type) {
		case shells.Zsh:
			histFile = filepath.Join(home, ".zsh_history")
		case shells.Fish:
			return readFishHistory()
		default:
			histFile = filepath.Join(home, ".bash_history")
		}
	}

	return readHistoryFile(histFile)
}

// readHistoryFile reads a history file (bash or zsh format) and returns
// deduplicated command strings, most recent first.
func readHistoryFile(path string) []string {
	f, err := os.Open(path) // #nosec G304,G703 -- path is derived from $HISTFILE or well-known shell defaults
	if err != nil {
		return nil
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		// Strip zsh extended history format: `: 1234567890:0;actual command`
		if strings.HasPrefix(line, ": ") {
			if idx := strings.Index(line, ";"); idx != -1 {
				line = line[idx+1:]
			}
		}
		line = strings.TrimSpace(line)
		if line != "" {
			lines = append(lines, line)
		}
	}

	// Reverse to get most recent first, then deduplicate preserving order.
	for i, j := 0, len(lines)-1; i < j; i, j = i+1, j-1 {
		lines[i], lines[j] = lines[j], lines[i]
	}
	return deduplicateStrings(lines)
}

// readFishHistory runs `fish -c 'builtin history'` and returns the entries.
func readFishHistory() []string {
	out, err := exec.Command("fish", "-c", "builtin history").Output() // #nosec G204
	if err != nil {
		return nil
	}
	var lines []string
	seen := map[string]bool{}
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if line != "" && !seen[line] {
			seen[line] = true
			lines = append(lines, line)
		}
	}
	return lines
}

// deduplicateStrings returns a slice with duplicates removed, preserving order.
func deduplicateStrings(in []string) []string {
	seen := make(map[string]bool, len(in))
	out := make([]string, 0, len(in))
	for _, s := range in {
		if !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	return out
}
