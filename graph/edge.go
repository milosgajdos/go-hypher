package graph

import (
	"fmt"
	"image/color"
	"strings"
	"sync"

	"github.com/google/uuid"
	gonum "gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding"

	"github.com/milosgajdos/go-hypher"
)

const (
	// DefaultEdgeLabel is the default edge label.
	DefaultEdgeLabel = "HypherEdge"
	// DefaultEdgeWeight is the default edge weight.
	DefaultEdgeWeight = 1.0
)

// Edge is a weighted graph edge.
type Edge struct {
	uid    string
	label  string
	from   hypher.Node
	to     hypher.Node
	weight float64
	attrs  map[string]any
	style  Style
	mu     sync.RWMutex
}

// NewEdge creates a new edge and returns it.
func NewEdge(from, to hypher.Node, opts ...Option) (*Edge, error) {
	eopts := Options{
		UID:    uuid.New().String(),
		Weight: DefaultEdgeWeight,
		Label:  DefaultEdgeLabel,
		Attrs:  make(map[string]any),
		Style:  DefaultEdgeStyle(),
	}

	for _, apply := range opts {
		apply(&eopts)
	}

	edge := &Edge{
		uid:    eopts.UID,
		from:   from,
		to:     to,
		weight: eopts.Weight,
		label:  eopts.Label,
		attrs:  eopts.Attrs,
		style:  eopts.Style,
	}

	if g := eopts.Graph; g != nil {
		if err := g.SetEdge(edge); err != nil {
			return nil, err
		}
	}

	return edge, nil
}

// UID returns edge UID.
func (e *Edge) UID() string {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.uid
}

// SetUID sets edge UID.
func (e *Edge) SetUID(uid string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.uid = uid
}

// Label returns edge label.
func (e *Edge) Label() string {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.label
}

// SetLabel sets edge label.
func (e *Edge) SetLabel(l string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.label = l
}

// From returns the from node of the first non-nil edge, or nil.
func (e *Edge) From() gonum.Node {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.from
}

// To returns the to node of the first non-nil edge, or nil.
func (e *Edge) To() gonum.Node {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.to
}

// Weight returns edge weight
func (e *Edge) Weight() float64 {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.weight
}

// SetWeight sets edge weight.
func (e *Edge) SetWeight(w float64) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.weight = w
}

// ReversedEdge returns a new edge with end points of the pair swapped.
func (e *Edge) ReversedEdge() gonum.Edge {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return &Edge{
		uid:    e.uid,
		from:   e.to,
		to:     e.from,
		label:  e.label,
		weight: e.weight,
		attrs:  e.attrs,
		style:  e.style,
	}
}

// Attrs returns node attributes.
func (e *Edge) Attrs() map[string]any {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.attrs
}

// Style returns edge style.
func (e *Edge) Style() string {
	return e.style.Type
}

// Shape returns edge shape.
func (e *Edge) Shape() string {
	return e.style.Shape
}

// Color returns edge color.
func (e *Edge) Color() color.RGBA {
	return e.style.Color
}

// Attributes returns node DOT attributes.
func (e *Edge) Attributes() []encoding.Attribute {
	e.mu.RLock()
	defer e.mu.RUnlock()

	styleAttrs := []encoding.Attribute{
		{Key: "label", Value: e.label},
		{Key: "shape", Value: e.style.Shape},
		{Key: "style", Value: e.style.Type},
	}

	a := AttrsToStringMap(e.attrs)
	attributes := make([]encoding.Attribute, 0, len(a))

	for k, v := range a {
		attributes = append(attributes, encoding.Attribute{Key: k, Value: v})
	}
	attributes = append(attributes, styleAttrs...)

	return attributes
}

// String implements fmt.Stringer.
func (e *Edge) String() string {
	e.mu.RLock()
	defer e.mu.RUnlock()

	var b strings.Builder
	fmt.Fprintf(&b, "Edge: %s\n", e.label)
	fmt.Fprintf(&b, "  UID: %s\n", e.uid)
	fmt.Fprintf(&b, "  From: Node(%d/%s)\n", e.from.ID(), e.from.UID())
	fmt.Fprintf(&b, "  To: Node(%d/%s)\n", e.to.ID(), e.to.UID())
	fmt.Fprintf(&b, "  Weight: %.2f\n", e.weight)

	if len(e.attrs) > 0 {
		fmt.Fprintf(&b, "  Attributes:\n")
		for k, v := range e.attrs {
			fmt.Fprintf(&b, "    %s: %v\n", k, v)
		}
	}

	return b.String()
}
