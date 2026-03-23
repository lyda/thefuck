package ui

import (
	"fmt"
	"os"

	"golang.org/x/term"

	"github.com/lyda/thefuck/internal/types"
)

// SelectCommand displays the top suggestion on stderr and waits for a keypress.
// Enter confirms and returns the script; Ctrl-C aborts and returns "".
func SelectCommand(suggestions []types.CorrectedCommand) string {
	if len(suggestions) == 0 {
		fmt.Fprintln(os.Stderr, "No fucks given")
		return ""
	}
	top := suggestions[0]
	fmt.Fprintf(os.Stderr, "\x1b[1mthefuck\x1b[0m: %s [enter/ctrl-c] ", top.Script)

	fd := int(os.Stdin.Fd()) // #nosec G115 -- fd fits in int on all supported platforms
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		// Not a real TTY (e.g. piped/scripted) — auto-confirm.
		fmt.Fprintln(os.Stderr)
		return top.Script
	}
	defer term.Restore(fd, oldState) //nolint:errcheck

	buf := make([]byte, 1)
	for {
		_, err := os.Stdin.Read(buf)
		if err != nil {
			fmt.Fprint(os.Stderr, "\r\n")
			return ""
		}
		switch buf[0] {
		case '\r', '\n':
			fmt.Fprint(os.Stderr, "\r\n")
			return top.Script
		case 3: // Ctrl-C
			fmt.Fprint(os.Stderr, "\r\nAborted\r\n")
			return ""
		}
	}
}
