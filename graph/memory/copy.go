package memory

import (
	"maps"

	"gonum.org/v1/gonum/graph/simple"
)

// NodeDeepCopy makes a deep copy of Node and returns it.
// It does not copy node inputs or outputs.
func NodeDeepCopy(n *Node) *Node {
	return &Node{
		id:    n.id,
		uid:   n.uid,
		dotid: n.dotid,
		label: n.label,
		attrs: maps.Clone(n.attrs),
		graph: n.graph,
		style: n.style,
	}
}

// EdgeDeepCopy makes a deep copy of Edge and returns it
func EdgeDeepCopy(e *Edge) *Edge {
	return &Edge{
		uid:    e.uid,
		label:  e.label,
		from:   NodeDeepCopy(e.From().(*Node)),
		to:     NodeDeepCopy(e.To().(*Node)),
		weight: e.weight,
		attrs:  maps.Clone(e.attrs),
		style:  e.style,
	}
}

// GraphDeepCopy return s deep copy of a memory graph.
func GraphDeepCopy(g *Graph) *Graph {
	g.mu.RLock()
	defer g.mu.RUnlock()

	cg := &Graph{
		WeightedDirectedGraph: simple.NewWeightedDirectedGraph(DefaultEdgeWeight, 0.0),
		uid:                   g.uid,
		dotid:                 g.dotid,
		label:                 g.label,
		attrs:                 maps.Clone(g.attrs),
		nodes:                 maps.Clone(g.nodes),
	}

	inputs := make([]*Node, 0, len(g.inputs))
	for _, n := range g.inputs {
		inputs = append(inputs, NodeDeepCopy(n))
	}
	cg.inputs = inputs

	outputs := make([]*Node, 0, len(g.outputs))
	for _, n := range g.outputs {
		outputs = append(outputs, NodeDeepCopy(n))
	}
	cg.outputs = outputs

	// copy all src nodes.
	nodes := g.Nodes()
	for nodes.Next() {
		n := nodes.Node().(*Node)
		node := NodeDeepCopy(n)
		if err := cg.AddNode(node); err != nil {
			panic("failed adding graph node")
		}
	}

	// copy all src edges.
	nodes.Reset()
	for nodes.Next() {
		nid := nodes.Node().ID()
		to := g.From(nid)
		for to.Next() {
			vid := to.Node().ID()
			e := g.WeightedEdge(nid, vid).(*Edge)
			edge := EdgeDeepCopy(e)
			cg.SetWeightedEdge(edge)
		}
	}

	return cg
}
