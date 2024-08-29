package sigma

import (
	"encoding/json"
	"fmt"

	"github.com/milosgajdos/go-hypher"
	"gonum.org/v1/gonum/graph/formats/sigmajs"
)

// Marshaler implements graph.Marshaler.
type Marshaler struct {
	name   string
	prefix string
	indent string
}

// NewMarshaler creates a new Marshaler and returns it.
func NewMarshaler(name, prefix, indent string) (*Marshaler, error) {
	return &Marshaler{
		name:   name,
		prefix: prefix,
		indent: indent,
	}, nil
}

// Marshal marshals g into format that can be used by
// SigmaJS. See here for more: http://sigmajs.org/
func (m *Marshaler) Marshal(g hypher.Graph) ([]byte, error) {
	c := sigmajs.Graph{
		Nodes: make([]sigmajs.Node, 0, g.Nodes().Len()),
		Edges: make([]sigmajs.Edge, 0, g.Edges().Len()),
	}

	nodes := g.Nodes()
	for nodes.Next() {
		n := nodes.Node().(hypher.Node)

		c.Nodes = append(c.Nodes, sigmajs.Node{
			ID:         fmt.Sprint(n.ID()),
			Attributes: n.Attrs(),
		})
	}

	edges := g.Edges()
	for edges.Next() {
		e := edges.Edge().(hypher.Edge)

		c.Edges = append(c.Edges, sigmajs.Edge{
			ID:         e.UID(),
			Source:     fmt.Sprint(e.From().ID()),
			Target:     fmt.Sprint(e.To().ID()),
			Attributes: e.Attrs(),
		})
	}

	return json.MarshalIndent(c, m.prefix, m.indent)
}
