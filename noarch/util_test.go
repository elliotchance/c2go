package noarch

import (
	"reflect"
	"testing"
)

func TestNullTerminatedBytePointer(t *testing.T) {
	type args struct {
		s *byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"nil slice", args{nil}, ""},
		{"single null-terminated", args{&[]byte{0}[0]}, ""},
		{"multi null-terminated", args{&[]byte{'a', 0, 'b', 0}[0]}, "a"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CStringToString(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CStringToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
