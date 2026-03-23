package types

import "strings"

// DefaultPriority is the priority assigned to rules that don't specify one.
// Lower priority values are preferred.
const DefaultPriority = 1000

// Command represents a failed shell command with its output.
type Command struct {
	Script string
	Output string
}

// ScriptParts splits Script into whitespace-delimited tokens.
func (c Command) ScriptParts() []string {
	return strings.Fields(c.Script)
}

// CorrectedCommand is a proposed fix for a Command.
type CorrectedCommand struct {
	Script   string
	Priority int
}
