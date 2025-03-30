package types

import (
	"math"
	"testing"
)

func TestNewPoint(t *testing.T) {
	p := NewPoint(3.0, 4.0)
	if p.X != 3.0 || p.Y != 4.0 {
		t.Errorf("NewPoint(3.0, 4.0) = %v; want {X:3.0, Y:4.0}", p)
	}
}

func TestPointDistanceTo(t *testing.T) {
	p1 := NewPoint(0, 0)
	p2 := NewPoint(3, 4)
	distance := p1.DistanceTo(p2)

	// Expected distance is 5.0 (3-4-5 triangle)
	if math.Abs(distance-5.0) > 1e-9 {
		t.Errorf("Distance from (0,0) to (3,4) = %v; want 5.0", distance)
	}
}

func TestPointAdd(t *testing.T) {
	p1 := NewPoint(1, 2)
	p2 := NewPoint(3, 4)
	sum := p1.Add(p2)

	if sum.X != 4.0 || sum.Y != 6.0 {
		t.Errorf("(1,2) + (3,4) = %v; want {X:4.0, Y:6.0}", sum)
	}
}

func TestPointScale(t *testing.T) {
	p := NewPoint(2, 3)
	scaled := p.Scale(2.5)

	if scaled.X != 5.0 || scaled.Y != 7.5 {
		t.Errorf("(2,3) * 2.5 = %v; want {X:5.0, Y:7.5}", scaled)
	}
}
