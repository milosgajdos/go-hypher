package graph

import (
	"reflect"
	"testing"
)

// MustEdge creates a new Edge and returns it, panicking if there's an error.
func MustEdge(t *testing.T, from, to *Node, opts ...Option) *Edge {
	e, err := NewEdge(from, to, opts...)
	if err != nil {
		t.Fatalf("failed to create new edge: %v", err)
	}
	return e
}

func TestNewEdge(t *testing.T) {
	g := MustGraph(t)
	n1 := MustNode(t, WithGraph(g))
	n2 := MustNode(t, WithGraph(g))

	e := MustEdge(t, n1, n2)

	if e.UID() == "" {
		t.Error("expected non-empty UID")
	}

	newUID := "newUID"
	e.SetUID(newUID)
	if uid := e.UID(); uid != newUID {
		t.Errorf("expected UID: %s, got: %s", newUID, uid)
	}

	if e.Label() != DefaultEdgeLabel {
		t.Errorf("expected label: %s, got: %s", DefaultNodeLabel, e.Label())
	}

	if e.From() != n1 {
		t.Errorf("expected From node to be %v, got: %v", n1, e.From())
	}

	if e.To() != n2 {
		t.Errorf("expected To node to be %v, got: %v", n2, e.To())
	}

	if e.Weight() != DefaultEdgeWeight {
		t.Errorf("expected weight: %f, got: %f", DefaultEdgeWeight, e.Weight())
	}

	newWeight := 20.0
	e.SetWeight(newWeight)
	if w := e.Weight(); w != newWeight {
		t.Errorf("expected weight: %f, got: %f", newWeight, w)
	}

	newLabel := "newLabel"
	e.SetLabel(newLabel)
	if l := e.Label(); l != newLabel {
		t.Errorf("expected label: %s, got: %s", newLabel, l)
	}

	if s := e.Style(); s != DefaultEdgeStyleType {
		t.Errorf("expected style: %s, got: %s", DefaultEdgeStyleType, s)
	}

	if s := e.Shape(); s != DefaultEdgeShape {
		t.Errorf("expected shape: %s, got: %s", DefaultEdgeShape, s)
	}

	if c := e.Color(); c != DefaultEdgeColor {
		t.Errorf("expected color: %v, got: %v", DefaultEdgeColor, c)
	}
}

func TestNewEdgeWithOptions(t *testing.T) {
	g := MustGraph(t)

	var (
		testWeight = 2.0
		testLabel  = "test-label"
	)
	attrs := map[string]any{"foo": "bar"}

	from := MustNode(t, WithGraph(g), WithAttrs(attrs))
	to := MustNode(t, WithGraph(g), WithAttrs(attrs))

	e, err := NewEdge(from, to, WithLabel(testLabel), WithAttrs(attrs), WithWeight(testWeight))
	if err != nil {
		t.Fatalf("failed to create new edge: %v", err)
	}

	// TODO(milosgajdos): should we use big.NewFloat.Cmp?
	if w := e.Weight(); w != testWeight {
		t.Errorf("expected weight: %f, got: %f", testWeight, w)
	}

	if l := e.Label(); l != testLabel {
		t.Errorf("expected label: %s, got: %s", testLabel, l)
	}

	if a := e.Attrs(); !reflect.DeepEqual(a, attrs) {
		t.Errorf("expected attributes %v, got: %v", attrs, a)
	}

	// NOTE: we are adding some Style attributes to attrs
	// so the size of the returned map may be different from attrs.
	if a := e.Attributes(); len(a) < len(attrs) {
		t.Errorf("expected attributes count: %d, got: %d", len(attrs), len(a))
	}
}

func TestEdgeReversed(t *testing.T) {
	g := MustGraph(t)
	n1 := MustNode(t, WithGraph(g))
	n2 := MustNode(t, WithGraph(g))

	e1 := MustEdge(t, n1, n2)
	e2 := e1.ReversedEdge().(*Edge)

	if e1.From() != e2.To() || e1.To() != e2.From() {
		t.Error("reversed edge should swap From and To nodes")
	}

	if e1.UID() != e2.UID() {
		t.Error("reversed edge should maintain the same UID")
	}
}
