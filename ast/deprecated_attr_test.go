package ast

import (
	"testing"
)

func TestDeprecatedAttr(t *testing.T) {
	nodes := map[string]Node{
		`0x7fec4b0ab9c0 <line:180:48, col:63> "This function is provided for compatibility reasons only.  Due to security concerns inherent in the design of tempnam(3), it is highly recommended that you use mkstemp(3) instead." ""`: &DeprecatedAttr{
			Addr:        0x7fec4b0ab9c0,
			Pos:         NewPositionFromString("line:180:48, col:63"),
			Message1:    "This function is provided for compatibility reasons only.  Due to security concerns inherent in the design of tempnam(3), it is highly recommended that you use mkstemp(3) instead.",
			Message2:    "",
			IsInherited: false,
			ChildNodes:  []Node{},
		},
		`0xb75d00 <line:1107:12> "This function or variable may be unsafe. Consider using _snwprintf_s instead. To disable deprecation, use _CRT_SECURE_NO_WARNINGS. See online help for details." ""`: &DeprecatedAttr{
			Addr:        0xb75d00,
			Pos:         NewPositionFromString("line:1107:12"),
			Message1:    "This function or variable may be unsafe. Consider using _snwprintf_s instead. To disable deprecation, use _CRT_SECURE_NO_WARNINGS. See online help for details.",
			Message2:    "",
			IsInherited: false,
			ChildNodes:  []Node{},
		},
		`0xb75d00 <line:1107:12> Inherited "This function or variable may be unsafe. Consider using _snwprintf_s instead. To disable deprecation, use _CRT_SECURE_NO_WARNINGS. See online help for details." ""`: &DeprecatedAttr{
			Addr:        0xb75d00,
			Pos:         NewPositionFromString("line:1107:12"),
			Message1:    "This function or variable may be unsafe. Consider using _snwprintf_s instead. To disable deprecation, use _CRT_SECURE_NO_WARNINGS. See online help for details.",
			Message2:    "",
			IsInherited: true,
			ChildNodes:  []Node{},
		},
	}

	runNodeTests(t, nodes)
}
