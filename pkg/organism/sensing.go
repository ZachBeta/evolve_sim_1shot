package organism

import (
	"github.com/zachbeta/evolve_sim/pkg/types"
)

// SensorReadings represents the chemical concentration readings from the three sensors
type SensorReadings struct {
	Front float64
	Left  float64
	Right float64
}

// ReadSensors reads the chemical concentration at each sensor position
// Returns the concentration readings for the front, left, and right sensors
func ReadSensors(
	org *types.Organism,
	world interface{ GetConcentrationAt(types.Point) float64 },
	sensorDistance float64,
) SensorReadings {
	// Get sensor positions
	sensorPositions := org.GetSensorPositions(sensorDistance)

	// Read concentrations at each sensor position
	readings := SensorReadings{
		Front: world.GetConcentrationAt(sensorPositions[0]),
		Left:  world.GetConcentrationAt(sensorPositions[1]),
		Right: world.GetConcentrationAt(sensorPositions[2]),
	}

	return readings
}
