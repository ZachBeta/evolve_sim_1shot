package simulation

import (
	"testing"

	"github.com/zachbeta/evolve_sim/pkg/config"
	"github.com/zachbeta/evolve_sim/pkg/types"
	"github.com/zachbeta/evolve_sim/pkg/world"
)

// createTestConfig creates a configuration for testing
func createTestConfig() config.SimulationConfig {
	return config.SimulationConfig{
		World: config.WorldConfig{
			Width:  100.0,
			Height: 100.0,
		},
		Organism: config.OrganismConfig{
			Count:                        10,
			Speed:                        1.0,
			SensorDistance:               5.0,
			TurnSpeed:                    0.1,
			PreferenceDistributionMean:   30.0,
			PreferenceDistributionStdDev: 5.0,
		},
		Chemical: config.ChemicalConfig{
			Count:          1,
			MinStrength:    100.0,
			MaxStrength:    100.0,
			MinDecayFactor: 0.01,
			MaxDecayFactor: 0.01,
		},
		Render: config.RenderConfig{
			WindowWidth:  640,
			WindowHeight: 480,
			FrameRate:    60,
		},
		RandomSeed:      12345,
		SimulationSpeed: 1.0,
	}
}

func TestNewSimulator(t *testing.T) {
	// Create test config
	cfg := createTestConfig()

	// Create world
	w := world.NewWorld(cfg)

	// Create simulator
	sim := NewSimulator(w, cfg)

	// Verify initial state
	if sim.Time != 0.0 {
		t.Errorf("Expected initial time to be 0, got %f", sim.Time)
	}

	if sim.IsPaused {
		t.Errorf("Expected simulator to be initially unpaused")
	}

	if sim.TimeStep != 1.0/60.0 {
		t.Errorf("Expected time step to be 1/60, got %f", sim.TimeStep)
	}

	if sim.SimulationSpeed != cfg.SimulationSpeed {
		t.Errorf("Expected simulation speed to be %f, got %f",
			cfg.SimulationSpeed, sim.SimulationSpeed)
	}
}

func TestStep(t *testing.T) {
	// Create test config
	cfg := createTestConfig()

	// Create world
	w := world.NewWorld(cfg)

	// Add a test organism
	org := types.NewOrganism(
		types.Point{X: 50, Y: 50},
		0, // Heading east
		10.0,
		1.0,
		types.DefaultSensorAngles(),
	)
	w.AddOrganism(org)

	// Add a chemical source to create a gradient
	w.AddChemicalSource(types.ChemicalSource{
		Position:    types.Point{X: 75, Y: 50},
		Strength:    100.0,
		DecayFactor: 0.01,
	})

	// Create simulator
	sim := NewSimulator(w, cfg)

	// Take one step
	initialTime := sim.Time
	initialOrgPos := w.GetOrganisms()[0].Position

	sim.Step()

	// Verify time was updated
	if sim.Time <= initialTime {
		t.Errorf("Expected time to advance, but it didn't")
	}

	// Verify organism moved
	updatedOrgPos := w.GetOrganisms()[0].Position
	if updatedOrgPos.X == initialOrgPos.X && updatedOrgPos.Y == initialOrgPos.Y {
		t.Errorf("Expected organism to move, but it didn't")
	}
}

func TestPause(t *testing.T) {
	// Create test config
	cfg := createTestConfig()

	// Create world
	w := world.NewWorld(cfg)

	// Create simulator
	sim := NewSimulator(w, cfg)

	// Pause simulator
	sim.SetPaused(true)

	// Record initial time
	initialTime := sim.Time

	// Take a step
	sim.Step()

	// Verify time didn't change
	if sim.Time != initialTime {
		t.Errorf("Expected time to remain unchanged while paused")
	}

	// Unpause and verify it resumes
	sim.SetPaused(false)
	sim.Step()

	if sim.Time == initialTime {
		t.Errorf("Expected time to advance after unpausing")
	}
}

func TestReset(t *testing.T) {
	// Create test config
	cfg := createTestConfig()

	// Create world
	w := world.NewWorld(cfg)

	// Create simulator
	sim := NewSimulator(w, cfg)

	// Advance simulation and pause
	for i := 0; i < 10; i++ {
		sim.Step()
	}
	sim.SetPaused(true)

	// Verify time advanced and simulation is paused
	if sim.Time == 0.0 {
		t.Errorf("Test setup failed: Expected time to advance")
	}
	if !sim.IsPaused {
		t.Errorf("Test setup failed: Expected simulator to be paused")
	}

	// Reset simulator
	sim.Reset()

	// Verify reset state
	if sim.Time != 0.0 {
		t.Errorf("Expected time to be reset to 0")
	}
	if sim.IsPaused {
		t.Errorf("Expected simulator to be unpaused after reset")
	}
}

func TestSimulationSpeed(t *testing.T) {
	// Create test config
	cfg := createTestConfig()

	// Create world
	w := world.NewWorld(cfg)

	// Create simulator
	sim := NewSimulator(w, cfg)

	// Record initial step rate
	sim.Step()
	initialTimeDelta := sim.Time

	// Reset
	sim.Reset()

	// Double speed
	sim.SetSimulationSpeed(2.0)

	// Take one step
	sim.Step()

	// Verify time advanced at double rate
	doubleSpeedDelta := sim.Time

	if doubleSpeedDelta <= initialTimeDelta*1.9 || doubleSpeedDelta >= initialTimeDelta*2.1 {
		t.Errorf("Expected time to advance at double rate, initial: %f, double: %f",
			initialTimeDelta, doubleSpeedDelta)
	}

	// Test speed limits
	sim.SetSimulationSpeed(0.01) // Below minimum
	if sim.SimulationSpeed != 0.1 {
		t.Errorf("Expected minimum speed limit to be enforced")
	}

	sim.SetSimulationSpeed(15.0) // Above maximum
	if sim.SimulationSpeed != 10.0 {
		t.Errorf("Expected maximum speed limit to be enforced")
	}
}

func TestFullSimulationEnergyBalance(t *testing.T) {
	// Create a simulation with a specific configuration for testing
	cfg := config.SimulationConfig{
		World: config.WorldConfig{
			Width:  500,
			Height: 500,
		},
		Organism: config.OrganismConfig{
			Count:                        50,
			Speed:                        2.0,
			SensorDistance:               10.0,
			TurnSpeed:                    0.3,
			PreferenceDistributionMean:   50.0,
			PreferenceDistributionStdDev: 10.0,
		},
		Chemical: config.ChemicalConfig{
			Count:                   10,
			MinStrength:             100,
			MaxStrength:             200,
			MinDecayFactor:          0.001,
			MaxDecayFactor:          0.01,
			DepletionRate:           1.0, // Higher depletion rate for testing
			RegenerationProbability: 0.1,
			TargetSystemEnergy:      100000,
		},
		RandomSeed: 42, // Fixed seed for deterministic testing
	}

	testWorld := world.NewWorld(cfg)
	simulator := NewSimulator(testWorld, cfg)

	// Manually set one source to have very low energy so it definitely gets depleted
	sources := testWorld.GetChemicalSources()
	if len(sources) > 0 {
		// Set the energy level of the first source to be very low
		testWorld.ChemicalSources[0].Energy = 1.0
	}

	// Record initial state after modification
	initialSourceCount := len(testWorld.GetChemicalSources())
	_, targetEnergy := testWorld.GetSystemEnergyInfo()

	// Run simulation for many steps to observe energy balance
	steps := 500 // Fewer steps should be enough with the forced depletion
	for i := 0; i < steps; i++ {
		simulator.Step()

		// Add debugging for the first few steps
		if i < 5 || i%100 == 0 {
			t.Logf("Step %d: Source 0 energy = %.2f, active = %v",
				i, testWorld.ChemicalSources[0].Energy, testWorld.ChemicalSources[0].IsActive)
		}
	}

	// After running for a while:

	// 1. Check total energy is still within reasonable bounds
	finalEnergy, _ := testWorld.GetSystemEnergyInfo()

	// Energy should be within a reasonable percentage of the target
	// (Note: this might need adjustment based on your specific implementation)
	maxDeviation := targetEnergy * 0.3 // Allow up to 30% deviation
	if finalEnergy < targetEnergy-maxDeviation || finalEnergy > targetEnergy+maxDeviation {
		t.Errorf("System energy out of expected bounds: got %v, expected within %v of %v",
			finalEnergy, maxDeviation, targetEnergy)
	}

	// 2. Check that source energy levels have changed
	currentSources := testWorld.GetChemicalSources()

	// Track energy state of sources
	activeCount := 0
	partiallyDepletedCount := 0

	for _, src := range currentSources {
		if src.IsActive {
			activeCount++

			// Check if source has been partially depleted
			if src.Energy < src.MaxEnergy*0.95 {
				partiallyDepletedCount++
			}
		}
	}

	// We should have at least some partially depleted sources
	if partiallyDepletedCount == 0 {
		t.Error("No sources were depleted at all")
	}

	// 3. Check that the number of sources may have changed (either through depletion or creation)
	if len(currentSources) < initialSourceCount {
		t.Logf("Source count decreased: %d -> %d", initialSourceCount, len(currentSources))
	} else if len(currentSources) > initialSourceCount {
		t.Logf("Source count increased: %d -> %d", initialSourceCount, len(currentSources))
	}

	// 4. Organism population should survive
	populationCount, avgEnergy := testWorld.GetPopulationInfo()
	if populationCount == 0 {
		t.Error("All organisms died during simulation")
	}

	t.Logf("Final state: Energy=%v/%v, Sources=%v/%v active, %v partially depleted, Population=%v, AvgEnergy=%v",
		finalEnergy, targetEnergy, activeCount, len(currentSources), partiallyDepletedCount, populationCount, avgEnergy)
}
