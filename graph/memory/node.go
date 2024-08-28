package memory

import (
	"context"
	"fmt"
	"image/color"
	"maps"
	"strings"
	"sync"

	"github.com/google/uuid"
	"gonum.org/v1/gonum/graph/encoding"

	"github.com/milosgajdos/go-hypher/graph/attrs"
	"github.com/milosgajdos/go-hypher/graph/style"
)

const (
	// DefaultNodeLabel is the default node label.
	DefaultNodeLabel = "InMemoryNode"
	// NoneID is non-existent ID.
	// Thanks Go for not having optionals!
	NoneID int64 = -1
)

// Nodes is a slice of Nodes.
type Nodes []*Node

// Node is a graph node.
type Node struct {
	id    int64
	uid   string
	dotid string
	label string
	attrs map[string]any
	graph *Graph
	style style.Style
	// node Op
	op Op
	// Node I/O
	inputs  []Value
	outputs []Value
	mu      sync.RWMutex
}

// NewNode creates a new Node and returns it.
func NewNode(opts ...Option) (*Node, error) {
	uid := uuid.New().String()
	nopts := Options{
		ID:    NoneID,
		UID:   uid,
		DotID: uid,
		Attrs: make(map[string]any),
		Style: style.DefaultNode(),
		Op:    NoOp{},
	}

	for _, apply := range opts {
		apply(&nopts)
	}

	node := &Node{
		id:      nopts.ID,
		uid:     nopts.UID,
		dotid:   nopts.DotID,
		label:   nopts.Label,
		attrs:   nopts.Attrs,
		graph:   nopts.Graph,
		style:   nopts.Style,
		op:      nopts.Op,
		inputs:  []Value{},
		outputs: []Value{},
	}

	if g := node.graph; g != nil {
		if err := g.AddNode(node); err != nil {
			return nil, err
		}
	}

	return node, nil
}

// ID returns node ID.
func (n *Node) ID() int64 {
	return n.id
}

// UID returns node UID.
func (n *Node) UID() string {
	n.mu.RLock()
	defer n.mu.RUnlock()

	return n.uid
}

// SetUID sets UID.
func (n *Node) SetUID(uid string) {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.uid = uid
}

// Label returns node label.
func (n *Node) Label() string {
	n.mu.RLock()
	defer n.mu.RUnlock()

	return n.label
}

// SetLabel sets node label.
func (n *Node) SetLabel(l string) {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.label = l
}

// Attrs returns node attributes.
func (n *Node) Attrs() map[string]any {
	n.mu.RLock()
	defer n.mu.RUnlock()

	return n.attrs
}

// Graph returns the node graph.
func (n *Node) Graph() *Graph {
	n.mu.RLock()
	defer n.mu.RUnlock()

	return n.graph
}

// Inputs return node inputs.
func (n *Node) Inputs() []Value {
	n.mu.RLock()
	defer n.mu.RUnlock()

	return n.inputs
}

// SetInputs sets the node inputs.
func (n *Node) SetInputs(inputs ...Value) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.inputs = inputs
	return nil
}

// Outputs returns node outputs.
func (n *Node) Outputs() []Value {
	n.mu.RLock()
	defer n.mu.RUnlock()

	return n.outputs
}

// Op returns node Op.
func (n *Node) Op() Op {
	n.mu.RLock()
	defer n.mu.RUnlock()

	return n.op
}

// Type returns the type of node style.
func (n *Node) Type() string {
	return n.style.Type
}

// Shape returns node shape.
func (n *Node) Shape() string {
	return n.style.Shape
}

// Color returns node color.
func (n *Node) Color() color.RGBA {
	return n.style.Color
}

// DOTID returns GraphVIz DOT ID.
func (n *Node) DOTID() string {
	n.mu.RLock()
	defer n.mu.RUnlock()

	return n.dotid
}

// SetDOTID sets GraphVIz DOT ID.
func (n *Node) SetDOTID(dotid string) {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.dotid = dotid
}

// Attributes returns node DOT attributes.
func (n *Node) Attributes() []encoding.Attribute {
	n.mu.RLock()
	defer n.mu.RUnlock()

	styleAttrs := []encoding.Attribute{
		{Key: "label", Value: n.label},
		{Key: "shape", Value: n.style.Shape},
		{Key: "style", Value: n.style.Type},
	}

	a := attrs.ToStringMap(n.attrs)
	attributes := make([]encoding.Attribute, 0, len(a))

	for k, v := range a {
		attributes = append(attributes, encoding.Attribute{Key: k, Value: v})
	}
	attributes = append(attributes, styleAttrs...)

	return attributes
}

// Node clones a node and returns it.
// The cloned node has a new UID.
// The node ID is reset to NoneID.
// The graph is not copied to the cloned node.
// No inputs or outputs are copied either.
func (n *Node) Clone() (*Node, error) {
	n.mu.RLock()
	defer n.mu.RUnlock()

	options := []Option{
		WithUID(uuid.New().String()),
		WithLabel(n.label),
		WithAttrs(maps.Clone(n.attrs)),
	}
	n2, err := NewNode(options...)
	if err != nil {
		return nil, err
	}

	return n2, nil
}

// CloneTo clones a node to graph g.
// The cloned node has a new UID
// even if g is the same as n.Graph().
func (n *Node) CloneTo(g *Graph) (*Node, error) {
	n.mu.RLock()
	defer n.mu.RUnlock()

	if g == nil {
		return nil, fmt.Errorf("invalid graph: %v", g)
	}

	n2, err := n.Clone()
	if err != nil {
		return nil, err
	}
	if err := g.AddNode(n2); err != nil {
		return nil, err
	}

	return n2, nil
}

// Exec executes a node Op and returns its result.
// It appends the output of the Op to its outputs.
func (n *Node) Exec(ctx context.Context, inputs ...Value) (Value, error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	opInputs := make([]Value, len(n.inputs)+len(inputs))
	copy(opInputs, n.inputs)
	copy(opInputs[len(n.inputs):], inputs)

	output, err := n.op.Do(ctx, opInputs...)
	if err != nil {
		return nil, fmt.Errorf("node %s op: %s error: %v", n.UID(), n.op.Desc(), err)
	}
	n.outputs = append(n.outputs, output)

	return output, nil
}

// String implements fmt.Stringer.
func (n *Node) String() string {
	n.mu.RLock()
	defer n.mu.RUnlock()

	var b strings.Builder
	fmt.Fprintf(&b, "Node: %s\n", n.label)
	fmt.Fprintf(&b, "  ID: %d\n", n.id)
	fmt.Fprintf(&b, "  UID: %s\n", n.uid)
	fmt.Fprintf(&b, "  DOTID: %s\n", n.dotid)
	if n.graph != nil {
		fmt.Fprintf(&b, "  Graph: %s\n", n.graph.UID())
	} else {
		fmt.Fprintf(&b, "  Graph: <not associated>\n")
	}

	if len(n.inputs) > 0 {
		fmt.Fprintf(&b, "  Inputs: %d\n", len(n.inputs))
	}
	if len(n.outputs) > 0 {
		fmt.Fprintf(&b, "  Outputs: %d\n", len(n.outputs))
	}

	if n.op != nil {
		fmt.Fprintf(&b, "  Op: %s, Desc: %s\n", n.op.Type(), n.op.Desc())
	}

	if len(n.attrs) > 0 {
		fmt.Fprintf(&b, "  Attributes:\n")
		for k, v := range n.attrs {
			fmt.Fprintf(&b, "    %s: %v\n", k, v)
		}
	}

	return b.String()
}
