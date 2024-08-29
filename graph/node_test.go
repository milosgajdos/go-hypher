package graph

import (
	"reflect"
	"testing"

	"github.com/milosgajdos/go-hypher"
)

func MustNode(t *testing.T, opts ...Option) *Node {
	n, err := NewNode(opts...)
	if err != nil {
		t.Fatalf("failed to create new node: %v", err)
	}
	return n
}

func TestNewNode(t *testing.T) {
	n := MustNode(t)

	if n.ID() != NoneID {
		t.Errorf("expected ID to be NoneID, got: %d", n.ID())
	}

	if n.UID() == "" {
		t.Error("expected non-empty UID")
	}

	if n.Label() != DefaultNodeLabel {
		t.Errorf("expected label: %s, got: %s", DefaultNodeLabel, n.Label())
	}

	newLabel := "newLabel"
	n.SetLabel(newLabel)
	if l := n.Label(); l != newLabel {
		t.Errorf("expected label: %s, got: %s", newLabel, l)
	}

	if n.Graph() != nil {
		t.Error("expected nil graph")
	}

	if s := n.Type(); s != DefaultNodeStyleType {
		t.Errorf("expected type: %s, got: %s", DefaultNodeStyleType, s)
	}

	if s := n.Shape(); s != DefaultNodeShape {
		t.Errorf("expected shape: %s, got: %s", DefaultNodeShape, s)
	}

	if c := n.Color(); c != DefaultNodeColor {
		t.Errorf("expected color: %v, got: %v", DefaultNodeColor, c)
	}

	if d := n.DOTID(); d != n.UID() {
		t.Errorf("expected dotid: %s, got: %s", n.UID(), d)
	}

	testDotID := "testDotID"
	n.SetDOTID(testDotID)

	if d := n.DOTID(); d != testDotID {
		t.Errorf("expected dotid: %s, got: %s", testDotID, d)
	}

	newUID := "newUID"
	n.SetUID(newUID)
	if uid := n.UID(); uid != newUID {
		t.Errorf("expected UID: %s, got: %s", newUID, uid)
	}

	if n.Op() == nil {
		t.Error("expected non-nil Op")
	}

	if len(n.Inputs()) != 0 {
		t.Errorf("expected empty inputs, got: %v", n.Inputs())
	}

	if len(n.Outputs()) != 0 {
		t.Errorf("expected empty outputs, got: %v", n.Outputs())
	}
}

func TestNewNodeWithOptions(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	var (
		testDotID = "testDotID"
		testLabel = "testLabel"
		testUID   = "testUID"
		testID    = int64(100)
	)

	attrs := map[string]interface{}{
		"foo": "bar",
	}

	style := Style{
		Type:  "foo",
		Shape: "bar",
		Color: DefaultNodeColor,
	}

	opts := []Option{
		WithID(testID),
		WithLabel(testLabel),
		WithAttrs(attrs),
		WithUID(testUID),
		WithDotID(testDotID),
		WithStyle(style),
	}

	n, err := NewNode(opts...)
	if err != nil {
		t.Fatalf("failed to create new node: %v", err)
	}

	if s := n.Type(); s != style.Type {
		t.Errorf("expected style: %s, got: %s", style.Type, s)
	}

	if s := n.Shape(); s != style.Shape {
		t.Errorf("expected shape: %s, got: %s", style.Shape, s)
	}

	if c := n.Color(); c != style.Color {
		t.Errorf("expected color: %v, got: %v", style.Color, c)
	}

	if id := n.ID(); id != testID {
		t.Errorf("expected id: %d, got: %d", testID, id)
	}

	if d := n.DOTID(); d != testDotID {
		t.Errorf("expected dotid: %s, got: %s", testDotID, d)
	}

	if u := n.UID(); u != testUID {
		t.Errorf("expected uid: %s, got: %s", testUID, u)
	}
}

func TestNodeWithGraph(t *testing.T) {
	g := MustGraph(t)
	n := MustNode(t, WithGraph(g))

	if n.ID() == NoneID {
		t.Error("expected non-NoneID")
	}

	if n.Graph() != g {
		t.Error("expected node to be added to the graph")
	}
}

func TestNodeClone(t *testing.T) {
	n1 := MustNode(t, WithLabel("TestNode"), WithAttrs(map[string]any{"key": "value"}))
	n2, err := n1.Clone()
	if err != nil {
		t.Fatalf("failed to clone node: %v", err)
	}

	if n1.UID() == n2.UID() {
		t.Error("cloned node should have a different UID")
	}

	if n1.Label() != n2.Label() {
		t.Errorf("expected same label, got: %s and %s", n1.Label(), n2.Label())
	}

	if !reflect.DeepEqual(n1.Attrs(), n2.Attrs()) {
		t.Errorf("expected same attributes, got: %v and %v", n1.Attrs(), n2.Attrs())
	}

	if n2.Graph() != nil {
		t.Error("cloned node should not be associated with a graph")
	}
}

func TestNodeCloneTo(t *testing.T) {
	g1 := MustGraph(t)
	n1 := MustNode(t, WithGraph(g1))

	g2 := MustGraph(t)
	n2, err := n1.CloneTo(g2)
	if err != nil {
		t.Fatalf("failed to clone node to new graph: %v", err)
	}

	if n2.Graph() != g2 {
		t.Error("cloned node should be associated with the new graph")
	}

	if n1.UID() == n2.UID() {
		t.Error("cloned node should have a different UID in the new graph")
	}
}

func TestNodeSetInputs(t *testing.T) {
	n := MustNode(t)
	inputs := []hypher.Value{{"foo": 10}}

	err := n.SetInputs(inputs...)
	if err != nil {
		t.Fatalf("failed to set inputs: %v", err)
	}

	if !reflect.DeepEqual(n.Inputs(), inputs) {
		t.Errorf("expected inputs: %v, got: %v", inputs, n.Inputs())
	}
}
