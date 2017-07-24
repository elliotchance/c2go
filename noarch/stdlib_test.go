package noarch

import (
	"testing"
)

// web link with clarification:
// http://www.cplusplus.com/reference/cstdlib/strtol/
func TestStrtol(t *testing.T) {

	tests := []struct {
		inputBytes     []byte
		inputBase      int
		expectedValue  int32
		expectedString []byte
	}{
		{
			inputBytes:     []byte("2030300 This is test"),
			inputBase:      10,
			expectedValue:  2030300,
			expectedString: []byte(" This is test"),
		},
		{
			inputBytes:     []byte("2030300 This is test"),
			inputBase:      6,
			expectedValue:  97308,
			expectedString: []byte(" This is test"),
		},
		{
			inputBytes:     []byte("10011101011 This is test"),
			inputBase:      2,
			expectedValue:  1259,
			expectedString: []byte(" This is test"),
		},
		{
			inputBytes:     []byte("10011101011 This is test"),
			inputBase:      7,
			expectedValue:  283433599,
			expectedString: []byte(" This is test"),
		},
		{
			inputBytes:     []byte("2030300 This is test"),
			inputBase:      0,
			expectedValue:  2030300,
			expectedString: []byte(" This is test"),
		},
	}

	for _, tt := range tests {
		a := tt.inputBytes
		var b []byte
		c := tt.inputBase
		ret := Strtol(a, &b, c)
		if ret != tt.expectedValue {
			t.Errorf("Strtol() return %v. expected = %v", ret, tt.expectedValue)
		}
		if len(b) != len(tt.expectedString) {
			t.Errorf("Strtol() by length return %#v. expected = %#v", string(b), string(tt.expectedString))
		}
		for i := range b {
			if b[i] != tt.expectedString[i] {
				t.Errorf("Strtol() by body return %#v. expected = %#v", string(b), string(tt.expectedString))
			}
		}
	}
}
