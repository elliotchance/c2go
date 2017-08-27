package ast

import (
	"testing"
)

func TestDeprecatedAttr(t *testing.T) {
	nodes := map[string]Node{
		`0x7fec4b0ab9c0 <line:180:48, col:63> "This function is provided for compatibility reasons only.  Due to security concerns inherent in the design of tempnam(3), it is highly recommended that you use mkstemp(3) instead." ""`: &DeprecatedAttr{
			Addr:     0x7fec4b0ab9c0,
			Position: "line:180:48, col:63",
			Message1: "This function is provided for compatibility reasons only.  Due to security concerns inherent in the design of tempnam(3), it is highly recommended that you use mkstemp(3) instead.",
			Message2: "",
			Children: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
