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

func TestNotInt(t *testing.T) {
	tests := []struct {
		input int
		want  int
	}{
		{0, 1},
		{1, 0},
		{42, 0},
	}
	for _, tt := range tests {
		name := fmt.Sprintf("%#v", tt.input)

		t.Run(name, func(t *testing.T) {
			if got := NotInt(tt.input); got != tt.want {
				t.Errorf("NotInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func funcB() interface{} {
	return "b"
}
func funcC() interface{} {
	return "c"
}

func TestTernary(t *testing.T) {
	if result := Ternary(true, funcB, funcC); result.(string) != "b" {
		t.Errorf("Ternary - true is Fail")
	}
	if result := Ternary(false, funcB, funcC); result.(string) != "c" {
		t.Errorf("Ternary - false is Fail")
	}
}
