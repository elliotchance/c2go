package util

import (
	"testing"
)

var floats = []struct {
	in  float64
	out string
}{

	// Zeros have no decimal
	{0, "0"},
	{-0, "0"},

	{3.14159265, "3.14159265"},

	// + sign is included in positive exponents
	{3.14159265e123, "3.14159265e+123"},
	{3.14159265e+123, "3.14159265e+123"},

	{3.14159265e-123, "3.14159265e-123"},

	// "Small" exponents are not stored as exponents
	{3.14159265e+2, "314.159265"},
	{3.14159265e+5, "314159.265"},
}

func TestFloatLit(t *testing.T) {
	for _, tt := range floats {
		actual := NewFloatLit(tt.in).Value
		if tt.out != actual {
			t.Errorf("input: %v", tt.in)
			t.Errorf("  expected: %v", tt.out)
			t.Errorf("  actual:   %v", actual)
		}
	}
}
