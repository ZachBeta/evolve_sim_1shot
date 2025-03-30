package simulation

import (
	"math/rand"
	"time"

	"github.com/zachbeta/evolve_sim/pkg/config"
	"github.com/zachbeta/evolve_sim/pkg/organism"
	"github.com/zachbeta/evolve_sim/pkg/types"
	"github.com/zachbeta/evolve_sim/pkg/world"
)

// ReproductionEventHandler is a function that handles reproduction events
type ReproductionEventHandler func(types.Point)

// Simulator handles the simulation loop and organism updates
type Simulator struct {
	World           *world.World
	Config          config.SimulationConfig
	Time            float64                  // Simulation time in seconds
	TimeStep        float64                  // Fixed time step in seconds
	IsPaused        bool                     // Flag to pause/resume simulation
	SimulationSpeed float64                  // Speed multiplier
	rng             *rand.Rand               // Random number generator
	OnReproduction  ReproductionEventHandler // Optional handler for reproduction events
}

// NewSimulator creates a new simulation engine with the given world and config
func NewSimulator(world *world.World, config config.SimulationConfig) *Simulator {
	// Create RNG
	var seed int64
	if config.RandomSeed != 0 {
		seed = config.RandomSeed
	} else {
		seed = time.Now().UnixNano()
	}
	rng := rand.New(rand.NewSource(seed))

	return &Simulator{
		World:           world,
		Config:          config,
		Time:            0.0,
		TimeStep:        1.0 / 60.0, // Default to 60 FPS
		IsPaused:        false,
		SimulationSpeed: config.SimulationSpeed,
		rng:             rng,
		OnReproduction:  nil,
	}
}

// SetReproductionHandler sets a function to be called when reproduction events occur
func (s *Simulator) SetReproductionHandler(handler ReproductionEventHandler) {
	s.OnReproduction = handler
}

// Step advances the simulation by one time step
func (s *Simulator) Step() {
	if s.IsPaused {
		return
	}

	// Adjust time step based on simulation speed
	adjustedTimeStep := s.TimeStep * s.SimulationSpeed

	// Get world bounds
	bounds := s.World.GetBounds()

	// Update chemical sources
	s.World.UpdateChemicalSources(adjustedTimeStep, s.rng)

	// Update each organism
	organisms := s.World.GetOrganisms()
	for i := range organisms {
		organism.Update(
			&organisms[i],
			s.World,
			bounds,
			s.Config.Organism.SensorDistance,
			s.Config.Organism.TurnSpeed,
			adjustedTimeStep,
		)
	}

	// Update world with modified organisms
	s.World.UpdateOrganisms(organisms)

	// Remove dead organisms (those with no energy)
	s.World.RemoveDeadOrganisms()

	// Process reproduction with our configuration
	reproCount, reproPositions := s.World.ProcessReproductionWithConfig(s.Config.Reproduction)

	// If reproduction events occurred and we have a handler, call it for each event
	if reproCount > 0 && s.OnReproduction != nil {
		for _, pos := range reproPositions {
			s.OnReproduction(pos)
		}
	}

	// Update simulation time
	s.Time += adjustedTimeStep
}

// Reset resets the simulation to its initial state
func (s *Simulator) Reset() {
	// Reset simulation time
	s.Time = 0.0

	// Reset the world
	s.World.Reset(s.Config)

	// Unpause the simulation
	s.IsPaused = false
}

// SetPaused sets the pause state of the simulation
func (s *Simulator) SetPaused(paused bool) {
	s.IsPaused = paused
}

// SetSimulationSpeed sets the simulation speed
func (s *Simulator) SetSimulationSpeed(speed float64) {
	// Enforce minimum speed
	if speed < 0.1 {
		speed = 0.1
	}

	// Enforce maximum speed
	if speed > 20.0 {
		speed = 20.0
	}

	s.SimulationSpeed = speed
}
