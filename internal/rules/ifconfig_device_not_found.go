package rules

import (
	"bufio"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

// getNetworkInterfaces returns a list of network interface names.
// On Linux it reads /proc/net/dev; elsewhere it runs `ifconfig -l`.
func getNetworkInterfaces() []string {
	if runtime.GOOS == "linux" {
		f, err := os.Open("/proc/net/dev")
		if err == nil {
			defer f.Close()
			var ifaces []string
			scanner := bufio.NewScanner(f)
			// Skip the two header lines
			scanner.Scan()
			scanner.Scan()
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				if line == "" {
					continue
				}
				// Format: "eth0:  ..."
				iface := strings.SplitN(line, ":", 2)[0]
				ifaces = append(ifaces, strings.TrimSpace(iface))
			}
			if len(ifaces) > 0 {
				return ifaces
			}
		}
	}

	// Fallback: use `ifconfig -a` and collect lines that don't start with space/tab.
	out, err := exec.Command("ifconfig", "-a").Output() // #nosec G204
	if err != nil {
		return nil
	}
	var ifaces []string
	for _, line := range strings.Split(string(out), "\n") {
		if line == "" || line == "\n" {
			continue
		}
		if !strings.HasPrefix(line, " ") && !strings.HasPrefix(line, "\t") {
			iface := strings.Fields(line)[0]
			// Strip trailing colon (Linux ifconfig format)
			iface = strings.TrimSuffix(iface, ":")
			if iface != "" {
				ifaces = append(ifaces, iface)
			}
		}
	}
	return ifaces
}

func init() {
	register(Rule{
		Name: "ifconfig_device_not_found",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			if len(parts) == 0 || parts[0] != "ifconfig" {
				return false
			}
			if _, err := exec.LookPath("ifconfig"); err != nil {
				return false
			}
			output := strings.ToLower(cmd.Output)
			return strings.Contains(output, "error fetching interface information: device not found") ||
				strings.Contains(output, "no such interface") ||
				strings.Contains(output, "does not exist")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			// The Python rule grabs the first word of the output minus trailing colon.
			// The interface name is typically the first word of the command (after ifconfig).
			parts := cmd.ScriptParts()
			if len(parts) < 2 {
				return nil
			}
			wrongIface := parts[1]
			ifaces := getNetworkInterfaces()
			closest := getCloseMatches(wrongIface, ifaces, 0.6)
			scripts := make([]string, 0, len(closest))
			for _, iface := range closest {
				scripts = append(scripts, replaceArgument(cmd.Script, wrongIface, iface))
			}
			if len(scripts) == 0 {
				return nil
			}
			return multi(scripts)
		},
	})
}
