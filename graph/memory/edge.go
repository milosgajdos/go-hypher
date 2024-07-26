package memory

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/google/uuid"

	gonum "gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding"

	"github.com/milosgajdos/go-hypher/graph/attrs"
	"github.com/milosgajdos/go-hypher/graph/style"
)

const (
	// DefaultEdgeLabel is the default edge label.
	DefaultEdgeLabel = "InMemoryEdge"
	// DefaultEdgeWeight is the default edge weight.
	DefaultEdgeWeight = 1.0
)

// Edge is a weighted graph edge.
type Edge struct {
	uid    string
	label  string
	from   *Node
	to     *Node
	weight float64
	attrs  map[string]any
	style  style.Style
}

// NewEdge creates a new edge and returns it.
func NewEdge(from, to *Node, opts ...Option) (*Edge, error) {
	eopts := Options{
		UID:    uuid.New().String(),
		Weight: DefaultEdgeWeight,
		Label:  DefaultEdgeLabel,
		Attrs:  make(map[string]any),
		Style:  style.DefaultEdge(),
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
func (e Edge) UID() string {
	return e.uid
}

// SetUID sets edge UID.
func (e *Edge) SetUID(uid string) {
	e.uid = uid
}

// Label returns edge label.
func (e Edge) Label() string {
	return e.label
}

// SetLabel sets edge label.
func (e *Edge) SetLabel(l string) {
	e.label = l
}

// From returns the from node of the first non-nil edge, or nil.
func (e *Edge) From() gonum.Node {
	return e.from
}

// To returns the to node of the first non-nil edge, or nil.
func (e *Edge) To() gonum.Node {
	return e.to
}

// Weight returns edge weight
func (e Edge) Weight() float64 {
	return e.weight
}

// SetWeight sets edge weight.
func (e *Edge) SetWeight(w float64) {
	e.weight = w
}

// ReversedEdge returns a new edge with end points of the pair swapped.
func (e *Edge) ReversedEdge() gonum.Edge {
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
	return e.attrs
}

// Style returns edge style.
func (e Edge) Style() string {
	return e.style.Type
}

// Shape returns edge shape.
func (e Edge) Shape() string {
	return e.style.Shape
}

// Color returns edge color.
func (e Edge) Color() color.RGBA {
	return e.style.Color
}

// Attributes returns node DOT attributes.
func (e Edge) Attributes() []encoding.Attribute {
	a := attrs.ToStringMap(e.attrs)
	attributes := make([]encoding.Attribute, len(a))

	i := 0
	for k, v := range a {
		attributes[i] = encoding.Attribute{Key: k, Value: v}
		i++
	}
	return attributes
}

// String implements fmt.Stringer.
func (e Edge) String() string {
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
