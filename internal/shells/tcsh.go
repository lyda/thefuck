package shells

import "strings"

// Tcsh implements Shell for tcsh/csh.
type Tcsh struct{}

func (Tcsh) And(commands ...string) string {
	return strings.Join(commands, " && ")
}

func (Tcsh) Or(commands ...string) string {
	return strings.Join(commands, " || ")
}

func (Tcsh) Quote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'"'"'`) + "'"
}

func (Tcsh) InitScript() string {
	return `alias fuck 'setenv TF_SHELL tcsh && set fucked_cmd=` + "`" + `history -h 2 | head -n 1` + "`" + ` && eval ` + "`" + `thefuck "$fucked_cmd"` + "`" + `'`
}
