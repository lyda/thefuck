package shells

import "strings"

// Zsh implements Shell for zsh.
type Zsh struct{}

func (z Zsh) And(commands ...string) string {
	return strings.Join(commands, " && ")
}

func (z Zsh) Or(commands ...string) string {
	return strings.Join(commands, " || ")
}

func (z Zsh) Quote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}

func (z Zsh) InitScript() string {
	return `fuck() {
    TF_CMD=$(TF_SHELL=zsh thefuck "$(fc -ln -1)") && eval "$TF_CMD"
}`
}
