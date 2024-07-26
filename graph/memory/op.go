package memory

// Value is an I/O value.
type Value map[string]any

// Op is an operation run by a Node.
type Op interface {
	// Type of the Op.
	Type() string
	// Desc describes the Op.
	Desc() string
	// Do runs the Op.
	Do(...Value) (Value, error)
}

// NoOp is a no-op Op.
type NoOp struct{}

func (op NoOp) Type() string { return "NoOp" }
func (op NoOp) Desc() string { return "NoOp does nothing" }

func (op NoOp) Do(_ ...Value) (Value, error) {
	return Value{}, nil
}
