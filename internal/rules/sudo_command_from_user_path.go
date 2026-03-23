package rules

import (
	"os/exec"
	"regexp"

	"github.com/lyda/thefuck/internal/types"
)

var sudoCmdNotFoundRe = regexp.MustCompile(`sudo: (.*): command not found`)

func getSudoCommandName(output string) string {
	m := sudoCmdNotFoundRe.FindStringSubmatch(output)
	if m != nil {
		return m[1]
	}
	return ""
}

func init() {
	register(Rule{
		Name: "sudo_command_from_user_path",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			if len(parts) == 0 || parts[0] != "sudo" {
				return false
			}
			if !sudoCmdNotFoundRe.MatchString(cmd.Output) {
				return false
			}
			cmdName := getSudoCommandName(cmd.Output)
			if cmdName == "" {
				return false
			}
			_, err := exec.LookPath(cmdName)
			return err == nil
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			cmdName := getSudoCommandName(cmd.Output)
			if cmdName == "" {
				return nil
			}
			replacement := `env "PATH=$PATH" ` + cmdName
			return single(replaceArgument(cmd.Script, cmdName, replacement))
		},
	})
}
