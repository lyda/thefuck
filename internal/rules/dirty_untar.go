package rules

import (
	"archive/tar"
	"compress/bzip2"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

var tarExtensions = []string{
	".tar", ".tar.Z", ".tar.bz2", ".tar.gz", ".tar.lz",
	".tar.lzma", ".tar.xz", ".taz", ".tb2", ".tbz", ".tbz2",
	".tgz", ".tlz", ".txz", ".tz",
}

func isTarExtract(script string) bool {
	if strings.Contains(script, "--extract") {
		return true
	}
	parts := strings.Fields(script)
	return len(parts) > 1 && strings.Contains(parts[1], "x")
}

// tarFileAndDir finds the tar archive in script parts and returns (filename, dirname).
func tarFileAndDir(parts []string) (string, string) {
	for _, c := range parts {
		for _, ext := range tarExtensions {
			if strings.HasSuffix(c, ext) {
				return c, c[:len(c)-len(ext)]
			}
		}
	}
	return "", ""
}

// isDirtyTar opens the tar archive and checks whether all entries share a
// single top-level directory prefix. Returns true when the archive would
// splatter files into the current directory (i.e. is "dirty").
func isDirtyTar(filename string) bool {
	f, err := os.Open(filename) // #nosec G304 -- filename is from command args
	if err != nil {
		return false
	}
	defer f.Close()

	var tr *tar.Reader
	switch {
	case strings.HasSuffix(filename, ".gz") || strings.HasSuffix(filename, ".tgz"):
		gr, err := gzip.NewReader(f)
		if err != nil {
			return false
		}
		defer gr.Close()
		tr = tar.NewReader(gr)
	case strings.HasSuffix(filename, ".bz2") || strings.HasSuffix(filename, ".tbz") ||
		strings.HasSuffix(filename, ".tbz2") || strings.HasSuffix(filename, ".tb2"):
		tr = tar.NewReader(bzip2.NewReader(f))
	default:
		tr = tar.NewReader(f)
	}

	var topDir string
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return false
		}
		name := filepath.ToSlash(hdr.Name)
		// Get the top-level component.
		parts := strings.SplitN(name, "/", 2)
		top := parts[0]
		if top == "" || top == "." {
			continue
		}
		if topDir == "" {
			topDir = top
		} else if topDir != top {
			// Multiple top-level entries — dirty.
			return true
		}
	}
	return false
}

func init() {
	register(Rule{
		Name: "dirty_untar",
		Match: func(cmd types.Command) bool {
			if strings.Contains(cmd.Script, "-C") {
				return false
			}
			if !isTarExtract(cmd.Script) {
				return false
			}
			filename, dir := tarFileAndDir(cmd.ScriptParts())
			if filename == "" || dir == "" {
				return false
			}
			return isDirtyTar(filename)
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			_, dir := tarFileAndDir(cmd.ScriptParts())
			if dir == "" {
				return nil
			}
			return single(shellAnd("mkdir -p "+dir, cmd.Script+" -C "+dir))
		},
	})
}
