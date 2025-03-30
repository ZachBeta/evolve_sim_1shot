package simulation

import (
	"github.com/zachbeta/evolve_sim/pkg/config"
	"github.com/zachbeta/evolve_sim/pkg/organism"
	"github.com/zachbeta/evolve_sim/pkg/world"
)

// Simulator handles the simulation loop and organism updates
type Simulator struct {
	World           *world.World
	Config          config.SimulationConfig
	Time            float64 // Simulation time in seconds
	TimeStep        float64 // Fixed time step in seconds
	IsPaused        bool    // Flag to pause/resume simulation
	SimulationSpeed float64 // Speed multiplier
}

// NewSimulator creates a new simulation engine with the given world and config
func NewSimulator(world *world.World, config config.SimulationConfig) *Simulator {
	return &Simulator{
		World:           world,
		Config:          config,
		Time:            0.0,
		TimeStep:        1.0 / 60.0, // Default to 60 FPS
		IsPaused:        false,
		SimulationSpeed: config.SimulationSpeed,
	}
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

	// Update simulation time
	s.Time += adjustedTimeStep
}

// Reset resets the simulation to its initial state
func (s *Simulator) Reset() {
	// Reset simulation time
	s.Time = 0.0

	// Reset the world
	s.World.Reset()

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
	if speed > 10.0 {
		speed = 10.0
	}

	s.SimulationSpeed = speed
}
