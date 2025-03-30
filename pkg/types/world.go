package types

// World represents the simulation environment containing organisms and chemical sources
type World struct {
	Width           float64          // Width of the world
	Height          float64          // Height of the world
	Organisms       []Organism       // Collection of organisms in the world
	ChemicalSources []ChemicalSource // Collection of chemical sources in the world
	Boundaries      Rect             // Rectangular boundary of the world
}

// NewWorld creates a new world with the specified dimensions
func NewWorld(width, height float64) World {
	return World{
		Width:           width,
		Height:          height,
		Organisms:       make([]Organism, 0),
		ChemicalSources: make([]ChemicalSource, 0),
		Boundaries:      NewRect(0, 0, width, height),
	}
}

// AddOrganism adds an organism to the world
// Returns true if the organism was added successfully, false if it's outside world boundaries
func (w *World) AddOrganism(org Organism) bool {
	if !w.Boundaries.Contains(org.Position) {
		return false
	}

	w.Organisms = append(w.Organisms, org)
	return true
}

// AddChemicalSource adds a chemical source to the world
// Returns true if the source was added successfully, false if it's outside world boundaries
func (w *World) AddChemicalSource(source ChemicalSource) bool {
	if !w.Boundaries.Contains(source.Position) {
		return false
	}

	w.ChemicalSources = append(w.ChemicalSources, source)
	return true
}

// GetWorldBounds returns the boundaries of the world
func (w *World) GetWorldBounds() Rect {
	return w.Boundaries
}

// GetConcentrationAt calculates the total chemical concentration at a given point
func (w *World) GetConcentrationAt(point Point) float64 {
	var totalConcentration float64 = 0

	for _, source := range w.ChemicalSources {
		totalConcentration += source.GetConcentrationAt(point)
	}

	return totalConcentration
}

// OrganismCount returns the number of organisms in the world
func (w *World) OrganismCount() int {
	return len(w.Organisms)
}

// ChemicalSourceCount returns the number of chemical sources in the world
func (w *World) ChemicalSourceCount() int {
	return len(w.ChemicalSources)
}
