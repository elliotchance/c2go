package noarch

import (
	"reflect"
	"testing"
)

func TestNullTerminatedByteSlice(t *testing.T) {
	type args struct {
		s []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"nil slice", args{nil}, ""},
		{"empty slice", args{[]byte{}}, ""},
		{"single null-terminated", args{[]byte{0}}, ""},
		{"multi null-terminated", args{[]byte{'a', 0, 'b', 0}}, "a"},
		{"not null-terminated", args{[]byte{'a', 'b', 'c'}}, "abc"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NullTerminatedByteSlice(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NullTerminatedByteSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}
