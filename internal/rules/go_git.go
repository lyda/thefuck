package rules

import (
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

// gitSubcommands is the set of common `git` subcommands (including aliases like co).
var gitSubcommands = map[string]bool{
	"add": true, "am": true, "archive": true, "bisect": true, "blame": true,
	"branch": true, "checkout": true, "cherry-pick": true, "clean": true,
	"clone": true, "co": true, "commit": true, "describe": true, "diff": true,
	"fetch": true, "gc": true, "grep": true, "init": true, "log": true,
	"merge": true, "mv": true, "pull": true, "push": true, "rebase": true,
	"remote": true, "reset": true, "restore": true, "revert": true, "rm": true,
	"show": true, "stash": true, "status": true, "submodule": true,
	"switch": true, "tag": true, "worktree": true,
}

func init() {
	register(Rule{
		Name: "go_git",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			return len(parts) >= 2 && parts[0] == "go" &&
				strings.Contains(cmd.Output, "unknown command") &&
				gitSubcommands[parts[1]]
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			return single(strings.Replace(cmd.Script, "go ", "git ", 1))
		},
	})
}
