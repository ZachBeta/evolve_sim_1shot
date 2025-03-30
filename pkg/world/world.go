package world

import (
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/zachbeta/evolve_sim/pkg/config"
	"github.com/zachbeta/evolve_sim/pkg/types"
)

// World extends the basic types.World with additional functionality
type World struct {
	types.World
	config            config.WorldConfig
	mutex             sync.RWMutex
	concentrationGrid *ConcentrationGrid
}

// NewWorld creates a new world with the specified configuration
func NewWorld(cfg config.SimulationConfig) *World {
	baseWorld := types.NewWorld(cfg.World.Width, cfg.World.Height)
	world := &World{
		World:  baseWorld,
		config: cfg.World,
	}

	// Populate the world with organisms and chemical sources
	world.PopulateWorld(cfg)

	// Initialize the concentration grid for faster lookups
	world.InitializeConcentrationGrid(5.0)

	return world
}

// GetConfig returns the world configuration
func (w *World) GetConfig() config.WorldConfig {
	return w.config
}

// AddOrganism adds an organism to the world thread-safely
func (w *World) AddOrganism(org types.Organism) bool {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	return w.World.AddOrganism(org)
}

// AddChemicalSource adds a chemical source to the world thread-safely
// and invalidates the concentration grid
func (w *World) AddChemicalSource(source types.ChemicalSource) bool {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	success := w.World.AddChemicalSource(source)
	if success {
		// Invalidate the concentration grid
		w.concentrationGrid = nil
	}
	return success
}

// GetOrganisms returns a copy of the organisms slice to avoid concurrent modification
func (w *World) GetOrganisms() []types.Organism {
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	// Create a copy to avoid concurrent modification
	orgCopy := make([]types.Organism, len(w.Organisms))
	copy(orgCopy, w.Organisms)
	return orgCopy
}

// GetChemicalSources returns a copy of the chemical sources slice to avoid concurrent modification
func (w *World) GetChemicalSources() []types.ChemicalSource {
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	// Create a copy to avoid concurrent modification
	sourcesCopy := make([]types.ChemicalSource, len(w.ChemicalSources))
	copy(sourcesCopy, w.ChemicalSources)
	return sourcesCopy
}

// GetOrganismAt returns the organism at the specified index
func (w *World) GetOrganismAt(index int) (types.Organism, bool) {
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	if index < 0 || index >= len(w.Organisms) {
		return types.Organism{}, false
	}

	return w.Organisms[index], true
}

// UpdateOrganism updates an organism at the specified index
func (w *World) UpdateOrganism(index int, org types.Organism) bool {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if index < 0 || index >= len(w.Organisms) {
		return false
	}

	// Ensure the new position is within bounds
	if !w.Boundaries.Contains(org.Position) {
		return false
	}

	w.Organisms[index] = org
	return true
}

// GetConcentrationAt calculates the total chemical concentration at a given point
// Uses the concentration grid if available for faster lookups
func (w *World) GetConcentrationAt(point types.Point) float64 {
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	// If we have a concentration grid, use it
	if w.concentrationGrid != nil {
		return w.concentrationGrid.GetConcentrationAt(point)
	}

	// Otherwise calculate directly (slower)
	return w.World.GetConcentrationAt(point)
}

// GetConcentrationGradientAt calculates the gradient (direction of concentration change)
// at the specified point
func (w *World) GetConcentrationGradientAt(point types.Point) types.Point {
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	// If we have a concentration grid, use it for faster gradient calculation
	if w.concentrationGrid != nil {
		return w.concentrationGrid.GetGradientAt(point)
	}

	// Otherwise, calculate numerically
	const delta = 0.5 // Small distance for finite difference

	// Calculate concentrations at points slightly offset from the original
	cCenter := w.World.GetConcentrationAt(point)
	cRight := w.World.GetConcentrationAt(types.Point{X: point.X + delta, Y: point.Y})
	cUp := w.World.GetConcentrationAt(types.Point{X: point.X, Y: point.Y + delta})

	// Calculate partial derivatives
	dCdx := (cRight - cCenter) / delta
	dCdy := (cUp - cCenter) / delta

	// Return the gradient vector
	gradient := types.Point{X: dCdx, Y: dCdy}

	// Normalize if not zero
	length := math.Sqrt(gradient.X*gradient.X + gradient.Y*gradient.Y)
	if length > 1e-9 {
		gradient.X /= length
		gradient.Y /= length
	}

	return gradient
}

// InitializeConcentrationGrid initializes the concentration grid for faster lookups
func (w *World) InitializeConcentrationGrid(resolution float64) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	grid := NewConcentrationGrid(w.Width, w.Height, resolution)

	// Calculate concentration at each grid point
	for x := 0; x < grid.NumCellsX; x++ {
		for y := 0; y < grid.NumCellsY; y++ {
			worldX := float64(x) * grid.CellSize
			worldY := float64(y) * grid.CellSize
			point := types.Point{X: worldX, Y: worldY}
			concentration := w.World.GetConcentrationAt(point)
			grid.SetConcentration(x, y, concentration)
		}
	}

	w.concentrationGrid = grid
}

// GetBounds returns the world boundaries as a Rect
func (w *World) GetBounds() types.Rect {
	return types.NewRect(0, 0, w.Width, w.Height)
}

// UpdateOrganisms replaces all organisms in the world with a new set
func (w *World) UpdateOrganisms(organisms []types.Organism) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	// Validate that all organisms are within bounds
	validOrganisms := make([]types.Organism, 0, len(organisms))
	for _, org := range organisms {
		if w.Boundaries.Contains(org.Position) {
			validOrganisms = append(validOrganisms, org)
		}
	}

	// Replace the organisms
	w.Organisms = validOrganisms
}

// PopulateWorld fills the world with organisms and chemical sources based on configuration
func (w *World) PopulateWorld(cfg config.SimulationConfig) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	// Create a random number generator with the provided seed
	rng := rand.New(rand.NewSource(cfg.RandomSeed))
	if cfg.RandomSeed == 0 {
		// If no seed is provided, use current time
		rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	// Add chemical sources
	for i := 0; i < cfg.Chemical.Count; i++ {
		// Random position within world bounds
		x := rng.Float64() * w.Width
		y := rng.Float64() * w.Height

		// Random strength within configured range
		strength := cfg.Chemical.MinStrength + rng.Float64()*(cfg.Chemical.MaxStrength-cfg.Chemical.MinStrength)

		// Random decay factor within configured range
		decayFactor := cfg.Chemical.MinDecayFactor + rng.Float64()*(cfg.Chemical.MaxDecayFactor-cfg.Chemical.MinDecayFactor)

		// Create and add chemical source
		source := types.NewChemicalSource(types.Point{X: x, Y: y}, strength, decayFactor)
		w.World.AddChemicalSource(source)
	}

	// Add organisms
	for i := 0; i < cfg.Organism.Count; i++ {
		// Evenly distribute organisms in a grid-like pattern with some randomness
		rows := int(math.Sqrt(float64(cfg.Organism.Count)))
		cols := (cfg.Organism.Count + rows - 1) / rows

		row := i / cols
		col := i % cols

		// Calculate base position
		baseX := w.Width * float64(col+1) / float64(cols+1)
		baseY := w.Height * float64(row+1) / float64(rows+1)

		// Add some random offset to avoid perfect grid alignment
		offsetX := (rng.Float64() - 0.5) * w.Width * 0.2 / float64(cols)
		offsetY := (rng.Float64() - 0.5) * w.Height * 0.2 / float64(rows)

		x := baseX + offsetX
		y := baseY + offsetY

		// Make sure organism is within bounds
		x = math.Max(1.0, math.Min(w.Width-1.0, x))
		y = math.Max(1.0, math.Min(w.Height-1.0, y))

		// Random heading
		heading := rng.Float64() * 2 * math.Pi

		// Normal distribution for chemical preference
		preference := rng.NormFloat64()*cfg.Organism.PreferenceDistributionStdDev + cfg.Organism.PreferenceDistributionMean

		// Create and add organism
		organism := types.NewOrganism(
			types.Point{X: x, Y: y},
			heading,
			preference,
			cfg.Organism.Speed,
			types.DefaultSensorAngles(),
		)
		w.World.AddOrganism(organism)
	}

	// Reset the concentration grid
	w.concentrationGrid = nil
}

// Reset resets the world to its initial state
func (w *World) Reset(cfg config.SimulationConfig) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	// Clear organisms and chemical sources
	w.Organisms = []types.Organism{}
	w.ChemicalSources = []types.ChemicalSource{}

	// Reset concentration grid
	w.concentrationGrid = nil

	// Unlock mutex temporarily to allow nested locks in PopulateWorld
	w.mutex.Unlock()

	// Repopulate the world
	w.PopulateWorld(cfg)

	// Re-initialize the concentration grid
	w.InitializeConcentrationGrid(5.0)

	// Re-lock mutex to satisfy defer w.mutex.Unlock()
	w.mutex.Lock()
}
