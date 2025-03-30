package world

import (
	"math/rand"
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

func TestDepleteEnergyFromSourcesAt(t *testing.T) {
	// Create a test world with a few chemical sources
	world := setupTestWorld()

	// Create chemical sources with known properties
	source1 := types.NewChemicalSource(types.Point{X: 50, Y: 50}, 100, 0.01)
	source2 := types.NewChemicalSource(types.Point{X: 150, Y: 150}, 200, 0.005)

	world.AddChemicalSource(source1)
	world.AddChemicalSource(source2)

	// Get the chemical sources from the world to record initial state
	sources := world.GetChemicalSources()

	// Store initial energy levels
	initialEnergies := make([]float64, len(sources))
	for i, source := range sources {
		initialEnergies[i] = source.Energy
	}

	// Test position that should be closer to source1
	testPosition := types.Point{X: 60, Y: 60}

	// Deplete some energy at this position
	depletionAmount := 100.0
	world.DepleteEnergyFromSourcesAt(testPosition, depletionAmount)

	// Get updated sources
	updatedSources := world.GetChemicalSources()

	// Verify energy was depleted from both sources but more from source1
	depletedEnergy1 := initialEnergies[len(initialEnergies)-2] - updatedSources[len(updatedSources)-2].Energy
	depletedEnergy2 := initialEnergies[len(initialEnergies)-1] - updatedSources[len(updatedSources)-1].Energy

	// Source1 should have lost more energy than source2
	if depletedEnergy1 <= depletedEnergy2 {
		t.Errorf("Expected source1 to lose more energy than source2, but got: source1=%v, source2=%v",
			depletedEnergy1, depletedEnergy2)
	}

	// Total depletion should be proportional to depletionAmount
	totalDepleted := depletedEnergy1 + depletedEnergy2
	// Account for the multiplier in the depletion method
	expectedTotal := depletionAmount * 2.0

	// Use a reasonable tolerance since floating point calculations are involved
	tolerance := 0.01 * expectedTotal
	if totalDepleted < expectedTotal-tolerance || totalDepleted > expectedTotal+tolerance {
		t.Errorf("Total depleted energy %v doesn't match expected %v within tolerance %v",
			totalDepleted, expectedTotal, tolerance)
	}
}

func TestDepleteEnergySourceDeactivation(t *testing.T) {
	// Create a test world with a chemical source
	world := setupTestWorld()

	// Create a chemical source with low energy
	lowEnergySource := types.NewChemicalSource(types.Point{X: 50, Y: 50}, 100, 0.01)
	lowEnergySource.Energy = 10.0 // Override with a low energy value

	world.AddChemicalSource(lowEnergySource)

	// Test position near the source
	testPosition := types.Point{X: 50, Y: 55}

	// Deplete more energy than the source has
	depletionAmount := 100.0
	world.DepleteEnergyFromSourcesAt(testPosition, depletionAmount)

	// Get updated sources
	updatedSources := world.GetChemicalSources()
	lastSource := updatedSources[len(updatedSources)-1]

	// Verify the source is depleted and inactive
	if lastSource.Energy != 0 {
		t.Errorf("Expected source energy to be 0, but got %v", lastSource.Energy)
	}

	if lastSource.IsActive {
		t.Error("Expected source to be inactive after full depletion")
	}

	// Verify that concentration at the test position is now 0
	concentration := world.GetConcentrationAt(testPosition)
	if concentration > 0 {
		t.Errorf("Expected concentration to be 0 after source deactivation, but got %v", concentration)
	}
}

func TestSystemEnergyTracking(t *testing.T) {
	// Create a world with a known configuration
	cfg := config.SimulationConfig{
		World: config.WorldConfig{
			Width:  1000,
			Height: 1000,
		},
		Chemical: config.ChemicalConfig{
			Count:              5,
			MinStrength:        100,
			MaxStrength:        200,
			MinDecayFactor:     0.001,
			MaxDecayFactor:     0.01,
			DepletionRate:      0.2,
			TargetSystemEnergy: 50000, // Explicitly set target energy
		},
	}

	world := NewWorld(cfg)

	// Get initial energy values
	totalEnergy, targetEnergy := world.GetSystemEnergyInfo()

	// Target energy should match configuration
	if targetEnergy != cfg.Chemical.TargetSystemEnergy {
		t.Errorf("Target energy = %v; want %v", targetEnergy, cfg.Chemical.TargetSystemEnergy)
	}

	// Total energy should match target initially
	if totalEnergy != targetEnergy {
		t.Errorf("Initial total energy = %v; want %v", totalEnergy, targetEnergy)
	}

	// Simulate depletion from organism consumption
	testPosition := types.Point{X: 500, Y: 500}
	world.DepleteEnergyFromSourcesAt(testPosition, 1000)

	// Get updated energy values
	newTotalEnergy, _ := world.GetSystemEnergyInfo()

	// Total energy should have decreased
	if newTotalEnergy >= totalEnergy {
		t.Errorf("Total energy after depletion = %v; want less than %v", newTotalEnergy, totalEnergy)
	}

	// The decrease should be proportional to the depletion amount (account for multiplier)
	expectedDecrease := 1000 * 2.0
	actualDecrease := totalEnergy - newTotalEnergy

	// Use a reasonable tolerance for floating-point comparisons
	// Energy distribution depends on concentration, so we can't predict exact value
	if actualDecrease < expectedDecrease*0.5 || actualDecrease > expectedDecrease*1.5 {
		t.Errorf("Energy decrease = %v; expected roughly %v", actualDecrease, expectedDecrease)
	}
}

func TestSourceCreation(t *testing.T) {
	// Create a world with specific configuration
	cfg := config.SimulationConfig{
		World: config.WorldConfig{
			Width:  1000,
			Height: 1000,
		},
		Chemical: config.ChemicalConfig{
			Count:                   0, // Start with no sources
			MinStrength:             100,
			MaxStrength:             200,
			MinDecayFactor:          0.001,
			MaxDecayFactor:          0.01,
			DepletionRate:           0.2,
			RegenerationProbability: 1.0, // High probability for testing
			TargetSystemEnergy:      10000,
		},
	}

	world := NewWorld(cfg)

	// Manually set energy deficit
	world.totalSystemEnergy = 0

	// Create a real RNG with a fixed seed for deterministic testing
	rng := rand.New(rand.NewSource(42))

	// Create a source
	initialSourceCount := len(world.GetChemicalSources())
	world.CreateChemicalSource(rng)

	// Get updated sources
	sources := world.GetChemicalSources()

	// Should have added one source
	if len(sources) != initialSourceCount+1 {
		t.Errorf("Expected %d sources after creation, got %d", initialSourceCount+1, len(sources))
	}

	// Total energy should have increased
	totalEnergy, _ := world.GetSystemEnergyInfo()
	if totalEnergy <= 0 {
		t.Errorf("Expected total energy to increase after source creation, got %v", totalEnergy)
	}

	// The new source should be active
	newSource := sources[len(sources)-1]
	if !newSource.IsActive {
		t.Error("Newly created source should be active")
	}
}

func TestUpdateChemicalSources(t *testing.T) {
	// Create a world with sources
	cfg := config.SimulationConfig{
		World: config.WorldConfig{
			Width:  1000,
			Height: 1000,
		},
		Chemical: config.ChemicalConfig{
			Count:                   3,
			MinStrength:             100,
			MaxStrength:             200,
			MinDecayFactor:          0.001,
			MaxDecayFactor:          0.01,
			DepletionRate:           0.2,
			RegenerationProbability: 1.0, // High probability for testing
			TargetSystemEnergy:      100000,
		},
	}

	world := NewWorld(cfg)

	// Get initial state
	initialSources := world.GetChemicalSources()
	initialEnergy, _ := world.GetSystemEnergyInfo()

	// Set one source to low energy to test deactivation
	if len(initialSources) > 0 {
		world.ChemicalSources[0].Energy = 0.1 // Almost depleted
	}

	// Create a real RNG with a fixed seed for deterministic testing
	rng := rand.New(rand.NewSource(42))

	// Update with large delta time to ensure depletion
	world.UpdateChemicalSources(10.0, rng)

	// Get updated state
	updatedSources := world.GetChemicalSources()
	updatedEnergy, _ := world.GetSystemEnergyInfo()

	// The first source should now be inactive
	if len(updatedSources) > 0 && updatedSources[0].IsActive {
		t.Error("Expected first source to become inactive after update")
	}

	// System energy should have changed
	if updatedEnergy == initialEnergy {
		t.Errorf("Expected system energy to change after update, but got %v before and after", updatedEnergy)
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

// setupTestWorld creates a new test world with basic configuration
func setupTestWorld() *World {
	cfg := config.SimulationConfig{
		World: config.WorldConfig{
			Width:  1000,
			Height: 1000,
		},
		Chemical: config.ChemicalConfig{
			Count:          0, // Start with no chemical sources, we'll add them manually
			MinStrength:    100,
			MaxStrength:    200,
			MinDecayFactor: 0.001,
			MaxDecayFactor: 0.01,
			DepletionRate:  0.2,
		},
	}

	return NewWorld(cfg)
}
