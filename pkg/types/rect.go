package types

// Rect represents a rectangular boundary with position and dimensions
type Rect struct {
	X      float64 // X coordinate of the top-left corner
	Y      float64 // Y coordinate of the top-left corner
	Width  float64 // Width of the rectangle
	Height float64 // Height of the rectangle
}

// NewRect creates a new Rect with the specified position and dimensions
func NewRect(x, y, width, height float64) Rect {
	return Rect{
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
	}
}

// Contains checks if a point is within the boundaries of this rectangle
func (r Rect) Contains(p Point) bool {
	return p.X >= r.X && p.X < r.X+r.Width &&
		p.Y >= r.Y && p.Y < r.Y+r.Height
}

// Center returns the center point of the rectangle
func (r Rect) Center() Point {
	return Point{
		X: r.X + r.Width/2,
		Y: r.Y + r.Height/2,
	}
}

// GetX returns the minimum X coordinate (left edge)
func (r Rect) GetX() float64 {
	return r.X
}

// GetY returns the minimum Y coordinate (top edge)
func (r Rect) GetY() float64 {
	return r.Y
}

// GetMaxX returns the maximum X coordinate (right edge)
func (r Rect) GetMaxX() float64 {
	return r.X + r.Width
}

// GetMaxY returns the maximum Y coordinate (bottom edge)
func (r Rect) GetMaxY() float64 {
	return r.Y + r.Height
}
