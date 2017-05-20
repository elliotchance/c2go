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
		want []byte
	}{
		{"nil slice", args{nil}, nil},
		{"empty slice", args{[]byte{}}, []byte{}},
		{"single null-terminated", args{[]byte{0}}, []byte{}},
		{"multi null-terminated", args{[]byte{1, 0, 2, 0}}, []byte{1}},
		{"not null-terminated", args{[]byte{1, 2, 3}}, []byte{1, 2, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NullTerminatedByteSlice(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NullTerminatedByteSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}
