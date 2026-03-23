package rules

import (
	"fmt"
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

func init() {
	register(Rule{
		Name: "docker_image_being_used_by_container",
		Match: func(cmd types.Command) bool {
			parts := cmd.ScriptParts()
			if len(parts) < 1 || parts[0] != "docker" {
				return false
			}
			return strings.Contains(cmd.Output, "image is being used by running container")
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			fields := strings.Fields(strings.TrimSpace(cmd.Output))
			if len(fields) == 0 {
				return nil
			}
			containerID := fields[len(fields)-1]
			return single(shellAnd(fmt.Sprintf("docker container rm -f %s", containerID), cmd.Script))
		},
	})
}
