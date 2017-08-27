package ast

import (
	"reflect"
	"testing"

	"github.com/elliotchance/c2go/util"
)

func TestArrayFiller(t *testing.T) {
	expected := &ArrayFiller{
		ChildNodes: []Node{},
	}
	actual := Parse(`array filler`)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("%s", util.ShowDiff(formatMultiLine(expected),
			formatMultiLine(actual)))
	}
}
