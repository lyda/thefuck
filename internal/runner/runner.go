package runner

import (
	"context"
	"os"
	"os/exec"
	"strings"
	"time"
)

// DefaultTimeout is how long to wait for a command before giving up.
const DefaultTimeout = 10 * time.Second

// GetOutput re-runs script in a shell and returns combined stdout+stderr.
// Returns an empty string on timeout; non-zero exit codes are not errors
// (that's why we're correcting the command).
func GetOutput(script string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "/bin/sh", "-c", script) // #nosec G204 -- intentional: re-running the user's failed command
	cmd.Env = buildEnv()
	cmd.Stdin = nil

	out, _ := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return "", ctx.Err()
	}
	return string(out), nil
}

// buildEnv returns os.Environ with LC_ALL, LANG, and GIT_TRACE overridden.
// GIT_TRACE=1 causes git to emit alias expansion info in its output,
// which the git rules rely on.
func buildEnv() []string {
	overrides := map[string]string{
		"LC_ALL":    "C",
		"LANG":      "C",
		"GIT_TRACE": "1",
	}
	result := make([]string, 0, len(os.Environ())+len(overrides))
	for _, e := range os.Environ() {
		before, _, ok := strings.Cut(e, "=")
		if !ok {
			result = append(result, e)
			continue
		}
		if _, skip := overrides[before]; !skip {
			result = append(result, e)
		}
	}
	for k, v := range overrides {
		result = append(result, k+"="+v)
	}
	return result
}

// ShellBin returns the path to /bin/sh or falls back to "sh".
func ShellBin() string {
	if _, err := exec.LookPath("/bin/sh"); err == nil {
		return "/bin/sh"
	}
	return "sh"
}
