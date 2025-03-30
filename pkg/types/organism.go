package types

import (
	"math"
)

// Organism represents a single-cell organism in the simulation
type Organism struct {
	Position        Point      // Current position in the world
	Heading         float64    // Direction the organism is facing (in radians)
	PreviousHeading float64    // Previous heading for smooth rotation animation
	ChemPreference  float64    // Preferred chemical concentration
	Speed           float64    // Movement speed (units per step)
	SensorAngles    [3]float64 // Angles of sensors relative to heading (front, left, right)
}

// NewOrganism creates a new organism with the given parameters
func NewOrganism(position Point, heading, chemPreference, speed float64, sensorAngles [3]float64) Organism {
	return Organism{
		Position:        position,
		Heading:         heading,
		PreviousHeading: heading, // Initialize previous heading to current heading
		ChemPreference:  chemPreference,
		Speed:           speed,
		SensorAngles:    sensorAngles,
	}
}

// DefaultSensorAngles returns the default angles for sensors: [0, -π/4, π/4]
// This corresponds to front (0°), left (-45°), and right (45°)
func DefaultSensorAngles() [3]float64 {
	return [3]float64{0, -math.Pi / 4, math.Pi / 4}
}

// GetSensorPositions calculates the positions of the organism's sensors
// based on its current position, heading, and sensor configuration
func (o Organism) GetSensorPositions(sensorDistance float64) [3]Point {
	var positions [3]Point

	for i, angle := range o.SensorAngles {
		// Calculate absolute angle by adding sensor angle to heading
		absoluteAngle := o.Heading + angle

		// Calculate sensor offset using trigonometry
		dx := math.Cos(absoluteAngle) * sensorDistance
		dy := math.Sin(absoluteAngle) * sensorDistance

		// Calculate sensor position
		positions[i] = Point{
			X: o.Position.X + dx,
			Y: o.Position.Y + dy,
		}
	}

	return positions
}

// MoveForward moves the organism forward in its current heading direction
func (o *Organism) MoveForward(distance float64) {
	dx := math.Cos(o.Heading) * distance
	dy := math.Sin(o.Heading) * distance

	o.Position.X += dx
	o.Position.Y += dy
}

// Turn changes the organism's heading by the specified angle (in radians)
func (o *Organism) Turn(angle float64) {
	o.Heading += angle

	// Normalize heading to [0, 2π)
	o.Heading = math.Mod(o.Heading, 2*math.Pi)
	if o.Heading < 0 {
		o.Heading += 2 * math.Pi
	}
}
