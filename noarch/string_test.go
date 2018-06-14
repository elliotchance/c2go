package noarch

import (
	"reflect"
	"testing"
)

func TestStringCopy(t *testing.T) {
	tests := []struct {
		name   string
		dst    *byte
		src    *byte
		length int32
		want   string
	}{
		{"src longer than length", &make([]byte, 4)[0], &[]byte("asdf")[0], 2, "as"},
		{"src shorter length", &make([]byte, 4)[0], &append([]byte("as"), 0)[0], 4, "as"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CStringToString(Strncpy(tt.dst, tt.src, tt.length)); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CStringToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
