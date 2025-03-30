package types

import "testing"

func TestNewRect(t *testing.T) {
	r := NewRect(10, 20, 30, 40)
	if r.X != 10 || r.Y != 20 || r.Width != 30 || r.Height != 40 {
		t.Errorf("NewRect(10, 20, 30, 40) = %v; want {X:10, Y:20, Width:30, Height:40}", r)
	}
}

func TestRectContains(t *testing.T) {
	r := NewRect(10, 10, 20, 20)

	// Test points inside
	insidePoints := []Point{
		{X: 10, Y: 10}, // Top-left corner (inclusive)
		{X: 29, Y: 10}, // Top-right (exclusive on right edge)
		{X: 10, Y: 29}, // Bottom-left (exclusive on bottom edge)
		{X: 20, Y: 20}, // Middle
	}

	for i, p := range insidePoints {
		if !r.Contains(p) {
			t.Errorf("Point %d %v should be inside rect %v", i, p, r)
		}
	}

	// Test points outside
	outsidePoints := []Point{
		{X: 9, Y: 10},  // Left of left edge
		{X: 30, Y: 10}, // Right of right edge
		{X: 10, Y: 9},  // Above top edge
		{X: 10, Y: 30}, // Below bottom edge
		{X: 30, Y: 30}, // Outside corner
	}

	for i, p := range outsidePoints {
		if r.Contains(p) {
			t.Errorf("Point %d %v should be outside rect %v", i, p, r)
		}
	}
}

func TestRectCenter(t *testing.T) {
	r := NewRect(10, 20, 30, 40)
	center := r.Center()

	if center.X != 25 || center.Y != 40 {
		t.Errorf("Center of rect %v = %v; want {X:25, Y:40}", r, center)
	}
}

func TestRectEdges(t *testing.T) {
	r := NewRect(10, 20, 30, 40)

	if r.GetX() != 10 {
		t.Errorf("GetX() = %v; want 10", r.GetX())
	}

	if r.GetY() != 20 {
		t.Errorf("GetY() = %v; want 20", r.GetY())
	}

	if r.GetMaxX() != 40 {
		t.Errorf("GetMaxX() = %v; want 40", r.GetMaxX())
	}

	if r.GetMaxY() != 60 {
		t.Errorf("GetMaxY() = %v; want 60", r.GetMaxY())
	}
}
