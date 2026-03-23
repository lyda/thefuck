package rules

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

// getVirtualenvs returns directory names inside WORKON_HOME (default ~/.virtualenvs).
func getVirtualenvs() []string {
	workonHome := os.Getenv("WORKON_HOME")
	if workonHome == "" {
		home := os.Getenv("HOME")
		workonHome = filepath.Join(home, ".virtualenvs")
	}

	entries, err := os.ReadDir(workonHome)
	if err != nil {
		return nil
	}

	var envs []string
	for _, entry := range entries {
		if entry.IsDir() {
			envs = append(envs, entry.Name())
		}
	}
	return envs
}

func init() {
	register(Rule{
		Name: "workon_doesnt_exists",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			if len(parts) < 2 || parts[0] != "workon" {
				return false
			}
			envs := getVirtualenvs()
			misspelled := parts[1]
			for _, env := range envs {
				if env == misspelled {
					return false
				}
			}
			return true
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			parts := cmd.ScriptParts()
			if len(parts) < 2 {
				return nil
			}
			misspelled := parts[1]
			envs := getVirtualenvs()
			createNew := "mkvirtualenv " + misspelled

			if len(envs) == 0 {
				return single(createNew)
			}

			matches := getCloseMatches(misspelled, envs, 0.6)
			var results []string
			for _, m := range matches {
				results = append(results, strings.Replace(cmd.Script, misspelled, m, 1))
			}
			results = append(results, createNew)
			return multi(results)
		},
	})
}
