// Package go-hypher enables creating AI agents as computational graphs.
//
// An agent is represented as a weighted Directed Acyclic Graph (DAG)
// which consists of nodes that perform a specific operation during
// agent execution aka Hypher Run.
//
// Hypher Run executes all the nodes in the hypher Graph and computes
// the results of their operations. All the outputs of the nodes
// are then passed to theis successors which then use them as their inputs.
// This continues all the way down to the hypher graph output nodes where
// the agent result is stored and can be fetched from.
//
// Given the agents are DAGs, they can form ensambles of agents
// through additional edges that link the agent DAGs as long
// as the resulting graph is also a DAG.
package hypher

import (
	"context"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding"
)

// Graph is weighted graph.
type Graph interface {
	graph.Weighted
	// UID returns graph UID.
	UID() string
	// Edges returns graph edges iterator.
	Edges() graph.Edges
	// Label returns graph label.
	Label() string
	// Attrs are graph attributes.
	Attrs() map[string]any
	// String is useful for debugging.
	String() string
}

// DOTNGraph is Graphviz DOT graph.
type DOTGraph interface {
	Graph
	encoding.Attributer
	// DOTID returns DOT ID.
	DOTID() string
	// SetDOTID sets DOT ID.
	SetDOTID(dotid string)
	// DOTAttributers sets DOT graph attributes.
	DOTAttributers() (graph, node, edge encoding.Attributer)
}

// Node is a graph node.
type Node interface {
	graph.Node
	// UID returns node UID.
	UID() string
	// Label returns node label.
	Label() string
	// Attrs returns node attributes.
	Attrs() map[string]any
	// String is useful for debugging.
	String() string
}

// Nodes is a slice of Nodes.
type Nodes []Node

// DOTNode is Graphviz DOT node.
type DOTNode interface {
	Node
	encoding.Attributer
	// DOTID returns DOT ID.
	DOTID() string
	// SetDOTID sets DOT ID.
	SetDOTID(dotid string)
}

// Edge is a graph edge.
type Edge interface {
	graph.WeightedEdge
	// UID returns edge UID.
	UID() string
	// Label returns edge label.
	Label() string
	// Attrs returns node attributes.
	Attrs() map[string]any
	// String is useful for debugging.
	String() string
}

// LabelSetter sets label.
type LabelSetter interface {
	SetLabel(string)
}

// UIDSetter sets UID.
type UIDSetter interface {
	SetUID(string)
}

// WeightSetter sets weight.
type WeightSetter interface {
	SetWeight(float64)
}

// DOTEdge is Graphviz DOT edge.
type DOTEdge interface {
	Edge
	encoding.Attributer
}

// Adder allows to add edges and nodes to graph.
type Adder interface {
	Graph
	graph.NodeAdder
	graph.WeightedEdgeAdder
}

// Remover allows to remove nodes and edges from graph.
type Remover interface {
	Graph
	graph.NodeRemover
	graph.EdgeRemover
}

// Updater allows to update graph.
type Updater interface {
	Adder
	Remover
}

// NodeUpdater adds and removes nodes.
type NodeUpdater interface {
	Graph
	graph.NodeAdder
	graph.NodeRemover
}

// EdgeUpdater adds and removes edges.
type EdgeUpdater interface {
	Graph
	graph.WeightedEdgeAdder
	graph.EdgeRemover
}

// Marshaler is used for marshaling graphs.
type Marshaler interface {
	// Marshal marshals graph into bytes.
	Marshal(g Graph) ([]byte, error)
}

// Unmarshaler is used for unmarshaling graphs.
type Unmarshaler interface {
	// Unmarshal unmarshals arbitrary bytes into graph.
	Unmarshal([]byte, Graph) error
}

// Syncer syncs the graph to a database or a filesystem.
type Syncer interface {
	Sync(context.Context, Graph) error
}

// Loader loads a graph from a databse or a filesystem.
type Loader interface {
	Load(context.Context, string) (Graph, error)
}

// Value is an I/O value.
type Value map[string]any

// Inputer returns its inputs.
type Inputer interface {
	// Inputs returns input values.
	Inputs() []Value
}

// Outputer returns its outputs.
type Outputer interface {
	// Outputs returns output values.
	Outputs() []Value
}

// Reseter resets inputs and outputs.
type Reseter interface {
	// Reset inputs and outputs.
	Reset()
}

// Runner is used to trigger a hypher Run.
// This usually means running the hypher Graph, by executing all its nodes.
type Runner interface {
	// Run an operation with the given inputs and options.
	Run(ctx context.Context, inputs map[string]Value, opts ...Option) error
}

// Execer executes a hypher operation.
// This usually means executing the hypher Node, by running its Op.
type Execer interface {
	// Exec executes an operation with the given inputs and returns the results.
	Exec(ctx context.Context, inputs ...Value) ([]Value, error)
}

// Op is an operation run by a Node.
type Op interface {
	// Type of the Op.
	Type() string
	// Desc describes the Op.
	Desc() string
	// Do runs the Op.
	Do(ctx context.Context, inputs ...Value) ([]Value, error)
	// String is useful for debugging.
	String() string
}
