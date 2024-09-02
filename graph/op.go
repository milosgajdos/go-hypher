package graph

import (
	"context"
	"fmt"

	"github.com/milosgajdos/go-hypher"
)

// NoOp is a no-op Op.
type NoOp struct{}

func (op NoOp) Type() string   { return "NoOp" }
func (op NoOp) Desc() string   { return "NoOp does nothing" }
func (op NoOp) String() string { return fmt.Sprintf("Op: %s, Desc: %s", op.Type(), op.Desc()) }

func (op NoOp) Do(_ context.Context, _ ...hypher.Value) (hypher.Value, error) {
	return hypher.Value{}, nil
}
