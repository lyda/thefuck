package rules

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/lyda/thefuck/internal/types"
)

var (
	allExecutablesOnce sync.Once
	allExecutables     []string
)

// getAllExecutables walks each directory in PATH and collects executable names.
func getAllExecutables() []string {
	allExecutablesOnce.Do(func() {
		seen := make(map[string]struct{})
		pathEnv := os.Getenv("PATH")
		for dir := range strings.SplitSeq(pathEnv, string(os.PathListSeparator)) {
			if dir == "" {
				continue
			}
			entries, err := os.ReadDir(dir)
			if err != nil {
				continue
			}
			for _, entry := range entries {
				if entry.IsDir() {
					continue
				}
				info, err := entry.Info()
				if err != nil {
					continue
				}
				// Check execute bit (owner, group, or other)
				if info.Mode()&0111 == 0 {
					continue
				}
				name := filepath.Base(entry.Name())
				if _, ok := seen[name]; !ok {
					seen[name] = struct{}{}
					allExecutables = append(allExecutables, name)
				}
			}
		}
	})
	return allExecutables
}

// findConcatenatedExecutable checks whether scriptPart looks like two
// executables concatenated together (e.g. "gitcommit" → "git" + "commit").
// Returns the first executable prefix found, or "".
func findConcatenatedExecutable(scriptPart string) string {
	executables := getAllExecutables()
	for _, exe := range executables {
		if len(exe) > 1 && strings.HasPrefix(scriptPart, exe) && len(scriptPart) > len(exe) {
			// The remainder after exe must also be a valid executable.
			remainder := scriptPart[len(exe):]
			for _, exe2 := range executables {
				if exe2 == remainder {
					return exe
				}
			}
		}
	}
	return ""
}

func init() {
	register(Rule{
		Name:     "missing_space_before_subcommand",
		Priority: 4000,
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			if len(parts) == 0 {
				return false
			}
			scriptPart := parts[0]
			// The first token must not itself be a valid executable.
			executables := getAllExecutables()
			for _, exe := range executables {
				if exe == scriptPart {
					return false
				}
			}
			return findConcatenatedExecutable(scriptPart) != ""
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			parts := cmd.ScriptParts()
			if len(parts) == 0 {
				return nil
			}
			exe := findConcatenatedExecutable(parts[0])
			if exe == "" {
				return nil
			}
			return single(strings.Replace(cmd.Script, exe, exe+" ", 1))
		},
	})
}
