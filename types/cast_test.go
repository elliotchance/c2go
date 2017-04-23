package types

import (
	"testing"

	"fmt"

	"github.com/elliotchance/c2go/program"
)

func TestCast(t *testing.T) {
	p := program.NewProgram()

	type args struct {
		expr     string
		fromType string
		toType   string
	}
	tests := []struct {
		args args
		want string
	}{
		// Casting to the same type is not needed.
		{args{"1", "int", "int"}, "1"},
		{args{"2.3", "float", "float"}, "2.3"},

		// Casting between numeric types.
		{args{"1", "int", "float"}, "float32(1)"},
		{args{"1", "int", "double"}, "float64(1)"},
		{args{"1", "int", "__uint16_t"}, "uint16(1)"},

		// Casting to bool
		{args{"1", "int", "bool"}, "1 != 0"},

		// Casting from bool. This is a special case becuase C int and bool
		// values are very commonly used interchangably.
		{args{"1", "bool", "int"}, "noarch.BoolToInt(1)"},

		// String types
		// {args{"foo", "[3]char", "const char*"}, "1 != 0"},
	}

	for _, tt := range tests {
		name := fmt.Sprintf("%#v", tt.args)

		t.Run(name, func(t *testing.T) {
			if got := Cast(p, tt.args.expr, tt.args.fromType, tt.args.toType); got != tt.want {
				t.Errorf("Cast() = %v, want %v", got, tt.want)
			}
		})
	}
}
