package rules

import (
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

// goSubcommands is the set of valid top-level `go` subcommands.
var goSubcommands = map[string]bool{
	"bug": true, "build": true, "clean": true, "doc": true, "env": true,
	"fix": true, "fmt": true, "generate": true, "get": true, "install": true,
	"list": true, "mod": true, "run": true, "telemetry": true, "test": true,
	"tool": true, "version": true, "vet": true, "work": true,
}

func init() {
	register(Rule{
		Name: "git_go",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			return len(parts) >= 2 && parts[0] == "git" &&
				strings.Contains(cmd.Output, "is not a git command") &&
				goSubcommands[parts[1]]
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			return single(strings.Replace(cmd.Script, "git ", "go ", 1))
		},
	})
}
