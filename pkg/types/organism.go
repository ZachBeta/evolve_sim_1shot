package types

import (
	"math"
	"math/rand"
)

// MaxTrailLength defines the maximum number of positions to store in the trail
const MaxTrailLength = 30

// Constants for reproduction
const (
	ReproductionThreshold = 0.75 // Percentage of max energy required to reproduce
	ReproductionCooldown  = 5.0  // Seconds between reproduction attempts
	OffspringEnergyRatio  = 0.3  // Portion of parent's energy given to offspring
	MutationFactorSmall   = 0.05 // For small mutations (like preferences)
	MutationFactorMedium  = 0.1  // For medium mutations (like speed)
	MutationFactorLarge   = 0.2  // For large mutations (like sensor distance)
)

// Organism represents a single-cell organism in the simulation
type Organism struct {
	Position              Point      // Current position in the world
	Heading               float64    // Direction the organism is facing (in radians)
	PreviousHeading       float64    // Previous heading for smooth rotation animation
	ChemPreference        float64    // Preferred chemical concentration
	Speed                 float64    // Movement speed (units per step)
	SensorAngles          [3]float64 // Angles of sensors relative to heading (front, left, right)
	PositionHistory       []Point    // History of positions for drawing trails
	UpdateCounter         int        // Counter to control how often we record position
	Energy                float64    // Current energy level
	EnergyCapacity        float64    // Maximum energy capacity
	TimeSinceReproduction float64    // Time elapsed since last reproduction
}

// NewOrganism creates a new organism with the given parameters
func NewOrganism(position Point, heading, chemPreference, speed float64, sensorAngles [3]float64) Organism {
	// Default energy capacity based on size/speed
	energyCapacity := 100.0 + speed*10.0

	return Organism{
		Position:              position,
		Heading:               heading,
		PreviousHeading:       heading, // Initialize previous heading to current heading
		ChemPreference:        chemPreference,
		Speed:                 speed,
		SensorAngles:          sensorAngles,
		PositionHistory:       make([]Point, 0, MaxTrailLength),
		UpdateCounter:         0,
		Energy:                energyCapacity * 0.8, // Start with 80% of max energy
		EnergyCapacity:        energyCapacity,
		TimeSinceReproduction: 0,
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

// UpdateTrail adds the current position to the position history
// if enough movement has occurred since the last update
func (o *Organism) UpdateTrail() {
	// Only update every few frames to avoid too many points
	o.UpdateCounter++
	if o.UpdateCounter >= 5 { // Record every 5th update
		o.UpdateCounter = 0

		// Add current position to history
		o.PositionHistory = append(o.PositionHistory, o.Position)

		// Trim history if it exceeds max length
		if len(o.PositionHistory) > MaxTrailLength {
			o.PositionHistory = o.PositionHistory[1:]
		}
	}
}

// CanReproduce checks if the organism has enough energy and has waited the cooldown period
func (o *Organism) CanReproduce() bool {
	return o.Energy >= o.EnergyCapacity*ReproductionThreshold &&
		o.TimeSinceReproduction >= ReproductionCooldown
}

// Reproduce creates a new organism with slight mutations
// The parent loses some energy in the process
func (o *Organism) Reproduce() Organism {
	// Calculate how much energy to give the offspring
	offspringEnergy := o.Energy * OffspringEnergyRatio

	// Reduce parent's energy
	o.Energy -= offspringEnergy

	// Reset reproduction timer
	o.TimeSinceReproduction = 0

	// Create offspring with mutations
	// Position is set to be slightly offset from parent
	offsetDistance := 5.0 + rand.Float64()*5.0  // 5-10 units away
	offsetAngle := rand.Float64() * 2 * math.Pi // Random angle

	positionOffset := Point{
		X: math.Cos(offsetAngle) * offsetDistance,
		Y: math.Sin(offsetAngle) * offsetDistance,
	}

	offspringPosition := Point{
		X: o.Position.X + positionOffset.X,
		Y: o.Position.Y + positionOffset.Y,
	}

	// Apply small mutations to preferences and attributes
	// Using normal distribution for more realistic mutations
	prefMutation := rand.NormFloat64() * o.ChemPreference * MutationFactorSmall
	speedMutation := rand.NormFloat64() * o.Speed * MutationFactorMedium

	// Don't allow negative speed
	newSpeed := math.Max(0.1, o.Speed+speedMutation)

	// Random heading for the offspring
	newHeading := rand.Float64() * 2 * math.Pi

	// Slightly mutate sensor angles
	var newSensorAngles [3]float64
	for i, angle := range o.SensorAngles {
		mutation := rand.NormFloat64() * MutationFactorSmall
		newSensorAngles[i] = angle + mutation
	}

	// Calculate new energy capacity based on speed
	newEnergyCapacity := 100.0 + newSpeed*10.0

	// Create the offspring
	return Organism{
		Position:              offspringPosition,
		Heading:               newHeading,
		PreviousHeading:       newHeading,
		ChemPreference:        o.ChemPreference + prefMutation,
		Speed:                 newSpeed,
		SensorAngles:          newSensorAngles,
		PositionHistory:       make([]Point, 0, MaxTrailLength),
		UpdateCounter:         0,
		Energy:                offspringEnergy,
		EnergyCapacity:        newEnergyCapacity,
		TimeSinceReproduction: 0,
	}
}
