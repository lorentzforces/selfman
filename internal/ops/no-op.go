package ops

import "fmt"

type NoOp struct {
	TypeOfNoOp string
	Description string
}

func (self NoOp) Execute() (string, error) {
	return "Successfully did nothing", nil
}

func (self NoOp) Describe() OpDescription {
	return OpDescription {
		TopLine: fmt.Sprintf("No-op operation: %s", self.TypeOfNoOp),
		ContextLines: []string{
			self.Description,
		},
	}
}

var SkipCloneOp = NoOp{
	TypeOfNoOp: "clone",
	Description: "Skipping git clone, source already present",
}

var NoBuildOp = NoOp{
	TypeOfNoOp: "build",
	Description: "This application does not need to be built",
}

var NoUpdateOp = NoOp{
	TypeOfNoOp: "update",
	Description: "This application does not update",
}
