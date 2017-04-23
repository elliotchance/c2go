package types

import (
	"fmt"
	"testing"

	"github.com/elliotchance/c2go/program"
)

func TestResolveTypeForBinaryOperator(t *testing.T) {
	p := program.NewProgram()

	type args struct {
		operator  string
		leftType  string
		rightType string
	}
	tests := []struct {
		args args
		want string
	}{
		// Bitwise
		{args{"|", "int", "int"}, "int"},
		{args{"&", "int", "int"}, "int"},
		{args{"<<", "int", "int"}, "int"},
		{args{">>", "int", "int"}, "int"},

		// Comparison
		{args{"==", "int", "int"}, "bool"},
		{args{"==", "float", "int"}, "bool"},

		{args{"!=", "int", "int"}, "bool"},
		{args{"!=", "float", "int"}, "bool"},

		{args{">", "int", "int"}, "bool"},
		{args{">", "float", "int"}, "bool"},
		{args{">=", "int", "int"}, "bool"},
		{args{">=", "float", "int"}, "bool"},
		{args{"<", "int", "int"}, "bool"},
		{args{"<", "float", "int"}, "bool"},
		{args{"<=", "int", "int"}, "bool"},
		{args{"<=", "float", "int"}, "bool"},

		// Arithmetic
		{args{"+", "int", "int"}, "int"},
		{args{"+", "float", "float"}, "float"},

		{args{"-", "int", "int"}, "int"},
		{args{"-", "float", "float"}, "float"},

		{args{"*", "int", "int"}, "int"},
		{args{"*", "float", "float"}, "float"},

		{args{"/", "int", "int"}, "int"},
		{args{"/", "float", "float"}, "float"},
	}

	for _, tt := range tests {
		name := fmt.Sprintf("%#v", tt.args)

		t.Run(name, func(t *testing.T) {
			if got := ResolveTypeForBinaryOperator(p, tt.args.operator, tt.args.leftType, tt.args.rightType); got != tt.want {
				t.Errorf("ResolveTypeForBinaryOperator() = %v, want %v", got, tt.want)
			}
		})
	}
}
