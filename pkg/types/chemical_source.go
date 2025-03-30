package types

import (
	"math"
)

// ChemicalSource represents a point in the world that emits chemicals
type ChemicalSource struct {
	Position    Point   // The position of the chemical source
	Strength    float64 // The strength/concentration at the source
	DecayFactor float64 // How quickly the concentration decays with distance

	// New fields for energy balance
	Energy        float64 // Current energy level of the source
	MaxEnergy     float64 // Maximum energy capacity
	DepletionRate float64 // Base rate at which the source depletes (per second)
	IsActive      bool    // Whether the source is currently active
}

// NewChemicalSource creates a new chemical source with the given parameters
func NewChemicalSource(position Point, strength, decayFactor float64) ChemicalSource {
	maxEnergy := strength * 1000 // Scale max energy with strength

	return ChemicalSource{
		Position:      position,
		Strength:      strength,
		DecayFactor:   decayFactor,
		Energy:        maxEnergy, // Start with full energy
		MaxEnergy:     maxEnergy,
		DepletionRate: 0.2, // Base depletion rate (adjust as needed)
		IsActive:      true,
	}
}

// GetConcentrationAt calculates the chemical concentration at a given point
func (cs ChemicalSource) GetConcentrationAt(point Point) float64 {
	// If source is inactive, it produces no concentration
	if !cs.IsActive {
		return 0
	}

	// If strength is zero, concentration is always zero
	if cs.Strength <= 0 {
		return 0
	}

	dist := cs.Position.DistanceTo(point)

	// Avoid division by zero if point is at source
	if dist < 1e-9 {
		return cs.Strength * (cs.Energy / cs.MaxEnergy)
	}

	// Calculate concentration using inverse square law with decay factor
	concentration := cs.Strength / (1.0 + dist*dist*cs.DecayFactor)

	// Scale by energy percentage
	energyRatio := cs.Energy / cs.MaxEnergy

	return concentration * energyRatio
}

// Update updates the energy level of the chemical source
func (cs *ChemicalSource) Update(deltaTime float64, worldEnergy *float64) {
	// Skip inactive sources
	if !cs.IsActive {
		return
	}

	// Base depletion (happens regardless of organisms)
	baseDepletion := cs.DepletionRate * deltaTime

	// Don't deplete more energy than available
	baseDepletion = math.Min(baseDepletion, cs.Energy)

	// Deplete energy
	cs.Energy -= baseDepletion

	// Track total energy removed from the system
	*worldEnergy -= baseDepletion

	// Check if source is depleted
	if cs.Energy <= 0 {
		cs.Energy = 0
		cs.IsActive = false
	}
}
