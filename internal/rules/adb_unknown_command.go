package rules

import (
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

var adbCommands = []string{
	"backup",
	"bugreport",
	"connect",
	"devices",
	"disable-verity",
	"disconnect",
	"enable-verity",
	"emu",
	"forward",
	"get-devpath",
	"get-serialno",
	"get-state",
	"install",
	"install-multiple",
	"jdwp",
	"keygen",
	"kill-server",
	"logcat",
	"pull",
	"push",
	"reboot",
	"reconnect",
	"restore",
	"reverse",
	"root",
	"run-as",
	"shell",
	"sideload",
	"start-server",
	"sync",
	"tcpip",
	"uninstall",
	"unroot",
	"usb",
	"wait-for",
}

func init() {
	register(Rule{
		Name: "adb_unknown_command",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			return len(parts) >= 1 && parts[0] == "adb" &&
				strings.HasPrefix(cmd.Output, "Android Debug Bridge version")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			parts := cmd.ScriptParts()
			// Skip past adb flags (-s <serial>, -H <host>, -P <port>, -L <socket>,
			// or other single-letter flags like -d/-e/-a) to find the subcommand.
			for idx := 1; idx < len(parts); idx++ {
				arg := parts[idx]
				if strings.HasPrefix(arg, "-") {
					// Flags that consume the next argument
					if arg == "-s" || arg == "-H" || arg == "-P" || arg == "-L" {
						idx++ // skip the value
					}
					continue
				}
				// Check the previous token is not a flag that consumed this as its value
				prev := ""
				if idx > 0 {
					prev = parts[idx-1]
				}
				if prev == "-s" || prev == "-H" || prev == "-P" || prev == "-L" {
					continue
				}
				// This is the subcommand
				closest := getCloseMatches(arg, adbCommands, 0.6)
				if len(closest) == 0 {
					return nil
				}
				return single(replaceArgument(cmd.Script, arg, closest[0]))
			}
			return nil
		},
	})
}
