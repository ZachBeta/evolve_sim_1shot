package organism

import (
	"math"

	"github.com/zachbeta/evolve_sim/pkg/types"
)

// Energy system constants
const (
	ENERGY_GAIN_THRESHOLD = 0.7  // Minimum concentration fit to gain energy (0-1)
	MAX_ENERGY_GAIN       = 0.5  // Maximum energy gain per second
	MAX_CONCENTRATION     = 1000 // Maximum expected concentration for normalization
)

// Direction represents the three possible directions an organism can turn
type Direction int

const (
	// Continue moving in the same direction
	Continue Direction = iota
	// Turn left
	Left
	// Turn right
	Right
)

// DecideDirection determines the best direction for the organism to move
// based on its sensor readings and chemical preference
func DecideDirection(readings SensorReadings, preference float64) Direction {
	// Calculate the difference between each reading and the preference
	// We want to find the reading closest to the preference
	frontDiff := math.Abs(readings.Front - preference)
	leftDiff := math.Abs(readings.Left - preference)
	rightDiff := math.Abs(readings.Right - preference)

	// Find the minimum difference
	minDiff := math.Min(frontDiff, math.Min(leftDiff, rightDiff))

	// Return the direction with the minimum difference
	if minDiff == frontDiff {
		return Continue
	} else if minDiff == leftDiff {
		return Left
	} else {
		return Right
	}
}

// Update performs a complete update cycle for an organism:
// 1. Reads sensors
// 2. Decides direction
// 3. Turns if necessary
// 4. Moves forward
// 5. Updates energy based on environment
func Update(
	org *types.Organism,
	world interface {
		GetConcentrationAt(types.Point) float64
		DepleteEnergyFromSourcesAt(types.Point, float64)
	},
	bounds types.Rect,
	sensorDistance float64,
	turnSpeed float64,
	deltaTime float64,
) {
	// Apply sensing cost before reading sensors
	org.Energy -= org.SensingCost * org.EnergyEfficiency * deltaTime

	// Read sensors
	readings := ReadSensors(org, world, sensorDistance)

	// Decide direction
	direction := DecideDirection(readings, org.ChemPreference)

	// Turn if necessary
	switch direction {
	case Left:
		org.Turn(-turnSpeed * deltaTime)
	case Right:
		org.Turn(turnSpeed * deltaTime)
	case Continue:
		// Continue straight, no turning needed
	}

	// Move forward (this includes energy consumption for movement)
	Move(org, bounds, deltaTime)

	// Update energy status - gain from optimal environment, lose from metabolism
	org.UpdateEnergy(world, deltaTime)

	// If energy is depleted, mark for removal
	if org.Energy <= 0 {
		org.MarkForRemoval = true
	}

	// Update reproduction timer
	org.TimeSinceReproduction += deltaTime
}
