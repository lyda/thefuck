package shells

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Shell abstracts shell-specific syntax differences.
type Shell interface {
	// And joins commands with the shell AND operator (e.g. " && " for bash).
	And(commands ...string) string
	// Or joins commands with the shell OR operator (e.g. " || " for bash).
	Or(commands ...string) string
	// Quote returns a safely single-quoted version of s.
	Quote(s string) string
	// InitScript returns the shell function definition for the "fuck" alias.
	InitScript() string
}

var (
	detected   Shell
	detectOnce sync.Once
)

// Detect returns the Shell for the current environment.
// Reads TF_SHELL first (set by the alias), then $SHELL basename.
// Defaults to Bash.
func Detect() Shell {
	detectOnce.Do(func() {
		name := os.Getenv("TF_SHELL")
		if name == "" {
			name = filepath.Base(os.Getenv("SHELL"))
		}
		switch strings.ToLower(name) {
		case "zsh":
			detected = Zsh{}
		case "fish":
			detected = Fish{}
		case "tcsh", "csh":
			detected = Tcsh{}
		case "powershell", "pwsh", "powershell.exe":
			detected = PowerShell{}
		default:
			detected = Bash{}
		}
	})
	return detected
}
