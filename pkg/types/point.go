package types

import (
	"math"
)

// Point represents a 2D coordinate (X, Y)
type Point struct {
	X float64
	Y float64
}

// NewPoint creates a new Point with the specified coordinates
func NewPoint(x, y float64) Point {
	return Point{X: x, Y: y}
}

// DistanceTo calculates the Euclidean distance between two points
func (p Point) DistanceTo(other Point) float64 {
	dx := p.X - other.X
	dy := p.Y - other.Y
	return math.Sqrt(dx*dx + dy*dy)
}

// Add returns the result of adding another point to this point
func (p Point) Add(other Point) Point {
	return Point{
		X: p.X + other.X,
		Y: p.Y + other.Y,
	}
}

// Scale returns the result of scaling this point by a factor
func (p Point) Scale(factor float64) Point {
	return Point{
		X: p.X * factor,
		Y: p.Y * factor,
	}
}
