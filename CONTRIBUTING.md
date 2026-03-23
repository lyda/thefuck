# Reporting issues

If you have any issue with *The Fuck*, check if it's already been reported and
if not, open an issue on [GitHub](https://github.com/lyda/thefuck) with:

  * your shell and its version (`bash`, `zsh`, `fish`, `tcsh`);
  * your OS and version;
  * how to reproduce the bug;
  * the output of the failing command and what *The Fuck* did (or didn't do);
  * anything else you think is relevant.

# Pull requests

Pull requests are welcome for new rules, bug fixes, and improvements.

# Developing

## Prerequisites

- Always uses the most recent version of Go.
- `staticcheck` and `gosec` are managed as Go tool dependencies — no separate install needed

## Getting started

```bash
git clone https://github.com/lyda/thefuck
cd thefuck
make
```

## Make targets

| Target          | Description                                         |
|-----------------|-----------------------------------------------------|
| `make build`    | Build `./bin/thefuck`                               |
| `make check`    | Run all checks: fmt, vet, staticcheck, gosec, tests |
| `make fix`      | Format code, update dependencies, run `go fix`      |
| `make fmt`      | Check formatting (fails if any file needs `gofmt`)  |
| `make vet`      | Run `go vet`                                        |
| `make lint`     | Run `staticcheck`                                   |
| `make sec`      | Run `gosec` security checks                         |
| `make test`     | Run tests with race detector and coverage           |
| `make coverage` | Generate `coverage.html` with per-file breakdown    |
| `make install`  | Install binary to `$GOPATH/bin`                     |
| `make clean`    | Remove build artifacts                              |

Before submitting a PR, run `make check` to make sure everything passes. If
`make fmt` fails, run `make fix` to auto-format and tidy the code.

## Adding a rule

Each rule lives in its own file in `internal/rules/`. A minimal rule looks like:

```go
package rules

import (
    "strings"

    "github.com/lyda/thefuck/internal/types"
)

func init() {
    register(Rule{
        Name: "my_rule",
        Match: func(cmd types.Command) bool {
            return strings.Contains(cmd.Output, "some error")
        },
        GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
            return single(strings.Replace(cmd.Script, "wrong", "right", 1))
        },
    })
}
```

### Available helpers

Defined in `internal/rules/rules.go`:

- `single(script string)` — return one correction
- `multi(scripts []string)` — return multiple corrections
- `shellAnd(cmds ...string)` — join with shell AND operator (`&&` for bash/zsh/tcsh, `; and ` for fish)
- `replaceArgument(script, from, to string)` — replace a shell word in a command string
- `getCloseMatches(word string, possibilities []string, cutoff float64)` — fuzzy match, cutoff 0.0–1.0
- `getAllMatchedCommands(output string, separators []string)` — extract suggestion lines after a separator

### `types.Command` fields

- `cmd.Script` — the original command string
- `cmd.Output` — combined stdout+stderr from re-running the command
- `cmd.ScriptParts()` — the command split into shell words

### Caching subprocess calls

If your rule runs a subprocess to discover valid commands (e.g. parsing help
text), cache the result with `sync.Once`:

```go
var (
    myCommandsOnce sync.Once
    myCommands     []string
)

func getMyCommands() []string {
    myCommandsOnce.Do(func() {
        out, _ := exec.Command("mytool", "help").CombinedOutput() // #nosec G204
        // parse out...
    })
    return myCommands
}
```

### Security annotations

`gosec` will flag certain patterns that are intentional in this codebase:

- `exec.Command` with a variable argument &rarr; add `// #nosec G204`
- `int(fd)` conversion from `uintptr` &rarr; add `// #nosec G115`
- `os.Open` with a variable path &rarr; add `// #nosec G304` (or `G304,G703` if taint analysis also fires)

### Rule priority

Lower priority value = presented first. Default is 1000. Set a custom priority
if your rule should rank above or below the default:

```go
register(Rule{
    Name:     "my_rule",
    Priority: 100, // show before default-priority rules
    ...
})
```
