package cytoscape

import (
	"encoding/json"
	"fmt"

	"gonum.org/v1/gonum/graph/formats/cytoscapejs"

	"github.com/milosgajdos/go-hypher"
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
// CytoscapJS https://js.cytoscape.org/
func (m *Marshaler) Marshal(g hypher.Graph) ([]byte, error) {
	c := cytoscapejs.Elements{
		Nodes: make([]cytoscapejs.Node, 0, g.Nodes().Len()),
		Edges: make([]cytoscapejs.Edge, 0, g.Edges().Len()),
	}

	nodes := g.Nodes()
	for nodes.Next() {
		n := nodes.Node().(hypher.Node)

		ndata := cytoscapejs.NodeData{
			ID:         fmt.Sprint(n.ID()),
			Attributes: n.Attrs(),
		}

		c.Nodes = append(c.Nodes, cytoscapejs.Node{
			Data:       ndata,
			Selectable: true,
		})
	}

	edges := g.Edges()
	for edges.Next() {
		e := edges.Edge().(hypher.Edge)

		edata := cytoscapejs.EdgeData{
			ID:         e.UID(),
			Source:     fmt.Sprint(e.From().ID()),
			Target:     fmt.Sprint(e.To().ID()),
			Attributes: e.Attrs(),
		}

		c.Edges = append(c.Edges, cytoscapejs.Edge{
			Data:       edata,
			Selectable: true,
		})
	}

	return json.MarshalIndent(c, m.prefix, m.indent)
}
