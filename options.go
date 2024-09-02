package hypher

import "maps"

// RunMode is Graph run mode.
type RunMode int

const (
	// RunLevelMode runs node on the same graph level concurrently
	RunLevelMode RunMode = iota
	// RunAllMode runs all Graph nodes concurrently.
	RunAllMode
)

// Options configure graph.
type Options struct {
	// ID configures ID
	ID int64
	// UID configures UID
	UID string
	// Label configures Label.
	Label string
	// Attrs configures Attrs.
	Attrs map[string]any
	// DotID configures DOT ID
	DotID string
	// Weight configures Edge weight.
	Weight float64
	// Graph configures Node's graph
	Graph Graph
	// RunMode configures Graph run mode.
	RunMode RunMode
	// Op configures Node's Op.
	Op Op
}

// Option is functional graph option.
type Option func(*Options)

// WithID sets ID option.
func WithID(id int64) Option {
	return func(o *Options) {
		o.ID = id
	}
}

// WithUID sets UID option.
func WithUID(uid string) Option {
	return func(o *Options) {
		o.UID = uid
	}
}

// WithLabel sets Label option.
func WithLabel(label string) Option {
	return func(o *Options) {
		o.Label = label
	}
}

// WithAttrs sets Attrs option,
func WithAttrs(attrs map[string]any) Option {
	return func(o *Options) {
		o.Attrs = maps.Clone(attrs)
	}
}

// WithDotID sets DotID option.
func WithDotID(dotid string) Option {
	return func(o *Options) {
		o.DotID = dotid
	}
}

// WithWeight sets Weight option.
func WithWeight(weight float64) Option {
	return func(o *Options) {
		o.Weight = weight
	}
}

// WithGraph sets hypher Graph.
func WithGraph(g Graph) Option {
	return func(o *Options) {
		o.Graph = g
	}
}

// WithRunAll sets Parallel option.
func WithRunMode(mode RunMode) Option {
	return func(o *Options) {
		o.RunMode = mode
	}
}

// WithOp sets Op.
func WithOp(op Op) Option {
	return func(o *Options) {
		o.Op = op
	}
}
