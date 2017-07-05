package noarch

import (
	"testing"
)

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
			inputBytes:     []byte("2030300 This is test"),
			inputBase:      0,
			expectedValue:  2030300,
			expectedString: []byte(" This is test"),
		},
		{
			inputBytes:     []byte("a d 2030300 This is test"),
			inputBase:      10,
			expectedValue:  0,
			expectedString: []byte("a d 2030300 This is test"),
		},
	}

	for _, tt := range tests {
		a := tt.inputBytes
		var b []byte
		var c int = tt.inputBase
		ret := Strtol(a, b, c)
		if ret != tt.expectedValue {
			t.Errorf("Strtol() return %v. expected = %v", ret, tt.expectedValue)
		}
		if len(b) != len(tt.expectedString) {
			t.Errorf("Strtol() return %#v. expected = %#v", string(b), string(tt.expectedString))
		}
		for i := range b {
			if b[i] != tt.expectedString[i] {
				t.Errorf("Strtol() return %#v. expected = %#v", string(b), string(tt.expectedString))
			}
		}
	}
}
