package organism

import (
	"testing"

	"github.com/zachbeta/evolve_sim/pkg/types"
)

// mockWorld implements a simple world that returns predefined concentrations
type mockWorld struct {
	// Define a function that returns concentration at given point
	concentrationFn func(types.Point) float64
}

func (m mockWorld) GetConcentrationAt(p types.Point) float64 {
	return m.concentrationFn(p)
}

func TestReadSensors(t *testing.T) {
	// Define a constant concentration world for basic testing
	constantWorld := mockWorld{
		concentrationFn: func(p types.Point) float64 {
			return 10.0
		},
	}

	// Define a gradient world where concentration increases with x coordinate
	gradientWorld := mockWorld{
		concentrationFn: func(p types.Point) float64 {
			return p.X
		},
	}

	t.Run("Read constant concentration", func(t *testing.T) {
		// Create organism
		org := types.NewOrganism(
			types.Point{X: 50, Y: 50},
			0, // Heading east
			10,
			1.0,
			types.DefaultSensorAngles(),
		)

		// Read sensors with constant world
		readings := ReadSensors(&org, constantWorld, 5.0)

		// All readings should be 10.0
		if readings.Front != 10.0 || readings.Left != 10.0 || readings.Right != 10.0 {
			t.Errorf("Expected all readings to be 10.0, got Front: %f, Left: %f, Right: %f",
				readings.Front, readings.Left, readings.Right)
		}
	})

	t.Run("Read gradient concentration", func(t *testing.T) {
		// Create organism
		org := types.NewOrganism(
			types.Point{X: 50, Y: 50},
			0, // Heading east
			10,
			1.0,
			types.DefaultSensorAngles(),
		)

		// Read sensors with gradient world
		readings := ReadSensors(&org, gradientWorld, 5.0)

		// Front sensor should read higher concentration than left and right
		if readings.Front <= readings.Left || readings.Front <= readings.Right {
			t.Errorf("Expected front reading (%f) to be higher than left (%f) and right (%f)",
				readings.Front, readings.Left, readings.Right)
		}

		// Left and right readings should be approximately equal
		if readings.Left != readings.Right {
			t.Errorf("Expected left and right readings to be equal, got left: %f, right: %f",
				readings.Left, readings.Right)
		}
	})
}
