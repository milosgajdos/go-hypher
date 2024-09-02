package graph

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
	gonum "gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/topo"

	"github.com/milosgajdos/go-hypher"
)

const (
	// DefaultGraphLabel is the default label.
	DefaultGraphLabel = "HypherGraph"
)

// Graph is an in-memory graph.
type Graph struct {
	*simple.WeightedDirectedGraph
	// graph metadata
	uid   string
	dotid string
	label string
	attrs map[string]any
	// node cache
	nodes map[string]int64
	// input and output nodes
	inputs  []*Node
	outputs []*Node
	mu      sync.RWMutex
}

// NewGraph creates a new graph and returns it.
func NewGraph(opts ...hypher.Option) (*Graph, error) {
	uid := uuid.New().String()
	gopts := hypher.Options{
		UID:    uid,
		DotID:  uid,
		Label:  DefaultGraphLabel,
		Weight: DefaultEdgeWeight,
		Attrs:  make(map[string]any),
	}

	for _, apply := range opts {
		apply(&gopts)
	}

	return &Graph{
		WeightedDirectedGraph: simple.NewWeightedDirectedGraph(gopts.Weight, 0.0),
		uid:                   gopts.UID,
		dotid:                 gopts.DotID,
		label:                 gopts.Label,
		attrs:                 gopts.Attrs,
		nodes:                 make(map[string]int64),
		inputs:                []*Node{},
		outputs:               []*Node{},
	}, nil
}

// UID returns graph UID.
func (g *Graph) UID() string {
	return g.uid
}

// Label returns graph label.
func (g *Graph) Label() string {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return g.label
}

// SetLabel sets label.
func (g *Graph) SetLabel(l string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.label = l
}

// SetUID sets UID.
func (g *Graph) SetUID(uid string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.uid = uid
}

// Attrs returns graph attributes.
// TODO: consider cloning these
func (g *Graph) Attrs() map[string]any {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return g.attrs
}

// DOTID returns GraphVIz DOT ID.
func (g *Graph) DOTID() string {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return g.dotid
}

// SetDOTID sets GraphVIz DOT ID.
func (g *Graph) SetDOTID(dotid string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.dotid = dotid
}

// DOTAttributers are graph.Graph values that specify top-level DOT attributes
// TODO: figure out node and edge top level attributes
func (g *Graph) DOTAttributers() (graph, node, edge encoding.Attributer) {
	return g, nil, nil
}

// Attributes returns graph DOT attributes.
func (g *Graph) Attributes() []encoding.Attribute {
	g.mu.RLock()
	defer g.mu.RUnlock()

	a := AttrsToStringMap(g.attrs)
	attributes := make([]encoding.Attribute, 0, len(a))

	for k, v := range a {
		attributes = append(attributes, encoding.Attribute{Key: k, Value: v})
	}

	return attributes
}

// HasEdgeFromTo returns whether an edge exist between two nodoes with the given IDs.
func (g *Graph) HasEdgeFromTo(uid, vid int64) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return g.WeightedDirectedGraph.HasEdgeBetween(uid, vid)
}

// To returns all nodes that can reach directly to the node with the given ID.
func (g *Graph) To(id int64) gonum.Nodes {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return g.WeightedDirectedGraph.To(id)
}

// SetInputs sets graph input nodes.
func (g *Graph) SetInputs(nodes []*Node) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.inputs = nodes
}

// Inputs returns graph input nodes.
func (g *Graph) Inputs() []*Node {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return g.inputs
}

// SetOutputs sets graph output nodes.
func (g *Graph) SetOutputs(nodes []*Node) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.outputs = nodes
}

// Outputs returns graph output nodes.
// TODO: consider cloning outputs
func (g *Graph) Outputs() []*Node {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return g.outputs
}

// NewNode creates a new node and adds it to the graph.
// It returns the new node or fails with error.
func (g *Graph) NewNode(opts ...hypher.Option) (*Node, error) {
	opts = append(opts, hypher.WithGraph(g))
	return NewNode(opts...)
}

// nodeExists returns true if the node already exists in g.
func (g *Graph) nodeExists(n *Node) bool {
	if node := g.Node(n.ID()); node != nil {
		if _, ok := g.nodes[n.UID()]; !ok {
			g.nodes[n.UID()] = n.ID()
		}
		return true
	}

	return false
}

// AddNode adds a node to the graph or returns error.
// If the node's graph is the same as g it returns nil.
// Otherwise, it tries to preserve the ID of the node
// unless the ID is not set (NoneID) or a node
// with the same ID already exists in g, in which
// case a new ID is generated before the node is added
// to g. If node's graph is nil, it's set to g.
func (g *Graph) AddNode(n *Node) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if n.Graph() != nil {
		if g.uid == n.Graph().UID() {
			if g.nodeExists(n) {
				return nil
			}
		}
	}

	// if it has no ID or a node with the same ID already exists in g.
	if n.ID() == NoneID || g.Node(n.ID()) != nil {
		node := g.WeightedDirectedGraph.NewNode()
		n.id = node.ID()
	}

	n.graph = g
	g.WeightedDirectedGraph.AddNode(n)
	g.nodes[n.UID()] = n.ID()

	return nil
}

// NewEdge creates a new edge link its node in the graph.
// It returns the new edge or fails with error.
func (g *Graph) NewEdge(from, to hypher.Node, opts ...hypher.Option) (*Edge, error) {
	opts = append(opts, hypher.WithGraph(g))
	return NewEdge(from, to, opts...)
}

// SetEdge adds the edge e to the graph linking the edge nodes.
// It adds the edge nodes to the graph if they don't already exist.
// It returns error if the new edge creates a graph cycle.
func (g *Graph) SetEdge(e hypher.Edge) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	fromNode, ok := e.From().(*Node)
	if !ok {
		return fmt.Errorf("invalid From node: %T", e.From())
	}

	toNode, ok := e.To().(*Node)
	if !ok {
		return fmt.Errorf("invalid To node: %T", e.To())
	}

	if edge := g.Edge(e.From().ID(), e.To().ID()); edge != nil {
		return nil
	}

	fromNodeID, toNodeID := fromNode.id, toNode.id
	fromNodeGraph, toNodeGraph := fromNode.graph, toNode.graph

	var fromAdded, toAdded bool

	// if the nodes do not exist in the graph we must create them
	// we can't create an edge between nodes that are not in the graph.
	if fromNode.ID() == NoneID || g.Node(fromNode.ID()) == nil {
		fromNode.id = g.WeightedDirectedGraph.NewNode().ID()
		fromNode.graph = g
		g.WeightedDirectedGraph.AddNode(fromNode)
		g.nodes[fromNode.UID()] = fromNode.ID()
		fromAdded = true
	}
	if toNode.ID() == NoneID || g.Node(toNode.ID()) == nil {
		toNode.id = g.WeightedDirectedGraph.NewNode().ID()
		toNode.graph = g
		g.WeightedDirectedGraph.AddNode(toNode)
		g.nodes[toNode.UID()] = toNode.ID()
		toAdded = true
	}
	g.SetWeightedEdge(e)

	// check if there is a cycle
	if topo.PathExistsIn(g, g.Node(e.To().ID()), g.Node(e.From().ID())) {
		// remove the edge and the nodes that have just been created
		g.RemoveEdge(fromNode.ID(), toNode.ID())
		// remove nodes if they had been created
		// and reset their IDs and graphs
		if fromAdded {
			g.RemoveNode(fromNode.ID())
			fromNode.id = fromNodeID
			fromNode.graph = fromNodeGraph
			delete(g.nodes, fromNode.uid)
		}
		if toAdded {
			g.RemoveNode(toNode.ID())
			toNode.id = toNodeID
			toNode.graph = toNodeGraph
			delete(g.nodes, toNode.uid)
		}
		return fmt.Errorf("cycle detected when adding edge: %s", e)
	}

	return nil
}

func (g *Graph) buildSubGraph(sg *Graph, n *Node, outputNodes map[int64]struct{}) (bool, error) {
	if sg.Node(n.ID()) != nil {
		return true, nil
	}

	if _, isOutput := outputNodes[n.ID()]; isOutput {
		if err := sg.AddNode(n); err != nil {
			return false, err
		}
		return true, nil
	}

	nodeInPathToOut := false

	for _, succ := range gonum.NodesOf(g.From(n.ID())) {
		succNode := succ.(*Node)
		succInPathToOut, err := g.buildSubGraph(sg, succNode, outputNodes)
		if err != nil {
			return false, err
		}
		if succInPathToOut {
			nodeInPathToOut = true
			if sg.Node(n.ID()) == nil {
				if err := sg.AddNode(n); err != nil {
					return false, err
				}
			}
			if sg.Node(succNode.ID()) == nil {
				if err := sg.AddNode(succNode); err != nil {
					return false, err
				}
			}
			edge := g.Edge(n.ID(), succNode.ID()).(*Edge)
			if err := sg.SetEdge(edge); err != nil {
				return false, err
			}
		}
	}

	return nodeInPathToOut, nil
}

// SubGraph returns a sub-graph of g which contains all the nodes
// which are either outputNodes or are on the path to the outputNodes
// when starting the graph traversal in inputNodes, including the inputNodes.
func (g *Graph) SubGraph(inputNodes, outputNodes Nodes) (*Graph, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	sg, err := NewGraph()
	if err != nil {
		return nil, err
	}

	// Set of output node IDs for a quick lookup
	outputSet := make(map[int64]struct{})
	for _, out := range outputNodes {
		outputSet[out.ID()] = struct{}{}
	}

	for _, inputNode := range inputNodes {
		_, err := g.buildSubGraph(sg, inputNode, outputSet)
		if err != nil {
			return nil, err
		}
	}

	return sg, nil
}

// TopoSort performs a topological sort of the graph.
// It returns nodes sorted in ascending order of their
// incoming edge counts (i.e. in degree) starting
// with nodes with zero incoming edges aka "roots".
// If there isn't at least one node with zero
// incoming edge the graph must ehter have a cycle,
// or have no nodes, so an empty slice is returned.
func (g *Graph) TopoSort() ([]gonum.Node, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	sorted, err := topo.Sort(g)
	if err != nil {
		return nil, err
	}

	return sorted, nil
}

// TopoSortWithLevels does the same topological sort
// as TopoSort, but groups the sorted nodes to graph levels.
func (g *Graph) TopoSortWithLevels() ([][]gonum.Node, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	sorted, err := topo.Sort(g)
	if err != nil {
		return nil, err
	}

	levels := make(map[int64]int)
	maxLevel := 0

	for _, node := range sorted {
		level := 0
		predecessors := g.To(node.ID())
		for predecessors.Next() {
			pred := predecessors.Node()
			if levels[pred.ID()] >= level {
				level = levels[pred.ID()] + 1
			}
		}
		levels[node.ID()] = level
		if level > maxLevel {
			maxLevel = level
		}
	}

	components := make([][]gonum.Node, maxLevel+1)
	for _, node := range sorted {
		level := levels[node.ID()]
		components[level] = append(components[level], node)
	}

	return components, nil
}

func (g *Graph) execNodeWait(ctx context.Context, node *Node, nodeChans map[int64]chan struct{}) error {
	// TODO: handle the case where we might not want to
	// aggregate all the predecessor outputs into Exec input;
	// We might want to do an exec for each predecessor output
	var nodeInputs []hypher.Value
	to := g.To(node.ID())
	for to.Next() {
		pred := to.Node()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-nodeChans[pred.ID()]: // Wait for the predecessor to finish
			predOutputs := pred.(*Node).Outputs()
			nodeInputs = append(nodeInputs, predOutputs...)
		}
	}

	// exec the node
	if _, err := node.Exec(ctx, nodeInputs...); err != nil {
		return err
	}

	// Signal completion to dependent nodes aka successors
	close(nodeChans[node.ID()])

	return nil
}

func (g *Graph) runAll(ctx context.Context) error {
	// get the execution (sub)graph
	sg, err := g.SubGraph(g.inputs, g.outputs)
	if err != nil {
		return err
	}

	// sort the returned graph topologically
	nodes, err := sg.TopoSort()
	if err != nil {
		return err
	}

	eg, egCtx := errgroup.WithContext(ctx)

	// Create a map to store the channels for each node
	nodeChans := make(map[int64]chan struct{})
	for _, node := range nodes {
		nodeChans[node.ID()] = make(chan struct{})
	}

	// Start a goroutine for each node
	for _, node := range nodes {
		// NOTE: we could also just pass the node UID
		node := node
		eg.Go(func() error {
			return g.execNodeWait(egCtx, node.(*Node), nodeChans)
		})
	}

	// Wait for all goroutines to complete or for an error to occur
	if err := eg.Wait(); err != nil {
		return fmt.Errorf("graph run failed: %v", err)
	}

	return nil
}

func (g *Graph) execNode(ctx context.Context, node *Node) error {
	var nodeInputs []hypher.Value
	to := g.To(node.ID())

	for to.Next() {
		pred := to.Node()
		predOutputs := pred.(*Node).Outputs()
		nodeInputs = append(nodeInputs, predOutputs...)
	}

	// exec the node
	if _, err := node.Exec(ctx, nodeInputs...); err != nil {
		return err
	}

	return nil
}

func (g *Graph) run(ctx context.Context) error {
	// get the execution (sub)graph
	sg, err := g.SubGraph(g.inputs, g.outputs)
	if err != nil {
		return err
	}

	nodeLevels, err := sg.TopoSortWithLevels()
	if err != nil {
		return err
	}

	for _, nodes := range nodeLevels {
		// run all nodes on the same level in parallel
		eg, egCtx := errgroup.WithContext(ctx)
		for _, node := range nodes {
			node := node
			eg.Go(func() error {
				return g.execNode(egCtx, node.(*Node))
			})
		}
		if err := eg.Wait(); err != nil {
			return fmt.Errorf("graph run failed: %v", err)
		}
	}

	return nil
}

// Run runs the graph with the given inputs.
// The inputs are passed in to the input nodes.
// Run executes all the graph nodes operations.
// Run is a blocking call. It returns when
// when the graph execution finished or if any
// of the executed nodes Op failed with error.
func (g *Graph) Run(ctx context.Context, inputs map[string]hypher.Value, opts ...hypher.Option) error {
	// NOTE: we only read Parallel option.
	gopts := hypher.Options{}
	for _, apply := range opts {
		apply(&gopts)
	}

	// set inputs to all the input nodes.
	for _, node := range g.inputs {
		if nodeInput, ok := inputs[node.UID()]; ok {
			if err := node.SetInputs(nodeInput); err != nil {
				return err
			}
		}
	}

	if gopts.RunMode == hypher.RunAllMode {
		return g.runAll(ctx)
	}

	return g.run(ctx)
}

// String implements fmt.Stringer.
func (g *Graph) String() string {
	g.mu.RLock()
	defer g.mu.RUnlock()

	var b strings.Builder
	fmt.Fprintf(&b, "Graph: %s\n", g.label)
	fmt.Fprintf(&b, "  UID: %s\n", g.uid)
	fmt.Fprintf(&b, "  Nodes: %d\n", g.Nodes().Len())
	fmt.Fprintf(&b, "  Edges: %d\n", g.Edges().Len())

	if len(g.inputs) > 0 {
		fmt.Fprintf(&b, "  Input Nodes: %d\n", len(g.inputs))
	}
	if len(g.outputs) > 0 {
		fmt.Fprintf(&b, "  Output Nodes: %d\n", len(g.outputs))
	}

	if len(g.attrs) > 0 {
		fmt.Fprintf(&b, "  Attributes:\n")
		for k, v := range g.attrs {
			fmt.Fprintf(&b, "    %s: %v\n", k, v)
		}
	}

	return b.String()
}
