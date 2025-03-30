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
	world interface{ GetConcentrationAt(types.Point) float64 },
	bounds types.Rect,
	sensorDistance float64,
	turnSpeed float64,
	deltaTime float64,
) {
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

	// Move forward (this now includes energy consumption)
	Move(org, bounds, deltaTime)

	// Energy gain from being in preferred environment
	currentConcentration := world.GetConcentrationAt(org.Position)
	concentrationDiff := math.Abs(currentConcentration - org.ChemPreference)

	// Calculate how close the organism is to its preferred concentration (0-1 scale)
	// 1.0 means perfect match, 0.0 means furthest possible
	concentrationFit := 1.0 - math.Min(concentrationDiff/(org.ChemPreference+1.0), 1.0)

	if concentrationFit > ENERGY_GAIN_THRESHOLD {
		// Scale energy gain by how good the fit is
		gainFactor := (concentrationFit - ENERGY_GAIN_THRESHOLD) / (1.0 - ENERGY_GAIN_THRESHOLD)
		energyGain := gainFactor * MAX_ENERGY_GAIN * deltaTime

		// Add energy, capped at max capacity
		org.Energy = math.Min(org.Energy+energyGain, org.EnergyCapacity)
	}
}
