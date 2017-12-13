package types

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/elliotchance/c2go/program"
	"github.com/elliotchance/c2go/util"

	goast "go/ast"
	"go/token"
)

func toJSON(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}

func TestCast(t *testing.T) {
	p := program.NewProgram()

	type args struct {
		expr     goast.Expr
		fromType string
		toType   string
	}
	tests := []struct {
		args args
		want goast.Expr
	}{
		// Casting to the same type is not needed.
		{args{util.NewIntLit(1), "int", "int"}, util.NewIntLit(1)},
		{args{util.NewFloatLit(2.3), "float", "float"}, util.NewFloatLit(2.3)},

		// Casting between numeric types.
		{args{util.NewIntLit(1), "int", "float"}, util.NewCallExpr("float32", util.NewIntLit(1))},
		{args{util.NewIntLit(1), "int", "double"}, util.NewCallExpr("float64", util.NewIntLit(1))},
		{args{util.NewIntLit(1), "int", "__uint16_t"}, util.NewCallExpr("uint16", util.NewIntLit(1))},

		// Casting to bool
		{args{util.NewIntLit(1), "int", "bool"}, util.NewBinaryExpr(util.NewIntLit(1), token.NEQ, util.NewIntLit(0), "bool", false)},
	}

	for _, tt := range tests {
		name := fmt.Sprintf("%#v", tt.args)

		t.Run(name, func(t *testing.T) {
			got, err := CastExpr(p, tt.args.expr, tt.args.fromType, tt.args.toType)

			if err != nil {
				t.Error(err)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cast()%s\n", util.ShowDiff(toJSON(got), toJSON(tt.want)))
			}
		})
	}
}

func TestGetArrayTypeAndSize(t *testing.T) {
	tests := []struct {
		in    string
		cType string
		size  int
	}{
		{"int", "int", -1},
		{"int [4]", "int", 4},
		{"int [4][3]", "int [3]", 4},
		{"int [4][3][2]", "int [3][2]", 4},
		{"int [4][3][2][1]", "int [3][2][1]", 4},
		{"int *[4]", "int *", 4},
		{"int *[4][3]", "int *[3]", 4},
		{"int *[4][3][2]", "int *[3][2]", 4},
		{"int *[4][3][2][1]", "int *[3][2][1]", 4},
		{"char *const", "char *const", -1},
		{"char *const [6]", "char *const", 6},
		{"char *const [6][5]", "char *const [5]", 6},
	}

	for _, tt := range tests {
		cType, size := GetArrayTypeAndSize(tt.in)
		if cType != tt.cType {
			t.Errorf("Expected type '%s', got '%s'", tt.cType, cType)
		}

		if size != tt.size {
			t.Errorf("Expected size '%d', got '%d'", tt.size, size)
		}
	}
}
