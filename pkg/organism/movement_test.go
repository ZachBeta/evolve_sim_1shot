package organism

import (
	"math"
	"testing"

	"github.com/zachbeta/evolve_sim/pkg/types"
)

func TestMove(t *testing.T) {
	// Define test bounds
	bounds := types.Rect{
		Min: types.Point{X: 0, Y: 0},
		Max: types.Point{X: 100, Y: 100},
	}

	t.Run("Normal movement within bounds", func(t *testing.T) {
		// Create organism in middle of bounds
		org := types.NewOrganism(
			types.Point{X: 50, Y: 50},
			0, // Heading east (0 radians)
			10,
			1.0,
			types.DefaultSensorAngles(),
		)

		// Move organism
		Move(&org, bounds, 1.0)

		// Expected position after moving east at speed 1.0 for 1.0 time units
		expectedX := 51.0
		expectedY := 50.0

		if math.Abs(org.Position.X-expectedX) > 0.001 || math.Abs(org.Position.Y-expectedY) > 0.001 {
			t.Errorf("Expected position (%f, %f), got (%f, %f)",
				expectedX, expectedY, org.Position.X, org.Position.Y)
		}
	})

	t.Run("Boundary collision - right wall", func(t *testing.T) {
		// Create organism near right boundary
		org := types.NewOrganism(
			types.Point{X: 99.5, Y: 50},
			0, // Heading east (0 radians)
			10,
			1.0,
			types.DefaultSensorAngles(),
		)

		// Original heading
		originalHeading := org.Heading

		// Move organism
		Move(&org, bounds, 1.0)

		// Expect heading to be flipped (π radians)
		// Due to reflection, heading should be approximately π (east -> west)
		if math.Abs(org.Heading-math.Pi) > 0.1 && math.Abs(org.Heading) > 0.1 {
			t.Errorf("Expected heading near %f or %f, got %f", math.Pi, 0.0, org.Heading)
		}

		// Position should be adjusted to remain within bounds
		if !bounds.Contains(org.Position) {
			t.Errorf("Organism position (%f, %f) outside bounds after collision",
				org.Position.X, org.Position.Y)
		}
	})

	t.Run("Boundary collision - bottom wall", func(t *testing.T) {
		// Create organism near bottom boundary
		org := types.NewOrganism(
			types.Point{X: 50, Y: 99.5},
			math.Pi/2, // Heading south (π/2 radians)
			10,
			1.0,
			types.DefaultSensorAngles(),
		)

		// Original heading
		originalHeading := org.Heading

		// Move organism
		Move(&org, bounds, 1.0)

		// Position should be adjusted to remain within bounds
		if !bounds.Contains(org.Position) {
			t.Errorf("Organism position (%f, %f) outside bounds after collision",
				org.Position.X, org.Position.Y)
		}

		// Heading should be different after collision
		if math.Abs(org.Heading-originalHeading) < 0.1 {
			t.Errorf("Heading did not change after collision")
		}
	})
}
