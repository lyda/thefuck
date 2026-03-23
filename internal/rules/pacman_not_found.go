package rules

import (
	"os/exec"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

// pacmanSearchPackages runs `pacman -Ssq <name>` and returns matching package names.
func pacmanSearchPackages(name string) []string {
	out, err := exec.Command("pacman", "-Ssq", name).Output() // #nosec G204
	if err != nil {
		return nil
	}
	var pkgs []string
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if s := strings.TrimSpace(line); s != "" {
			pkgs = append(pkgs, s)
		}
	}
	return pkgs
}

func init() {
	register(Rule{
		Name: "pacman_not_found",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			if len(parts) == 0 {
				return false
			}
			firstTwo := parts[0]
			if len(parts) >= 2 {
				firstTwo = parts[0] + " " + parts[1]
			}
			isPacmanCmd := parts[0] == "yay" ||
				parts[0] == "pikaur" ||
				parts[0] == "yaourt" ||
				parts[0] == "pacman" ||
				firstTwo == "sudo pacman"
			if !isPacmanCmd {
				return false
			}
			return strings.Contains(cmd.Output, "error: target not found:")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			parts := cmd.ScriptParts()
			if len(parts) == 0 {
				return nil
			}
			// The package name is the last argument
			pkgName := parts[len(parts)-1]
			candidates := pacmanSearchPackages(pkgName)
			closest := getCloseMatches(pkgName, candidates, 0.6)
			if len(closest) == 0 {
				return nil
			}
			var results []string
			for _, c := range closest {
				results = append(results, replaceArgument(cmd.Script, pkgName, c))
			}
			return multi(results)
		},
	})
}
