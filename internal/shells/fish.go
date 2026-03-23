package shells

import "strings"

// Fish implements Shell for fish shell.
type Fish struct{}

func (f Fish) And(commands ...string) string {
	return strings.Join(commands, "; and ")
}

func (f Fish) Or(commands ...string) string {
	return strings.Join(commands, "; or ")
}

func (f Fish) Quote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `\'`) + "'"
}

func (f Fish) InitScript() string {
	return `function fuck
    set -l TF_CMD (TF_SHELL=fish thefuck (history | head -n1))
    if test $status -eq 0
        eval $TF_CMD
    end
end`
}
