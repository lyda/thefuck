package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/lyda/thefuck/internal/corrector"
	"github.com/lyda/thefuck/internal/runner"
	"github.com/lyda/thefuck/internal/shells"
	"github.com/lyda/thefuck/internal/types"
	"github.com/lyda/thefuck/internal/ui"

	// Import rules package for its init() side-effects (populates rules.All).
	_ "github.com/lyda/thefuck/internal/rules"
)

func main() {
	// Guess the shell with no args.
	if len(os.Args) == 1 || (len(os.Args) == 2 && os.Args[1] == "--alias") {
		if strings.HasSuffix(os.Getenv("SHELL"), "/zsh") {
			printInit("zsh")
		} else if strings.HasSuffix(os.Getenv("SHELL"), "/bash") {
			printInit("bash")
		} else if strings.HasSuffix(os.Getenv("SHELL"), "/fish") {
			printInit("fish")
		} else if strings.HasSuffix(os.Getenv("SHELL"), "/tcsh") ||
			strings.HasSuffix(os.Getenv("SHELL"), "/csh") {
			printInit("tcsh")
		}
		return
	}

	// Subcommand: thefuck init bash|zsh|fish
	if len(os.Args) == 3 && os.Args[1] == "init" {
		printInit(os.Args[2])
		return
	}

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: thefuck <command>")
		fmt.Fprintln(os.Stderr, "       thefuck init bash|zsh|fish")
		os.Exit(1)
	}

	script := os.Args[1]
	output, _ := runner.GetOutput(script)

	cmd := types.Command{Script: script, Output: output}
	suggestions := corrector.GetCorrectedCommands(cmd)

	selected := ui.SelectCommand(suggestions)
	if selected == "" {
		os.Exit(1)
	}
	// Print the corrected command to stdout so the shell alias can eval it.
	fmt.Println(selected)
}

func printInit(shellName string) {
	var sh shells.Shell
	switch shellName {
	case "bash":
		sh = shells.Bash{}
	case "zsh":
		sh = shells.Zsh{}
	case "fish":
		sh = shells.Fish{}
	case "tcsh", "csh":
		sh = shells.Tcsh{}
	case "powershell", "pwsh":
		sh = shells.PowerShell{}
	default:
		fmt.Fprintf(os.Stderr, "thefuck: unknown shell %q\n", shellName)
		fmt.Fprintln(os.Stderr, "supported: bash, zsh, fish, tcsh, powershell")
		os.Exit(1)
	}
	fmt.Println(sh.InitScript())
}
