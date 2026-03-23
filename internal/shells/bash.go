package shells

import "strings"

// Bash implements Shell for bash (and sh-compatible shells).
type Bash struct{}

func (b Bash) And(commands ...string) string {
	return strings.Join(commands, " && ")
}

func (b Bash) Or(commands ...string) string {
	return strings.Join(commands, " || ")
}

func (b Bash) Quote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}

func (b Bash) InitScript() string {
	return `fuck() {
    TF_CMD=$(TF_SHELL=bash thefuck "$(fc -ln -1 | sed 's/^ *//;s/ *$//')") && eval "$TF_CMD"
}`
}
