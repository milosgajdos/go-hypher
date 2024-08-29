package graph

import (
	"image/color"
)

var (
	// DefaultEdgeColor is default edge color.
	DefaultEdgeColor = color.RGBA{R: 0, G: 0, B: 0}
	// DefaultNodeColor is default color for unknown entity.
	DefaultNodeColor = color.RGBA{R: 230, G: 230, B: 230}
)

const (
	// DefaultNodeStyleType is default style type.
	DefaultNodeStyleType = "rounded,filled,solid"
	// DefaultNodeShape is default node shape.
	DefaultNodeShape = "hexagon"
	// DefaultEdgeStyleType is default edge style.
	DefaultEdgeStyleType = ""
	// EdgeShape is default edge shape.
	DefaultEdgeShape = "normal"
	// UnknownShape is unknown shape.
	UnknownShape = "unknown"
)

// Style defines styling.
type Style struct {
	// Type is style type.
	Type string
	// Shape is style shape.
	Shape string
	// Color is style color.
	Color color.RGBA
}

// DefaultNodeStyle returns default node style
func DefaultNodeStyle() Style {
	return Style{
		Type:  DefaultNodeStyleType,
		Shape: DefaultNodeShape,
		Color: DefaultNodeColor,
	}
}

// DefaultEdgeStyle returns default edge style
func DefaultEdgeStyle() Style {
	return Style{
		Type:  DefaultEdgeStyleType,
		Shape: DefaultEdgeShape,
		Color: DefaultEdgeColor,
	}
}
