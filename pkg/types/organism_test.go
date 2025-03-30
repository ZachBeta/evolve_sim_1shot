package types

import (
	"math"
	"testing"
)

func TestNewOrganism(t *testing.T) {
	position := NewPoint(10, 20)
	heading := math.Pi / 2 // 90 degrees
	chemPreference := 5.0
	speed := 2.0
	sensorAngles := DefaultSensorAngles()

	org := NewOrganism(position, heading, chemPreference, speed, sensorAngles)

	if org.Position.X != 10 || org.Position.Y != 20 {
		t.Errorf("Organism position = %v; want {X:10, Y:20}", org.Position)
	}

	if org.Heading != math.Pi/2 {
		t.Errorf("Organism heading = %v; want %v", org.Heading, math.Pi/2)
	}

	if org.ChemPreference != 5.0 {
		t.Errorf("Organism chemPreference = %v; want 5.0", org.ChemPreference)
	}

	if org.Speed != 2.0 {
		t.Errorf("Organism speed = %v; want 2.0", org.Speed)
	}

	if org.SensorAngles != sensorAngles {
		t.Errorf("Organism sensorAngles = %v; want %v", org.SensorAngles, sensorAngles)
	}
}

func TestDefaultSensorAngles(t *testing.T) {
	angles := DefaultSensorAngles()

	if angles[0] != 0 {
		t.Errorf("Default front sensor angle = %v; want 0", angles[0])
	}

	if angles[1] != -math.Pi/4 {
		t.Errorf("Default left sensor angle = %v; want %v", angles[1], -math.Pi/4)
	}

	if angles[2] != math.Pi/4 {
		t.Errorf("Default right sensor angle = %v; want %v", angles[2], math.Pi/4)
	}
}

func TestGetSensorPositions(t *testing.T) {
	// Create organism at origin facing right (0 radians)
	org := NewOrganism(NewPoint(0, 0), 0, 5.0, 1.0, DefaultSensorAngles())
	sensorDistance := 2.0

	positions := org.GetSensorPositions(sensorDistance)

	// Front sensor should be at (2, 0)
	if math.Abs(positions[0].X-2.0) > 1e-9 || math.Abs(positions[0].Y) > 1e-9 {
		t.Errorf("Front sensor position = %v; want approx {X:2, Y:0}", positions[0])
	}

	// Left sensor should be at approximately (√2, -√2)
	expectedLeftX := math.Cos(-math.Pi/4) * sensorDistance
	expectedLeftY := math.Sin(-math.Pi/4) * sensorDistance
	if math.Abs(positions[1].X-expectedLeftX) > 1e-9 ||
		math.Abs(positions[1].Y-expectedLeftY) > 1e-9 {
		t.Errorf("Left sensor position = %v; want approx {X:%v, Y:%v}",
			positions[1], expectedLeftX, expectedLeftY)
	}

	// Right sensor should be at approximately (√2, √2)
	expectedRightX := math.Cos(math.Pi/4) * sensorDistance
	expectedRightY := math.Sin(math.Pi/4) * sensorDistance
	if math.Abs(positions[2].X-expectedRightX) > 1e-9 ||
		math.Abs(positions[2].Y-expectedRightY) > 1e-9 {
		t.Errorf("Right sensor position = %v; want approx {X:%v, Y:%v}",
			positions[2], expectedRightX, expectedRightY)
	}
}

func TestMoveForward(t *testing.T) {
	// Test moving right (0 radians)
	org1 := NewOrganism(NewPoint(0, 0), 0, 5.0, 1.0, DefaultSensorAngles())
	org1.MoveForward(5.0)
	if math.Abs(org1.Position.X-5.0) > 1e-9 || math.Abs(org1.Position.Y) > 1e-9 {
		t.Errorf("After moving right, position = %v; want {X:5, Y:0}", org1.Position)
	}

	// Test moving up (π/2 radians)
	org2 := NewOrganism(NewPoint(0, 0), math.Pi/2, 5.0, 1.0, DefaultSensorAngles())
	org2.MoveForward(5.0)
	if math.Abs(org2.Position.X) > 1e-9 || math.Abs(org2.Position.Y-5.0) > 1e-9 {
		t.Errorf("After moving up, position = %v; want {X:0, Y:5}", org2.Position)
	}

	// Test moving at 45 degrees (π/4 radians)
	org3 := NewOrganism(NewPoint(0, 0), math.Pi/4, 5.0, 1.0, DefaultSensorAngles())
	org3.MoveForward(math.Sqrt(2))
	if math.Abs(org3.Position.X-1.0) > 1e-9 || math.Abs(org3.Position.Y-1.0) > 1e-9 {
		t.Errorf("After moving diagonally, position = %v; want approx {X:1, Y:1}", org3.Position)
	}
}

func TestTurn(t *testing.T) {
	// Test turning right
	org1 := NewOrganism(NewPoint(0, 0), 0, 5.0, 1.0, DefaultSensorAngles())
	org1.Turn(math.Pi / 2)
	if math.Abs(org1.Heading-math.Pi/2) > 1e-9 {
		t.Errorf("After turning right, heading = %v; want %v", org1.Heading, math.Pi/2)
	}

	// Test turning left
	org2 := NewOrganism(NewPoint(0, 0), math.Pi, 5.0, 1.0, DefaultSensorAngles())
	org2.Turn(-math.Pi / 2)
	if math.Abs(org2.Heading-math.Pi/2) > 1e-9 {
		t.Errorf("After turning left, heading = %v; want %v", org2.Heading, math.Pi/2)
	}

	// Test normalization of heading (> 2π)
	org3 := NewOrganism(NewPoint(0, 0), 0, 5.0, 1.0, DefaultSensorAngles())
	org3.Turn(3 * math.Pi)
	expected3 := math.Mod(3*math.Pi, 2*math.Pi)
	if math.Abs(org3.Heading-expected3) > 1e-9 {
		t.Errorf("After turning beyond 2π, heading = %v; want %v", org3.Heading, expected3)
	}

	// Test normalization of heading (< 0)
	org4 := NewOrganism(NewPoint(0, 0), math.Pi, 5.0, 1.0, DefaultSensorAngles())
	org4.Turn(-3 * math.Pi)
	// With our implementation, negative angles wrap to positive values in [0, 2π)
	expected4 := math.Mod(-3*math.Pi+math.Pi, 2*math.Pi)
	if expected4 < 0 {
		expected4 += 2 * math.Pi
	}
	if math.Abs(org4.Heading-expected4) > 1e-9 {
		t.Errorf("After turning below 0, heading = %v; want %v", org4.Heading, expected4)
	}
}
