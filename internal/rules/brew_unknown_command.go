package rules

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/lyda/thefuck/internal/types"
)

var (
	brewCommandsOnce sync.Once
	brewCmds         []string
)

var brewUnknownCmdRe = regexp.MustCompile(`Error: Unknown command: ([a-z]+)`)

// getBrewRepository returns the path prefix from `brew --repository`.
func getBrewRepository() string {
	out, err := exec.Command("brew", "--repository").Output() // #nosec G204
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// getBrewCommands lists all commands available in a Homebrew installation.
func getBrewCommands() []string {
	brewCommandsOnce.Do(func() {
		prefix := getBrewRepository()
		if prefix == "" {
			// Fallback list (based on Homebrew 0.9.5)
			brewCmds = []string{
				"info", "home", "options", "install", "uninstall",
				"search", "list", "update", "upgrade", "pin", "unpin",
				"doctor", "create", "edit", "cask",
			}
			return
		}

		// Core commands: Library/Homebrew/cmd/*.rb and *.sh
		cmdDir := filepath.Join(prefix, "Library", "Homebrew", "cmd")
		entries, err := os.ReadDir(cmdDir)
		if err == nil {
			for _, e := range entries {
				name := e.Name()
				if strings.HasSuffix(name, ".rb") || strings.HasSuffix(name, ".sh") {
					ext := filepath.Ext(name)
					brewCmds = append(brewCmds, strings.TrimSuffix(name, ext))
				}
			}
		}

		// Tap-specific commands: Library/Taps/<user>/homebrew-*/cmd/brew-*.rb
		tapsDir := filepath.Join(prefix, "Library", "Taps")
		users, err := os.ReadDir(tapsDir)
		if err != nil {
			return
		}
		for _, user := range users {
			if !user.IsDir() {
				continue
			}
			userDir := filepath.Join(tapsDir, user.Name())
			taps, err := os.ReadDir(userDir)
			if err != nil {
				continue
			}
			for _, tap := range taps {
				if !tap.IsDir() || !strings.HasPrefix(tap.Name(), "homebrew-") {
					continue
				}
				tapCmdDir := filepath.Join(userDir, tap.Name(), "cmd")
				files, err := os.ReadDir(tapCmdDir)
				if err != nil {
					continue
				}
				for _, f := range files {
					name := f.Name()
					if strings.HasPrefix(name, "brew-") && strings.HasSuffix(name, ".rb") {
						cmd := strings.TrimSuffix(strings.TrimPrefix(name, "brew-"), ".rb")
						brewCmds = append(brewCmds, cmd)
					}
				}
			}
		}

		if len(brewCmds) == 0 {
			brewCmds = []string{
				"info", "home", "options", "install", "uninstall",
				"search", "list", "update", "upgrade", "pin", "unpin",
				"doctor", "create", "edit", "cask",
			}
		}
	})
	return brewCmds
}

func init() {
	register(Rule{
		Name: "brew_unknown_command",
		Match: func(cmd types.Command) bool {
			if _, err := exec.LookPath("brew"); err != nil {
				return false
			}
			if !strings.Contains(cmd.Script, "brew") ||
				!strings.Contains(cmd.Output, "Unknown command") {
				return false
			}
			m := brewUnknownCmdRe.FindStringSubmatch(cmd.Output)
			if len(m) < 2 {
				return false
			}
			return len(getCloseMatches(m[1], getBrewCommands(), 0.6)) > 0
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			m := brewUnknownCmdRe.FindStringSubmatch(cmd.Output)
			if len(m) < 2 {
				return nil
			}
			brokenCmd := m[1]
			closest := getCloseMatches(brokenCmd, getBrewCommands(), 0.6)
			scripts := make([]string, 0, len(closest))
			for _, c := range closest {
				scripts = append(scripts, replaceArgument(cmd.Script, brokenCmd, c))
			}
			if len(scripts) == 0 {
				return nil
			}
			return multi(scripts)
		},
	})
}
