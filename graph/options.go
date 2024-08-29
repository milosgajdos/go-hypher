package graph

import (
	"maps"

	"github.com/milosgajdos/go-hypher"
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
	Attrs map[string]interface{}
	// DotID configures DOT ID
	DotID string
	// Weight configures weight.
	Weight float64
	// Graph configures node's graph
	Graph *Graph
	// Op configures node's Op.
	Op hypher.Op
	// Style configures style.
	Style Style
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
func WithAttrs(attrs map[string]interface{}) Option {
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

// WithGraph sets Graph options.
func WithGraph(g *Graph) Option {
	return func(o *Options) {
		o.Graph = g
	}
}

// WithOp sets Op.
func WithOp(op hypher.Op) Option {
	return func(o *Options) {
		o.Op = op
	}
}

// WithStyle sets Style option.
func WithStyle(s Style) Option {
	return func(o *Options) {
		o.Style = s
	}
}