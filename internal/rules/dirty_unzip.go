package rules

import (
	"archive/zip"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

// zipFile returns the zip archive path from the unzip command's arguments.
// unzip usage: unzip [-flags] file[.zip] [file(s)...] [-x file(s)...]
func zipFile(parts []string) string {
	for _, c := range parts[1:] {
		if strings.HasPrefix(c, "-") {
			continue
		}
		if strings.HasSuffix(c, ".zip") {
			return c
		}
		return c + ".zip"
	}
	return ""
}

// isBadZip returns true when the zip archive would splatter files into the
// current directory (more than one top-level entry, or entries not grouped
// under a single directory).
func isBadZip(filename string) bool {
	r, err := zip.OpenReader(filename) // #nosec G304 -- filename is from command args
	if err != nil {
		return false
	}
	defer r.Close()

	var topDir string
	for _, f := range r.File {
		name := f.Name
		// Get top-level component.
		parts := strings.SplitN(name, "/", 2)
		top := parts[0]
		if top == "" || top == "." {
			continue
		}
		if topDir == "" {
			topDir = top
		} else if topDir != top {
			return true
		}
	}
	return false
}

func init() {
	register(Rule{
		Name: "dirty_unzip",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			if len(parts) == 0 || parts[0] != "unzip" {
				return false
			}
			if strings.Contains(cmd.Script, " -d ") {
				return false
			}
			zf := zipFile(parts)
			if zf == "" {
				return false
			}
			return isBadZip(zf)
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			parts := cmd.ScriptParts()
			zf := zipFile(parts)
			if zf == "" {
				return nil
			}
			// Strip .zip suffix to get the directory name.
			dir := zf[:len(zf)-4]
			return single(cmd.Script + " -d " + dir)
		},
	})
}
