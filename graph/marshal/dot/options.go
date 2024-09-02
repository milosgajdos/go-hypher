package dot

// Options configure graph.
type Options struct {
	// NodeStyle configures Node style.
	NodeStyle Style
	// EdgeStyle configures Edge style.
	EdgeStyle Style
	// GraphStyle configures Graphe style.
	GraphStyle Style
}

// Option is functional graph option.
type Option func(*Options)

// WithNodeStyle sets Node Style.
func WithNodeStyle(s Style) Option {
	return func(o *Options) {
		o.NodeStyle = s
	}
}

// WithEdgeStyle sets Edge Style.
func WithEdgeStyle(s Style) Option {
	return func(o *Options) {
		o.EdgeStyle = s
	}
}

// WithGraphStyle sets Gra[h style.
func WithGraphStyle(s Style) Option {
	return func(o *Options) {
		o.GraphStyle = s
	}
}
