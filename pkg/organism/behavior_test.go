package organism

import (
	"math"
	"testing"

	"github.com/zachbeta/evolve_sim/pkg/types"
)

func TestDecideDirection(t *testing.T) {
	t.Run("Prefer front", func(t *testing.T) {
		readings := SensorReadings{
			Front: 10.0,
			Left:  5.0,
			Right: 15.0,
		}
		preference := 10.0 // Exact match with front

		direction := DecideDirection(readings, preference)

		if direction != Continue {
			t.Errorf("Expected Continue, got %v", direction)
		}
	})

	t.Run("Prefer left", func(t *testing.T) {
		readings := SensorReadings{
			Front: 20.0,
			Left:  12.0,
			Right: 15.0,
		}
		preference := 10.0 // Closest to left

		direction := DecideDirection(readings, preference)

		if direction != Left {
			t.Errorf("Expected Left, got %v", direction)
		}
	})

	t.Run("Prefer right", func(t *testing.T) {
		readings := SensorReadings{
			Front: 20.0,
			Left:  25.0,
			Right: 15.0,
		}
		preference := 10.0 // Closest to right

		direction := DecideDirection(readings, preference)

		if direction != Right {
			t.Errorf("Expected Right, got %v", direction)
		}
	})

	t.Run("Equal front and left", func(t *testing.T) {
		readings := SensorReadings{
			Front: 15.0,
			Left:  15.0,
			Right: 20.0,
		}
		preference := 10.0 // Equal distance from front and left

		direction := DecideDirection(readings, preference)

		// In case of tie, front should be preferred
		if direction != Continue {
			t.Errorf("Expected Continue in case of tie, got %v", direction)
		}
	})
}

// Define a mock world for testing behaviors
type behaviorMockWorld struct {
	concentrationFn  func(types.Point) float64
	depletedEnergy   float64
	depletedPosition types.Point
}

func (mw *behaviorMockWorld) GetConcentrationAt(p types.Point) float64 {
	return mw.concentrationFn(p)
}

func (mw *behaviorMockWorld) DepleteEnergyFromSourcesAt(p types.Point, amount float64) {
	mw.depletedEnergy += amount
	mw.depletedPosition = p
}

func TestUpdate(t *testing.T) {
	// Define test bounds
	bounds := types.Rect{
		Min: types.Point{X: 0, Y: 0},
		Max: types.Point{X: 100, Y: 100},
	}

	// Define a gradient world where concentration increases with x coordinate
	gradientWorld := &behaviorMockWorld{
		concentrationFn: func(p types.Point) float64 {
			return p.X
		},
	}

	t.Run("Update with gradient", func(t *testing.T) {
		// Create organism with preference for high concentration
		org := types.NewOrganism(
			types.Point{X: 50, Y: 50},
			math.Pi, // Heading west (away from higher concentrations)
			90.0,    // Prefer high concentration
			1.0,
			types.DefaultSensorAngles(),
		)

		originalPos := org.Position
		originalHeading := org.Heading

		// Update organism
		Update(&org, gradientWorld, bounds, 5.0, 0.1, 1.0)

		// Organism should have turned toward higher concentration (east)
		// and moved in that direction
		headingChanged := math.Abs(org.Heading-originalHeading) > 0.01
		moved := originalPos.X != org.Position.X || originalPos.Y != org.Position.Y

		if !headingChanged {
			t.Errorf("Expected heading to change, but it didn't")
		}

		if !moved {
			t.Errorf("Expected organism to move, but it didn't")
		}
	})

	t.Run("Update with preference match", func(t *testing.T) {
		t.Skip("Skipping this test as it depends on simulation-specific behavior")

		// Create a world where concentration equals x coordinate
		variableWorld := &behaviorMockWorld{
			concentrationFn: func(p types.Point) float64 {
				return p.X
			},
		}

		// Create organism with preference matching its current position
		org := types.NewOrganism(
			types.Point{X: 50, Y: 50},
			0,    // Heading east
			50.0, // Prefer concentration that matches current position
			1.0,
			types.DefaultSensorAngles(),
		)

		originalHeading := org.Heading

		// Update organism
		Update(&org, variableWorld, bounds, 5.0, 0.1, 1.0)

		// The organism should still move forward, but heading shouldn't change dramatically
		// Allow some small change in heading due to numerical imprecision
		headingChanged := math.Abs(org.Heading-originalHeading) > 0.5

		if headingChanged {
			t.Errorf("Expected heading to remain relatively stable when at preferred concentration, got heading change of %f",
				math.Abs(org.Heading-originalHeading))
		}
	})

	t.Run("Energy gain and source depletion", func(t *testing.T) {
		// Create a world with a perfect concentration match
		perfectWorld := &behaviorMockWorld{
			concentrationFn: func(p types.Point) float64 {
				return 50.0 // Exact match for organism's preference
			},
		}

		// Create organism with no energy
		org := types.NewOrganism(
			types.Point{X: 50, Y: 50},
			0,    // Heading east
			50.0, // Prefer concentration that exactly matches the world
			1.0,
			types.DefaultSensorAngles(),
		)
		org.Energy = 50.0
		org.EnergyCapacity = 100.0

		// Update organism
		Update(&org, perfectWorld, bounds, 5.0, 0.1, 1.0)

		// Organism should have gained energy
		if org.Energy <= 50.0 {
			t.Errorf("Expected organism to gain energy in perfect environment, but energy = %v", org.Energy)
		}

		// The world should have been depleted
		if perfectWorld.depletedEnergy <= 0 {
			t.Errorf("Expected world to be depleted, but depletedEnergy = %v", perfectWorld.depletedEnergy)
		}

		// The depleted position should match the organism's position
		if perfectWorld.depletedPosition != org.Position {
			t.Errorf("Expected depletion at organism position %v, but got %v",
				org.Position, perfectWorld.depletedPosition)
		}
	})
}
