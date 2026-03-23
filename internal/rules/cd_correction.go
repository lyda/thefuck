package rules

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

const cdCorrectionCutoff = 0.6

// getSubDirs returns immediate child directories of parent.
func getSubDirs(parent string) []string {
	entries, err := os.ReadDir(parent)
	if err != nil {
		return nil
	}
	var dirs []string
	for _, e := range entries {
		if e.IsDir() {
			dirs = append(dirs, e.Name())
		}
	}
	return dirs
}

func init() {
	register(Rule{
		Name: "cd_correction",
		Match: func(cmd types.Command) bool {
			if !strings.HasPrefix(cmd.Script, "cd ") {
				return false
			}
			lower := strings.ToLower(cmd.Output)
			return strings.Contains(lower, "no such file or directory") ||
				strings.Contains(lower, "cd: can't cd to") ||
				strings.Contains(lower, "does not exist")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			parts := cmd.ScriptParts()
			if len(parts) < 2 {
				return nil
			}

			// Split the destination path into components.
			dest := filepath.ToSlash(parts[1])
			components := strings.Split(dest, "/")
			// Drop trailing empty component from trailing slash.
			if len(components) > 0 && components[len(components)-1] == "" {
				components = components[:len(components)-1]
			}

			var cwd string
			startIdx := 0

			if len(components) > 0 && components[0] == "" {
				// Absolute path.
				cwd = string(os.PathSeparator)
				startIdx = 1
			} else {
				var err error
				cwd, err = os.Getwd()
				if err != nil {
					return nil
				}
			}

			for _, dir := range components[startIdx:] {
				switch dir {
				case ".":
					// Stay in cwd.
				case "..":
					cwd = filepath.Dir(cwd)
				default:
					matches := getCloseMatches(dir, getSubDirs(cwd), cdCorrectionCutoff)
					if len(matches) == 0 {
						// Fall back to cd_mkdir behaviour: mkdir -p + cd.
						return single(shellAnd("mkdir -p "+parts[1], "cd "+parts[1]))
					}
					cwd = filepath.Join(cwd, matches[0])
				}
			}
			return single(`cd "` + cwd + `"`)
		},
	})
}
