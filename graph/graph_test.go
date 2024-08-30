package graph

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/milosgajdos/go-hypher"
)

func MustGraph(t *testing.T, opts ...Option) *Graph {
	g, err := NewGraph(opts...)
	if err != nil {
		t.Fatal(err)
	}
	return g
}

func TestNewGraph(t *testing.T) {
	g, err := NewGraph()
	if err != nil {
		t.Fatalf("failed to create new graph: %v", err)
	}

	if uid := g.UID(); uid == "" {
		t.Error("expected non-empty UID")
	}

	if l := g.Label(); l != DefaultGraphLabel {
		t.Errorf("expected label: %s, got: %s", DefaultGraphLabel, l)
	}

	newLabel := "newLabel"
	g.SetLabel(newLabel)
	if l := g.Label(); l != newLabel {
		t.Errorf("expected label: %s, got: %s", newLabel, l)
	}

	newUID := "newUID"
	g.SetUID(newUID)
	if uid := g.UID(); uid != newUID {
		t.Errorf("expected UID: %s, got: %s", newUID, uid)
	}

	if a := g.Attrs(); a == nil {
		t.Error("expected non-empty attributes")
	}
}

func TestNewGraphWithOpts(t *testing.T) {
	uid := "FooID"
	label := "fooLabel"
	attrs := map[string]any{"foo": "bar"}

	g, err := NewGraph(
		WithUID(uid),
		WithLabel(label),
		WithAttrs(attrs),
	)
	if err != nil {
		t.Fatalf("failed to create new graph: %v", err)
	}

	if u := g.UID(); u != uid {
		t.Errorf("expected uid: %s, got: %s", uid, u)
	}

	if l := g.Label(); l != label {
		t.Errorf("expected label: %s, got: %s", label, l)
	}

	if !reflect.DeepEqual(g.Attrs(), attrs) {
		t.Errorf("expected attrs: %v, got: %v", attrs, g.Attrs())
	}
}

func TestGraphAddNode(t *testing.T) {
	g := MustGraph(t)
	n1 := MustNode(t)

	err := g.AddNode(n1)
	if err != nil {
		t.Fatalf("failed to add node to graph: %v", err)
	}

	if n1.ID() == NoneID {
		t.Errorf("AddNode must change the node ID; got: %v", n1.ID())
	}

	if n1.Graph() != g {
		t.Errorf("expected graph: %s, got: %s", g, n1.Graph())
	}

	err = g.AddNode(n1)
	if err != nil {
		t.Fatalf("failed to add existing node to graph: %v", err)
	}
}

func TestGraphSetEdge(t *testing.T) {
	g := MustGraph(t)
	n1 := MustNode(t)
	n2 := MustNode(t)
	e := MustEdge(t, n1, n2)

	err := g.SetEdge(e)
	if err != nil {
		t.Fatalf("failed to set edge in graph: %v", err)
	}

	if n1.Graph() != g || n2.Graph() != g {
		t.Error("nodes should be added to the graph when setting an edge")
	}

	// Adding the same edge again should not produce an error
	err = g.SetEdge(e)
	if err != nil {
		t.Fatalf("failed to set existing edge in graph: %v", err)
	}

	// Test cycle detection
	n3 := MustNode(t)
	e2 := MustEdge(t, n2, n3)
	e3 := MustEdge(t, n3, n1)

	err = g.SetEdge(e2)
	if err != nil {
		t.Fatalf("unexpected error setting edge %s: %v", e2, err)
	}
	err = g.SetEdge(e3)
	if err == nil {
		t.Error("expected cycle detection error when setting edge")
	}
}

func TestGraphInputsOutputs(t *testing.T) {
	g := MustGraph(t)
	n1 := MustNode(t, WithGraph(g))
	n2 := MustNode(t, WithGraph(g))

	g.SetInputs([]*Node{n1})
	g.SetOutputs([]*Node{n2})

	if !reflect.DeepEqual(g.Inputs(), []*Node{n1}) {
		t.Error("graph inputs not set correctly")
	}

	if !reflect.DeepEqual(g.Outputs(), []*Node{n2}) {
		t.Error("graph outputs not set correctly")
	}
}

func TestSubGraph(t *testing.T) {
	// Create a new graph
	g := MustGraph(t)

	// Create nodes
	nodes := make(map[int64]*Node)
	for i := int64(0); i < 7; i++ {
		node := MustNode(t, WithGraph(g))
		nodes[i] = node
	}

	// Create and add edges to the graph
	edges := []struct{ from, to int64 }{
		{0, 1}, {0, 3},
		{1, 5}, {1, 6},
		{2, 3}, {2, 4}, {2, 5},
		{3, 5},
	}
	for _, e := range edges {
		edge := MustEdge(t, nodes[e.from], nodes[e.to])
		if err := g.SetEdge(edge); err != nil {
			t.Fatalf("Failed to add edge %d->%d to graph: %v", e.from, e.to, err)
		}
	}

	t.Logf("graph nodes: %d, edges: %d", g.Nodes().Len(), g.Edges().Len())

	// Define input and output nodes
	inputNodes := Nodes{nodes[0], nodes[2]}
	outputNodes := Nodes{nodes[5], nodes[6]}
	//outputNodes := Nodes{nodes[5]}

	// Call SubGraph
	subgraph, err := g.SubGraph(inputNodes, outputNodes)
	if err != nil {
		t.Fatalf("SubGraph error: %v", err)
	}

	// Define expected nodes and edges in the subgraph
	expectedNodes := map[int64]bool{
		0: true,
		1: true,
		2: true,
		3: true,
		5: true,
		6: true,
	}
	expectedEdges := map[string]bool{
		"0->1": true, "0->3": true,
		"1->5": true, "1->6": true,
		"2->3": true, "2->5": true,
		"3->5": true,
	}

	// Check if all expected nodes are in the subgraph
	for id := range expectedNodes {
		if subgraph.Node(id) == nil {
			t.Errorf("Expected node %d to be in the subgraph, but it was not found", id)
		}
	}

	// Check if all nodes in the subgraph are expected
	for _, nodeID := range subgraph.nodes {
		if !expectedNodes[nodeID] {
			t.Errorf("Unexpected node %d found in the subgraph", nodeID)
		}
	}

	// Check if all expected edges are in the subgraph
	for _, e := range edges {
		edgeStr := fmt.Sprintf("%d->%d", e.from, e.to)
		if expectedEdges[edgeStr] {
			if subgraph.Edge(e.from, e.to) == nil {
				t.Errorf("Expected edge %s to be in the subgraph, but it was not found", edgeStr)
			}
			continue
		}
		if subgraph.Edge(e.from, e.to) != nil {
			t.Errorf("Unexpected edge %s found in the subgraph", edgeStr)
		}
	}

	// Check if all edges in the subgraph are expected
	sgEdges := subgraph.Edges()
	for sgEdges.Next() {
		edge := sgEdges.Edge().(*Edge)
		edgeStr := fmt.Sprintf("%d->%d", edge.From().ID(), edge.To().ID())
		if !expectedEdges[edgeStr] {
			t.Errorf("Unexpected edge %s found in the subgraph", edgeStr)
		}
	}
}

func TestGraphTopoSort(t *testing.T) {
	tests := []struct {
		name          string
		setupGraph    func() (*Graph, error)
		expectedOrder []int64
		expectError   bool
	}{
		{
			name: "Simple_graph",
			setupGraph: func() (*Graph, error) {
				g := MustGraph(t)
				n1 := MustNode(t, WithGraph(g))
				n2 := MustNode(t, WithGraph(g))
				n3 := MustNode(t, WithGraph(g))
				if err := g.SetEdge(MustEdge(t, n1, n2)); err != nil {
					return nil, err
				}
				if err := g.SetEdge(MustEdge(t, n2, n3)); err != nil {
					return nil, err
				}
				return g, nil
			},
			expectedOrder: []int64{0, 1, 2},
			expectError:   false,
		},
		{
			name: "Empty graph",
			setupGraph: func() (*Graph, error) {
				return MustGraph(t), nil
			},
			expectedOrder: []int64{},
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, err := tt.setupGraph()
			if err != nil {
				t.Fatalf("Failed to set up graph: %v", err)
			}

			sorted, err := g.TopoSort()

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected an error, but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if len(sorted) != len(tt.expectedOrder) {
					t.Errorf("Expected %d nodes, but got %d", len(tt.expectedOrder), len(sorted))
				}

				actualOrder := make([]int64, len(sorted))
				for i, node := range sorted {
					actualOrder[i] = node.ID()
				}

				if !reflect.DeepEqual(actualOrder, tt.expectedOrder) {
					t.Errorf("Expected order %v, but got %v", tt.expectedOrder, actualOrder)
				}
			}
		})
	}
}

const testOpKey = "test"

type testOp struct{}

func (t testOp) Type() string { return "testOp" }
func (t testOp) Desc() string { return "testOp sets inputs to outputs" }

func (t testOp) Do(_ context.Context, inputs ...hypher.Value) (hypher.Value, error) {
	return hypher.Value{testOpKey: inputs}, nil
}

func TestGraphRun(t *testing.T) {
	g := MustGraph(t)

	n0 := MustNode(t, WithGraph(g), WithOp(testOp{}))
	n1 := MustNode(t, WithGraph(g), WithOp(testOp{}))
	n2 := MustNode(t, WithGraph(g), WithOp(testOp{}))
	n3 := MustNode(t, WithGraph(g), WithOp(testOp{}))
	n4 := MustNode(t, WithGraph(g), WithOp(testOp{}))
	n5 := MustNode(t, WithGraph(g), WithOp(testOp{}))
	n6 := MustNode(t, WithGraph(g), WithOp(testOp{}))

	MustEdge(t, n0, n1, WithGraph(g))
	MustEdge(t, n0, n3, WithGraph(g))
	MustEdge(t, n1, n5, WithGraph(g))
	MustEdge(t, n1, n6, WithGraph(g))
	MustEdge(t, n2, n3, WithGraph(g))
	MustEdge(t, n2, n4, WithGraph(g))
	MustEdge(t, n2, n5, WithGraph(g))
	MustEdge(t, n3, n5, WithGraph(g))

	g.SetInputs([]*Node{n0, n2})
	g.SetOutputs([]*Node{n5})

	// node inputs
	_ = n1.SetInputs(hypher.Value{"ID": n1.ID()})
	_ = n3.SetInputs(hypher.Value{"ID": n3.ID()})
	_ = n4.SetInputs(hypher.Value{"ID": n4.ID()})
	_ = n5.SetInputs(hypher.Value{"ID": n5.ID()})
	_ = n6.SetInputs(hypher.Value{"ID": n6.ID()})

	graphInputs := map[string]hypher.Value{
		n0.UID(): {"ID": n0.ID()},
		n2.UID(): {"ID": n2.ID()},
	}

	ctx := context.Background()
	if err := g.Run(ctx, graphInputs); err != nil {
		t.Fatalf("run failed: %v", err)
	}

	expectedExecuted := map[int64]bool{
		n0.ID(): true,
		n1.ID(): true,
		n2.ID(): true,
		n3.ID(): true,
		n5.ID(): true,
	}

	// make sure only the right nodes have been executed
	for _, n := range []*Node{n0, n1, n2, n3, n4, n5, n6} {
		if !expectedExecuted[n.ID()] {
			if len(n.Outputs()) > 0 {
				t.Errorf("Node %d should not have been executed but has outputs", n.ID())
			}
			continue
		}
		if len(n.Outputs()) == 0 {
			t.Errorf("Node %d should have been executed but has no outputs", n.ID())
		}
	}

	// Check input propagation
	checkNodeOutput := func(n *Node, expectedInputCount int) {
		outputs := n.Outputs()
		if len(outputs) != 1 {
			t.Errorf("Node %d: expected 1 output, got %d", n.ID(), len(outputs))
			return
		}
		output := outputs[0]
		inputs, ok := output[testOpKey].([]hypher.Value)
		if !ok {
			t.Errorf("Node %d: output is not of type []Input", n.ID())
			return
		}
		if len(inputs) != expectedInputCount {
			t.Errorf("Node %d: expected %d inputs, got %d", n.ID(), expectedInputCount, len(inputs))
		}
	}

	checkNodeOutput(n0, 1) // Graph input
	checkNodeOutput(n2, 1) // Graph input
	checkNodeOutput(n1, 2) // Preset input + input from n0
	checkNodeOutput(n3, 3) // Preset input + inputs from n0 and n2
	checkNodeOutput(n5, 4) // Preset input + inputs from n1, n2, and n3
}
