package dot

import (
	"github.com/milosgajdos/go-hypher"
	"github.com/milosgajdos/go-hypher/graph"

	"gonum.org/v1/gonum/graph/encoding/dot"
)

// Marshaler is used for marshaling graph to DOT format.
type Marshaler struct {
	name       string
	prefix     string
	indent     string
	nodeStyle  Style
	edgeStyle  Style
	graphStyle Style
}

// NewMarshaler creates a new DOT graph marshaler and returns it.
func NewMarshaler(name, prefix, indent string, opts ...Option) (*Marshaler, error) {
	dotOpts := Options{
		NodeStyle:  DefaultNodeStyle(),
		EdgeStyle:  DefaultEdgeStyle(),
		GraphStyle: DefaultGraphStyle(),
	}

	for _, apply := range opts {
		apply(&dotOpts)
	}

	return &Marshaler{
		name:       name,
		prefix:     prefix,
		indent:     indent,
		nodeStyle:  dotOpts.NodeStyle,
		edgeStyle:  dotOpts.EdgeStyle,
		graphStyle: dotOpts.GraphStyle,
	}, nil
}

// Marshal marshal g into DOT and returns it.
func (m *Marshaler) Marshal(g hypher.Graph) ([]byte, error) {
	// Apply DOT styling

	hg := g.(*graph.Graph)
	hg.Attrs()["label"] = hg.Label()
	for k, v := range m.graphStyle.Attrs {
		hg.Attrs()[k] = v
	}

	nodes := g.Nodes()
	for nodes.Next() {
		n := nodes.Node().(*graph.Node)
		n.Attrs()["label"] = n.Label()
		n.Attrs()["shape"] = m.nodeStyle.Shape
		n.Attrs()["style"] = m.nodeStyle.Type
		for k, v := range m.nodeStyle.Attrs {
			n.Attrs()[k] = v
		}
	}

	edges := g.Edges()
	for edges.Next() {
		e := edges.Edge().(*graph.Edge)
		e.Attrs()["label"] = e.Label()
		e.Attrs()["shape"] = m.edgeStyle.Shape
		e.Attrs()["style"] = m.edgeStyle.Type
		for k, v := range m.edgeStyle.Attrs {
			e.Attrs()[k] = v
		}
	}

	return dot.Marshal(g, m.name, m.prefix, m.indent)
}
