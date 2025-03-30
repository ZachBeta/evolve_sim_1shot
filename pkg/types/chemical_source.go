package types

// ChemicalSource represents a point in the world that emits chemicals
type ChemicalSource struct {
	Position    Point   // The position of the chemical source
	Strength    float64 // The strength/concentration at the source
	DecayFactor float64 // How quickly the concentration decays with distance
}

// NewChemicalSource creates a new chemical source with the given parameters
func NewChemicalSource(position Point, strength, decayFactor float64) ChemicalSource {
	return ChemicalSource{
		Position:    position,
		Strength:    strength,
		DecayFactor: decayFactor,
	}
}

// GetConcentrationAt calculates the chemical concentration at a given point
// Uses inverse square law: concentration = strength / (1 + distance² * decayFactor)
func (cs ChemicalSource) GetConcentrationAt(point Point) float64 {
	dist := cs.Position.DistanceTo(point)

	// Avoid division by zero if point is at source
	if dist < 1e-9 {
		return cs.Strength
	}

	return cs.Strength / (1.0 + dist*dist*cs.DecayFactor)
}
