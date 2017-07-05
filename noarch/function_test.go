package noarch

import (
	"fmt"
	"testing"
)

func TestBoolToInt(t *testing.T) {
	tests := []struct {
		input bool
		want  int
	}{
		{true, 1},
		{false, 0},
	}
	for _, tt := range tests {
		name := fmt.Sprintf("%#v", tt.input)

		t.Run(name, func(t *testing.T) {
			if got := BoolToInt(tt.input); got != tt.want {
				t.Errorf("BoolToInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func testB() interface{} {
	return "b"
}
func testC() interface{} {
	return "c"
}

func TestTernary(t *testing.T) {
	if result := Ternary(true, testB, testC); result.(string) != "b" {
		t.Errorf("Ternary - true is Fail")
	}
	if result := Ternary(false, testB, testC); result.(string) != "c" {
		t.Errorf("Ternary - false is Fail")
	}
}
