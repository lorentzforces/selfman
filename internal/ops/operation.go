// The ops package contains the basic operations used by selfman.
// When the user invokes a command, that command emits a list of operations to be performed.
package ops

import "strings"

// An operation to be performed by selfman.
// These will typically be mutating operations to change selfman's managed state.
//
// TODO: still on the fence whether non-mutating, purely-informational things (like "list
// configured apps") should be implemented as Operations or if that's just a thing a command
// does at the end of its execution.
type Operation interface {
	// Execute the operation. If an error is returned, the operation has failed, and any context
	// should be included in the error itself. If err is non-nil, then msg should contain no useful
	// information and should be disregarded.
	Execute() (msg string, err error)

	// A human-readable description of what the operation will do when executed. Should include
	// context such as file names, destinations, etc.
	Describe() OpDescription
}

type OpDescription struct {
	TopLine string
	ContextLines []string
}

func (self OpDescription) String() string {
	var buf strings.Builder
	self.buildString(&buf)
	return buf.String()
}

const opDescriptionIndent string = "    "

func (self OpDescription) buildString(buf *strings.Builder) {
	buf.WriteString(self.TopLine)
	for _, contextLine := range self.ContextLines {
		buf.WriteString("\n")
		buf.WriteString(opDescriptionIndent)
		buf.WriteString(contextLine)
	}
}
