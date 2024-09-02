package dot

import (
	"reflect"
	"testing"
)

func TestDefaultNodeStyle(t *testing.T) {
	dn := DefaultNodeStyle()

	s := Style{
		Type:  DefaultNodeStyleType,
		Shape: DefaultNodeShape,
		Color: DefaultNodeColor,
		Attrs: make(map[string]any),
	}

	if !reflect.DeepEqual(dn, s) {
		t.Fatal("unexpected default node style")
	}
}

func TestDefaultEdgeStyle(t *testing.T) {
	de := DefaultEdgeStyle()

	s := Style{
		Type:  DefaultEdgeStyleType,
		Shape: DefaultEdgeShape,
		Color: DefaultEdgeColor,
		Attrs: make(map[string]any),
	}

	if !reflect.DeepEqual(de, s) {
		t.Fatal("unexpected default edge style")
	}
}

func TestDefaultGraphStyle(t *testing.T) {
	dg := DefaultGraphStyle()

	s := Style{
		Attrs: map[string]any{
			"labelloc": "t",
		},
	}

	if !reflect.DeepEqual(dg, s) {
		t.Fatal("unexpected default graph style")
	}
}
