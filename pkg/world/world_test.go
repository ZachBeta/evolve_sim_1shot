package world

import (
	"testing"

	"github.com/zachbeta/evolve_sim/pkg/config"
	"github.com/zachbeta/evolve_sim/pkg/types"
)

func TestNewWorld(t *testing.T) {
	cfg := config.SimulationConfig{
		World: config.WorldConfig{
			Width:  100.0,
			Height: 200.0,
		},
	}

	world := NewWorld(cfg)

	if world.Width != 100.0 {
		t.Errorf("World width = %v; want 100.0", world.Width)
	}

	if world.Height != 200.0 {
		t.Errorf("World height = %v; want 200.0", world.Height)
	}

	if world.OrganismCount() != 0 {
		t.Errorf("Initial organism count = %v; want 0", world.OrganismCount())
	}

	if world.ChemicalSourceCount() != 0 {
		t.Errorf("Initial chemical source count = %v; want 0", world.ChemicalSourceCount())
	}
}

func TestWorldAddAndGetOrganisms(t *testing.T) {
	world := NewWorld(config.SimulationConfig{
		World: config.WorldConfig{Width: 100.0, Height: 100.0},
	})

	// Add an organism
	org := types.NewOrganism(
		types.NewPoint(50, 50),
		0.0,
		5.0,
		1.0,
		types.DefaultSensorAngles(),
	)

	if !world.AddOrganism(org) {
		t.Error("Failed to add organism within bounds")
	}

	// Check count
	if world.OrganismCount() != 1 {
		t.Errorf("Organism count = %v; want 1", world.OrganismCount())
	}

	// Get organisms
	organisms := world.GetOrganisms()
	if len(organisms) != 1 {
		t.Errorf("GetOrganisms returned %v organisms; want 1", len(organisms))
	}

	// Check if it's a copy by modifying it and ensuring the original is unchanged
	organisms[0].Position.X = 999

	// Get organism at index
	orgAtIndex, ok := world.GetOrganismAt(0)
	if !ok {
		t.Error("GetOrganismAt(0) returned false; want true")
	}

	if orgAtIndex.Position.X != 50 {
		t.Errorf("Original organism was modified; X = %v, want 50", orgAtIndex.Position.X)
	}

	// Test invalid index
	_, ok = world.GetOrganismAt(999)
	if ok {
		t.Error("GetOrganismAt(999) returned true; want false")
	}
}

func TestWorldUpdateOrganism(t *testing.T) {
	world := NewWorld(config.SimulationConfig{
		World: config.WorldConfig{Width: 100.0, Height: 100.0},
	})

	// Add an organism
	org := types.NewOrganism(
		types.NewPoint(50, 50),
		0.0,
		5.0,
		1.0,
		types.DefaultSensorAngles(),
	)

	world.AddOrganism(org)

	// Update the organism
	updatedOrg := org
	updatedOrg.Position.X = 60
	updatedOrg.Position.Y = 70

	success := world.UpdateOrganism(0, updatedOrg)
	if !success {
		t.Error("UpdateOrganism returned false; want true")
	}

	// Verify the update
	orgAfterUpdate, _ := world.GetOrganismAt(0)
	if orgAfterUpdate.Position.X != 60 || orgAfterUpdate.Position.Y != 70 {
		t.Errorf("Organism position after update = (%v, %v); want (60, 70)",
			orgAfterUpdate.Position.X, orgAfterUpdate.Position.Y)
	}

	// Test updating invalid index
	success = world.UpdateOrganism(999, org)
	if success {
		t.Error("UpdateOrganism(999) returned true; want false")
	}

	// Test updating with out-of-bounds position
	outOfBoundsOrg := org
	outOfBoundsOrg.Position.X = 999
	success = world.UpdateOrganism(0, outOfBoundsOrg)
	if success {
		t.Error("UpdateOrganism with out-of-bounds position returned true; want false")
	}
}

func TestWorldAddAndGetChemicalSources(t *testing.T) {
	world := NewWorld(config.SimulationConfig{
		World: config.WorldConfig{Width: 100.0, Height: 100.0},
	})

	// Add a chemical source
	source := types.NewChemicalSource(
		types.NewPoint(50, 50),
		100.0,
		0.1,
	)

	if !world.AddChemicalSource(source) {
		t.Error("Failed to add chemical source within bounds")
	}

	// Check count
	if world.ChemicalSourceCount() != 1 {
		t.Errorf("Chemical source count = %v; want 1", world.ChemicalSourceCount())
	}

	// Get chemical sources
	sources := world.GetChemicalSources()
	if len(sources) != 1 {
		t.Errorf("GetChemicalSources returned %v sources; want 1", len(sources))
	}

	// Check if it's a copy by modifying it and ensuring the original is unchanged
	sources[0].Position.X = 999

	// Check original source is unchanged
	if world.ChemicalSources[0].Position.X != 50 {
		t.Errorf("Original source was modified; X = %v, want 50", world.ChemicalSources[0].Position.X)
	}
}

func TestWorldGetConcentrationAt(t *testing.T) {
	world := NewWorld(config.SimulationConfig{
		World: config.WorldConfig{Width: 100.0, Height: 100.0},
	})

	// Add two chemical sources
	source1 := types.NewChemicalSource(types.NewPoint(25, 25), 100.0, 0.1)
	source2 := types.NewChemicalSource(types.NewPoint(75, 75), 50.0, 0.2)
	world.AddChemicalSource(source1)
	world.AddChemicalSource(source2)

	// Test concentration at various points
	testPoints := []types.Point{
		{X: 25, Y: 25}, // At source 1
		{X: 75, Y: 75}, // At source 2
		{X: 50, Y: 50}, // Between sources
		{X: 0, Y: 0},   // Far from sources
	}

	for _, point := range testPoints {
		// Direct calculation
		concentration := world.GetConcentrationAt(point)

		// Expected concentration (sum of contributions from each source)
		expected := source1.GetConcentrationAt(point) + source2.GetConcentrationAt(point)

		if !approximatelyEqual(concentration, expected, 1e-9) {
			t.Errorf("Concentration at (%v, %v) = %v; want %v",
				point.X, point.Y, concentration, expected)
		}
	}
}

func TestWorldGradientCalculation(t *testing.T) {
	world := NewWorld(config.SimulationConfig{
		World: config.WorldConfig{Width: 100.0, Height: 100.0},
	})

	// Add a single chemical source at the center
	source := types.NewChemicalSource(types.NewPoint(50, 50), 100.0, 0.01)
	world.AddChemicalSource(source)

	// Test gradient at various points
	testPoints := []struct {
		point             types.Point
		expectedDirection types.Point // approximate expected direction
	}{
		{types.Point{X: 30, Y: 50}, types.Point{X: 1, Y: 0}},  // Left of source, should point right
		{types.Point{X: 70, Y: 50}, types.Point{X: -1, Y: 0}}, // Right of source, should point left
		{types.Point{X: 50, Y: 30}, types.Point{X: 0, Y: 1}},  // Above source, should point down
		{types.Point{X: 50, Y: 70}, types.Point{X: 0, Y: -1}}, // Below source, should point up
	}

	for _, tc := range testPoints {
		gradient := world.GetConcentrationGradientAt(tc.point)

		// Check if the gradient direction is approximately correct
		// For this test, we just check the sign of x and y components
		if !sameSign(gradient.X, tc.expectedDirection.X) || !sameSign(gradient.Y, tc.expectedDirection.Y) {
			t.Errorf("Gradient at (%v, %v) = (%v, %v); expected direction (%v, %v)",
				tc.point.X, tc.point.Y, gradient.X, gradient.Y,
				tc.expectedDirection.X, tc.expectedDirection.Y)
		}
	}
}

func TestConcentrationGrid(t *testing.T) {
	world := NewWorld(config.SimulationConfig{
		World: config.WorldConfig{Width: 100.0, Height: 100.0},
	})

	// Add a chemical source
	source := types.NewChemicalSource(types.NewPoint(50, 50), 100.0, 0.01)
	world.AddChemicalSource(source)

	// Initialize the concentration grid
	world.InitializeConcentrationGrid(5.0) // 5.0 units per cell

	// Test concentration at various points using the grid
	testPoints := []types.Point{
		{X: 25, Y: 25},
		{X: 50, Y: 50},
		{X: 75, Y: 75},
		{X: 10, Y: 10},
	}

	for _, point := range testPoints {
		// Direct calculation
		directConcentration := source.GetConcentrationAt(point)

		// Grid-based calculation
		gridConcentration := world.GetConcentrationAt(point)

		// Allow some error due to grid discretization
		if !approximatelyEqual(gridConcentration, directConcentration, 0.5) {
			t.Errorf("Grid concentration at (%v, %v) = %v; direct calculation = %v",
				point.X, point.Y, gridConcentration, directConcentration)
		}
	}
}

// Helper function to check if two float64 values are approximately equal
func approximatelyEqual(a, b, epsilon float64) bool {
	diff := a - b
	return diff < epsilon && diff > -epsilon
}

// Helper function to check if two float64 values have the same sign
func sameSign(a, b float64) bool {
	if a == 0 || b == 0 {
		return true // Treat zero as matching any sign
	}
	return (a > 0) == (b > 0)
}
