package rules

import (
	"slices"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

var sudoPatterns = []string{
	"permission denied",
	"eacces",
	"pkg: insufficient privileges",
	"you cannot perform this operation unless you are root",
	"non-root users cannot",
	"operation not permitted",
	"not super-user",
	"superuser privilege",
	"root privilege",
	"this command has to be run under the root user.",
	"this operation requires root.",
	"requested operation requires superuser privilege",
	"must be run as root",
	"must run as root",
	"must be superuser",
	"must be root",
	"need to be root",
	"need root",
	"needs to be run as root",
	"only root can ",
	"you don't have access to the history db.",
	"authentication is required",
	"edspermissionerror",
	"you don't have write permissions",
	"use `sudo`",
	"sudorequirederror",
	"error: insufficient privileges",
	"updatedb: can not open a temporary file",
}

func sudoMatch(cmd types.Command) bool {
	parts := cmd.ScriptParts()
	if len(parts) > 0 && parts[0] == "sudo" {
		// Already sudo — but only skip if no && (compound command)
		hasAnd := slices.Contains(parts, "&&")
		if !hasAnd {
			return false
		}
	}
	lower := strings.ToLower(cmd.Output)
	for _, p := range sudoPatterns {
		if strings.Contains(lower, p) {
			return true
		}
	}
	return false
}

func sudoGetNewCommand(cmd types.Command) []types.CorrectedCommand {
	if strings.Contains(cmd.Script, "&&") {
		parts := cmd.ScriptParts()
		filtered := make([]string, 0, len(parts))
		for _, p := range parts {
			if p != "sudo" {
				filtered = append(filtered, p)
			}
		}
		return single(`sudo sh -c "` + strings.Join(filtered, " ") + `"`)
	}
	if strings.Contains(cmd.Script, ">") {
		escaped := strings.ReplaceAll(cmd.Script, `"`, `\"`)
		return single(`sudo sh -c "` + escaped + `"`)
	}
	return single("sudo " + cmd.Script)
}

func init() {
	register(Rule{
		Name:          "sudo",
		Match:         sudoMatch,
		GetNewCommand: sudoGetNewCommand,
	})
}
