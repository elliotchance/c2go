package ast

import (
	"testing"
)

func TestAvailabilityAttr(t *testing.T) {
	nodes := map[string]Node{
		`0x7fc5ff8e5d18 </usr/include/AvailabilityInternal.h:21697:88, col:124> macos 10.10 0 0 "" ""`: &AvailabilityAttr{
			Addr:          0x7fc5ff8e5d18,
			Pos:           "/usr/include/AvailabilityInternal.h:21697:88, col:124",
			OS:            "macos",
			Version:       "10.10",
			Unknown1:      0,
			Unknown2:      0,
			IsUnavailable: false,
			Message1:      "",
			Message2:      "",
			IsInherited:   false,
			ChildNodes:    []Node{},
		},
		`0x7fc5ff8e60d0 </usr/include/Availability.h:215:81, col:115> watchos 3.0 0 0 "" ""`: &AvailabilityAttr{
			Addr:          0x7fc5ff8e60d0,
			Pos:           "/usr/include/Availability.h:215:81, col:115",
			OS:            "watchos",
			Version:       "3.0",
			Unknown1:      0,
			Unknown2:      0,
			IsUnavailable: false,
			Message1:      "",
			Message2:      "",
			IsInherited:   false,
			ChildNodes:    []Node{},
		},
		`0x7fc5ff8e6170 <col:81, col:115> tvos 10.0 0 0 "" ""`: &AvailabilityAttr{
			Addr:          0x7fc5ff8e6170,
			Pos:           "col:81, col:115",
			OS:            "tvos",
			Version:       "10.0",
			Unknown1:      0,
			Unknown2:      0,
			IsUnavailable: false,
			Message1:      "",
			Message2:      "",
			IsInherited:   false,
			ChildNodes:    []Node{},
		},
		`0x7fc5ff8e61d8 <col:81, col:115> ios 10.0 0 0 "" ""`: &AvailabilityAttr{
			Addr:          0x7fc5ff8e61d8,
			Pos:           "col:81, col:115",
			OS:            "ios",
			Version:       "10.0",
			Unknown1:      0,
			Unknown2:      0,
			IsUnavailable: false,
			Message1:      "",
			Message2:      "",
			IsInherited:   false,
			ChildNodes:    []Node{},
		},
		`0x7fc5ff8f0e18 </usr/include/sys/cdefs.h:275:50, col:99> swift 0 0 0 Unavailable "Use snprintf instead." ""`: &AvailabilityAttr{
			Addr:          0x7fc5ff8f0e18,
			Pos:           "/usr/include/sys/cdefs.h:275:50, col:99",
			OS:            "swift",
			Version:       "0",
			Unknown1:      0,
			Unknown2:      0,
			IsUnavailable: true,
			Message1:      "Use snprintf instead.",
			Message2:      "",
			IsInherited:   false,
			ChildNodes:    []Node{},
		},
		`0x7fc5ff8f1988 <line:275:50, col:99> swift 0 0 0 Unavailable "Use mkstemp(3) instead." ""`: &AvailabilityAttr{
			Addr:          0x7fc5ff8f1988,
			Pos:           "line:275:50, col:99",
			OS:            "swift",
			Version:       "0",
			Unknown1:      0,
			Unknown2:      0,
			IsUnavailable: true,
			Message1:      "Use mkstemp(3) instead.",
			Message2:      "",
			IsInherited:   false,
			ChildNodes:    []Node{},
		},
		`0x104035438 </usr/include/AvailabilityInternal.h:14571:88, col:124> macosx 10.10 0 0 ""`: &AvailabilityAttr{
			Addr:          0x104035438,
			Pos:           "/usr/include/AvailabilityInternal.h:14571:88, col:124",
			OS:            "macosx",
			Version:       "10.10",
			Unknown1:      0,
			Unknown2:      0,
			IsUnavailable: false,
			Message1:      "",
			Message2:      "",
			IsInherited:   false,
			ChildNodes:    []Node{},
		},
		`0x7f9bd588b1a8 </usr/include/gethostuuid.h:39:65, col:100> Inherited macos 10.5 0 0 "" ""`: &AvailabilityAttr{
			Addr:          0x7f9bd588b1a8,
			Pos:           "/usr/include/gethostuuid.h:39:65, col:100",
			OS:            "macos",
			Version:       "10.5",
			Unknown1:      0,
			Unknown2:      0,
			IsUnavailable: false,
			Message1:      "",
			Message2:      "",
			IsInherited:   true,
			ChildNodes:    []Node{},
		},
	}

	runNodeTests(t, nodes)
}
