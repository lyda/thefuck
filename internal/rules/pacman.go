package rules

import (
	"os/exec"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

// pacmanPkgManagerCmd returns the pacman-compatible package manager to use
// (yay, pikaur, yaourt, or "sudo pacman"), or "" if none is found.
func pacmanPkgManagerCmd() string {
	for _, mgr := range []string{"yay", "pikaur", "yaourt"} {
		if _, err := exec.LookPath(mgr); err == nil {
			return mgr
		}
	}
	if _, err := exec.LookPath("pacman"); err == nil {
		return "sudo pacman"
	}
	return ""
}

// pkgfilePackages runs `pkgfile -b -v <cmdName>` and returns the list of
// packages that provide it (the first field of each output line).
func pkgfilePackages(cmdName string) []string {
	out, err := exec.Command("pkgfile", "-b", "-v", cmdName).Output() // #nosec G204
	if err != nil {
		return nil
	}
	var pkgs []string
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) > 0 {
			pkgs = append(pkgs, fields[0])
		}
	}
	return pkgs
}

// pacmanCmdName extracts the bare command name from the script, stripping a
// leading "sudo " if present.
func pacmanCmdName(script string) string {
	s := strings.TrimPrefix(strings.TrimSpace(script), "sudo ")
	return strings.Fields(s)[0]
}

func init() {
	register(Rule{
		Name: "pacman",
		Match: func(cmd types.Command) bool {
			if _, err := exec.LookPath("pkgfile"); err != nil {
				return false
			}
			if !strings.Contains(cmd.Output, "not found") {
				return false
			}
			cmdName := pacmanCmdName(cmd.Script)
			if cmdName == "" {
				return false
			}
			return len(pkgfilePackages(cmdName)) > 0
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			mgr := pacmanPkgManagerCmd()
			if mgr == "" {
				return nil
			}
			cmdName := pacmanCmdName(cmd.Script)
			packages := pkgfilePackages(cmdName)
			if len(packages) == 0 {
				return nil
			}
			var results []string
			for _, pkg := range packages {
				install := mgr + " -S " + pkg
				results = append(results, shellAnd(install, cmd.Script))
			}
			return multi(results)
		},
	})
}
