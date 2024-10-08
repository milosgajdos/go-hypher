package graph

import (
	"reflect"
	"testing"

	"github.com/milosgajdos/go-hypher"
)

func TestCopy(t *testing.T) {
	t.Run("EmptyGraph", func(t *testing.T) {
		g, err := NewGraph(hypher.WithLabel("foo"))
		if err != nil {
			t.Fatalf("failed to create graph: %v", err)
		}

		g2 := DeepCopy(g)

		if !reflect.DeepEqual(g, g2) {
			t.Fatalf("expected graphs to be equal g: %#v, g2: %#v", g, g2)
		}
	})

	t.Run("Non-EmptyGraph", func(t *testing.T) {
		g := MustGraph(t)
		g2 := DeepCopy(g)

		if !reflect.DeepEqual(g, g2) {
			t.Fatalf("expected graphs to be equal g: %#v, g2: %#v", g, g2)
		}
	})
}
