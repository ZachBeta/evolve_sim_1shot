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
	config         config.WorldConfig
	chemicalConfig config.ChemicalConfig // Store chemical config separately

	// Replace single mutex with more granular locks
	sourceMutex   sync.RWMutex // For chemical sources
	organismMutex sync.RWMutex // For organisms
	gridMutex     sync.RWMutex // For concentration grid
	energyMutex   sync.RWMutex // For energy tracking

	concentrationGrid *ConcentrationGrid

	// New fields for energy balance
	totalSystemEnergy  float64
	targetSystemEnergy float64
}

// NewWorld creates a new world with the specified configuration
func NewWorld(cfg config.SimulationConfig) *World {
	baseWorld := types.NewWorld(cfg.World.Width, cfg.World.Height)
	world := &World{
		World:          baseWorld,
		config:         cfg.World,
		chemicalConfig: cfg.Chemical, // Store chemical config
	}

	// Populate the world with organisms and chemical sources
	world.PopulateWorld(cfg)

	// Calculate initial system energy
	// Use configured targetSystemEnergy if available, otherwise calculate based on sources
	if cfg.Chemical.TargetSystemEnergy > 0 {
		world.targetSystemEnergy = cfg.Chemical.TargetSystemEnergy
	} else {
		for _, source := range world.ChemicalSources {
			world.targetSystemEnergy += source.MaxEnergy
		}
	}

	// Initialize total energy to match target
	world.totalSystemEnergy = world.targetSystemEnergy

	// Initialize the concentration grid for faster lookups with larger cell size for better performance
	world.InitializeConcentrationGrid(10.0)

	return world
}

// GetConfig returns the world configuration
func (w *World) GetConfig() config.WorldConfig {
	return w.config
}

// AddOrganism adds an organism to the world thread-safely
func (w *World) AddOrganism(org types.Organism) bool {
	w.organismMutex.Lock()
	defer w.organismMutex.Unlock()

	return w.World.AddOrganism(org)
}

// AddChemicalSource adds a chemical source to the world thread-safely
// and invalidates the concentration grid
func (w *World) AddChemicalSource(source types.ChemicalSource) bool {
	w.sourceMutex.Lock()
	defer w.sourceMutex.Unlock()

	success := w.World.AddChemicalSource(source)
	if success {
		// Invalidate the concentration grid
		w.concentrationGrid = nil
	}
	return success
}

// GetOrganisms returns a copy of the organisms slice to avoid concurrent modification
func (w *World) GetOrganisms() []types.Organism {
	w.organismMutex.RLock()
	defer w.organismMutex.RUnlock()

	// Create a copy to avoid concurrent modification
	orgCopy := make([]types.Organism, len(w.Organisms))
	copy(orgCopy, w.Organisms)
	return orgCopy
}

// GetChemicalSources returns a copy of the chemical sources slice to avoid concurrent modification
func (w *World) GetChemicalSources() []types.ChemicalSource {
	w.sourceMutex.RLock()
	defer w.sourceMutex.RUnlock()

	// Create a copy to avoid concurrent modification
	sourcesCopy := make([]types.ChemicalSource, len(w.ChemicalSources))
	copy(sourcesCopy, w.ChemicalSources)
	return sourcesCopy
}

// GetOrganismAt returns the organism at the specified index
func (w *World) GetOrganismAt(index int) (types.Organism, bool) {
	w.organismMutex.RLock()
	defer w.organismMutex.RUnlock()

	if index < 0 || index >= len(w.Organisms) {
		return types.Organism{}, false
	}

	return w.Organisms[index], true
}

// UpdateOrganism updates an organism at the specified index
func (w *World) UpdateOrganism(index int, org types.Organism) bool {
	w.organismMutex.Lock()
	defer w.organismMutex.Unlock()

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
	// First check if we have a concentration grid
	w.gridMutex.RLock()
	grid := w.concentrationGrid
	w.gridMutex.RUnlock()

	if grid != nil {
		return grid.GetConcentrationAt(point)
	}

	// Otherwise calculate directly (slower)
	w.sourceMutex.RLock()
	defer w.sourceMutex.RUnlock()

	return w.World.GetConcentrationAt(point)
}

// GetConcentrationGradientAt calculates the gradient (direction of concentration change)
// at the specified point
func (w *World) GetConcentrationGradientAt(point types.Point) types.Point {
	w.gridMutex.RLock()
	defer w.gridMutex.RUnlock()

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
	w.gridMutex.Lock()
	defer w.gridMutex.Unlock()

	grid := NewConcentrationGrid(w.Width, w.Height, resolution)

	// Instead of calculating concentrations at each grid point,
	// just give the grid a reference to our chemical sources
	sources := w.GetChemicalSources()
	grid.SetSources(sources)

	w.concentrationGrid = grid
}

// GetBounds returns the world boundaries as a Rect
func (w *World) GetBounds() types.Rect {
	return types.NewRect(0, 0, w.Width, w.Height)
}

// UpdateOrganisms replaces all organisms in the world with a new set
func (w *World) UpdateOrganisms(organisms []types.Organism) {
	w.organismMutex.Lock()
	defer w.organismMutex.Unlock()

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
	w.organismMutex.Lock()
	defer w.organismMutex.Unlock()

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

		// Create organism config from simulation config
		organismConfig := types.OrganismConfig{
			InitialEnergy:         cfg.Energy.InitialEnergy,
			MaximumEnergy:         cfg.Energy.MaximumEnergy,
			BaseMetabolicRate:     cfg.Energy.BaseMetabolicRate,
			MovementCostFactor:    cfg.Energy.MovementCostFactor,
			SensingCostBase:       cfg.Energy.SensingCostBase,
			OptimalEnergyGainRate: cfg.Energy.OptimalEnergyGainRate,
			EnergyEfficiencyRange: cfg.Energy.EnergyEfficiencyRange,
		}

		// Create and add organism with energy configuration
		organism := types.NewOrganismWithConfig(
			types.Point{X: x, Y: y},
			heading,
			preference,
			cfg.Organism.Speed,
			types.DefaultSensorAngles(),
			organismConfig,
		)
		w.World.AddOrganism(organism)
	}

	// Reset the concentration grid
	w.concentrationGrid = nil
}

// Reset resets the world to its initial state
func (w *World) Reset(cfg config.SimulationConfig) {
	w.organismMutex.Lock()
	defer w.organismMutex.Unlock()

	// Clear organisms and chemical sources
	w.Organisms = []types.Organism{}
	w.ChemicalSources = []types.ChemicalSource{}

	// Reset concentration grid
	w.concentrationGrid = nil

	// Unlock mutex temporarily to allow nested locks in PopulateWorld
	w.organismMutex.Unlock()

	// Repopulate the world
	w.PopulateWorld(cfg)

	// Re-initialize the concentration grid
	w.InitializeConcentrationGrid(10.0)

	// Re-lock mutex to satisfy defer w.organismMutex.Unlock()
	w.organismMutex.Lock()
}

// GetConcentrationGrid returns the current concentration grid
func (w *World) GetConcentrationGrid() *ConcentrationGrid {
	w.gridMutex.RLock()
	defer w.gridMutex.RUnlock()

	// Ensure the grid is initialized
	if w.concentrationGrid == nil {
		// Release the read lock
		w.gridMutex.RUnlock()

		// Acquire a write lock to initialize the grid
		w.gridMutex.Lock()
		// Check again in case another thread initialized it while we were waiting
		if w.concentrationGrid == nil {
			w.InitializeConcentrationGrid(10.0)
		}
		// Downgrade to a read lock
		w.gridMutex.Unlock()
		w.gridMutex.RLock()
	}

	return w.concentrationGrid
}

// RemoveOrganism removes an organism at the specified index
func (w *World) RemoveOrganism(index int) bool {
	w.organismMutex.Lock()
	defer w.organismMutex.Unlock()

	if index < 0 || index >= len(w.Organisms) {
		return false
	}

	// Remove the organism by replacing it with the last one and truncating
	w.Organisms[index] = w.Organisms[len(w.Organisms)-1]
	w.Organisms = w.Organisms[:len(w.Organisms)-1]
	return true
}

// RemoveDeadOrganisms removes all organisms with zero or negative energy
func (w *World) RemoveDeadOrganisms() int {
	w.organismMutex.Lock()
	defer w.organismMutex.Unlock()

	aliveOrganisms := make([]types.Organism, 0, len(w.Organisms))
	removedCount := 0

	// Keep only organisms with positive energy
	for _, org := range w.Organisms {
		if org.Energy > 0 {
			aliveOrganisms = append(aliveOrganisms, org)
		} else {
			removedCount++
		}
	}

	// Update the organisms list
	w.Organisms = aliveOrganisms
	return removedCount
}

// Reproduction and population constants
const (
	DefaultMaxOrganismCount = 1000 // Default maximum number of organisms allowed in the world
)

// ProcessReproduction checks all organisms for reproduction eligibility
// and creates offspring as needed
func (w *World) ProcessReproduction() int {
	return w.ProcessReproductionWithConfig(config.ReproductionConfig{
		MaxPopulation: DefaultMaxOrganismCount,
	})
}

// ProcessReproductionWithConfig checks all organisms for reproduction eligibility
// and creates offspring based on the provided configuration
func (w *World) ProcessReproductionWithConfig(cfg config.ReproductionConfig) int {
	w.organismMutex.Lock()
	defer w.organismMutex.Unlock()

	maxPopulation := cfg.MaxPopulation
	if maxPopulation <= 0 {
		maxPopulation = DefaultMaxOrganismCount
	}

	// If we've reached the max population, don't allow reproduction
	if len(w.Organisms) >= maxPopulation {
		return 0
	}

	// Create a slice to hold new organisms
	newOrganisms := make([]types.Organism, 0)

	// Track how many new organisms were created
	reproductionCount := 0

	// Check each organism for reproduction
	for i := range w.Organisms {
		if w.Organisms[i].CanReproduce() && len(w.Organisms)+len(newOrganisms) < maxPopulation {
			// Create a new organism
			offspring := w.Organisms[i].Reproduce()

			// Ensure the offspring is within world bounds
			if w.Boundaries.Contains(offspring.Position) {
				newOrganisms = append(newOrganisms, offspring)
				reproductionCount++
			}
		}
	}

	// Add all new organisms to the world
	w.Organisms = append(w.Organisms, newOrganisms...)

	return reproductionCount
}

// GetPopulationInfo returns information about the current population
func (w *World) GetPopulationInfo() (int, float64) {
	w.organismMutex.RLock()
	defer w.organismMutex.RUnlock()

	count := len(w.Organisms)
	avgEnergy := 0.0

	// Calculate average energy
	for _, org := range w.Organisms {
		avgEnergy += org.Energy
	}

	if count > 0 {
		avgEnergy /= float64(count)
	}

	return count, avgEnergy
}

// DepleteEnergyFromSourcesAt removes energy from chemical sources based on organism consumption
func (w *World) DepleteEnergyFromSourcesAt(position types.Point, amount float64) {
	w.sourceMutex.Lock()
	defer w.sourceMutex.Unlock()

	// Calculate how much each source contributes to the concentration at this position
	totalConcentration := 0.0
	sourceConcentrations := make([]float64, len(w.ChemicalSources))

	for i, source := range w.ChemicalSources {
		if source.IsActive {
			conc := source.GetConcentrationAt(position)
			sourceConcentrations[i] = conc
			totalConcentration += conc
		}
	}

	// No concentration means no sources to deplete
	if totalConcentration <= 0 {
		return
	}

	// Distribute depletion proportionally based on concentration contribution
	for i := range w.ChemicalSources {
		if sourceConcentrations[i] > 0 {
			// Calculate proportion of total concentration from this source
			proportion := sourceConcentrations[i] / totalConcentration

			// Calculate how much energy to remove from this source
			depletionAmount := amount * proportion * 50.0 // Increased from 5.0 to 50.0 for faster depletion

			// Don't deplete more than available
			originalEnergy := w.ChemicalSources[i].Energy
			if depletionAmount > originalEnergy {
				depletionAmount = originalEnergy
			}

			// Deplete the source
			w.ChemicalSources[i].Energy -= depletionAmount

			// Track total energy removed from the system
			w.totalSystemEnergy -= depletionAmount

			// Check for depletion
			if w.ChemicalSources[i].Energy <= 0 {
				w.ChemicalSources[i].Energy = 0
				w.ChemicalSources[i].IsActive = false

				// Invalidate the concentration grid when a source becomes inactive
				w.concentrationGrid = nil
			}
		}
	}
}

// UpdateChemicalSources updates all chemical sources in the world
func (w *World) UpdateChemicalSources(deltaTime float64, rng *rand.Rand) {
	w.sourceMutex.Lock()
	defer w.sourceMutex.Unlock()

	// Process each source
	for i := range w.ChemicalSources {
		// Skip inactive sources
		if !w.ChemicalSources[i].IsActive {
			continue
		}

		// Remember energy before update
		energyBefore := w.ChemicalSources[i].Energy

		// Update the source
		w.ChemicalSources[i].Update(deltaTime, &w.totalSystemEnergy)

		// If energy changed significantly, invalidate the concentration grid
		if math.Abs(energyBefore-w.ChemicalSources[i].Energy) > energyBefore*0.05 {
			w.concentrationGrid = nil
		}
	}

	// Check if we need to regenerate depleted sources
	regenerationProbability := w.chemicalConfig.RegenerationProbability * deltaTime
	if rng.Float64() < regenerationProbability {
		// Count inactive sources
		inactiveSources := 0
		for _, source := range w.ChemicalSources {
			if !source.IsActive {
				inactiveSources++
			}
		}

		// If we have inactive sources, try to regenerate one
		if inactiveSources > 0 {
			// Find a random inactive source
			inactiveIndices := make([]int, 0, inactiveSources)
			for i, source := range w.ChemicalSources {
				if !source.IsActive {
					inactiveIndices = append(inactiveIndices, i)
				}
			}

			if len(inactiveIndices) > 0 {
				// Choose a random inactive source
				randomIndex := inactiveIndices[rng.Intn(len(inactiveIndices))]

				// Regenerate it
				w.ChemicalSources[randomIndex].Energy = w.ChemicalSources[randomIndex].MaxEnergy
				w.ChemicalSources[randomIndex].IsActive = true

				// Update system energy
				w.totalSystemEnergy += w.ChemicalSources[randomIndex].Energy

				// Invalidate the concentration grid
				w.concentrationGrid = nil
			}
		} else if len(w.ChemicalSources) < w.chemicalConfig.Count {
			// Create a new source if we're below the target count
			w.CreateChemicalSource(rng)
		}
	}
}

// CreateChemicalSource creates a new chemical source at a random position
// to maintain energy balance in the system
func (w *World) CreateChemicalSource(rng *rand.Rand) {
	// Calculate energy deficit in the system
	energyDeficit := w.targetSystemEnergy - w.totalSystemEnergy

	// Don't create if the deficit is too small
	if energyDeficit < w.targetSystemEnergy*0.01 { // Reduced threshold (was 0.1)
		return
	}

	// Determine strength of new source based on deficit and configuration
	minStrength := w.chemicalConfig.MinStrength
	maxStrength := w.chemicalConfig.MaxStrength

	// Determine strength of new source based on deficit
	// Make it relatively strong to create interesting new hotspots
	strength := minStrength + rng.Float64()*(maxStrength-minStrength)

	// Scale based on deficit (larger deficit = stronger sources)
	deficitRatio := energyDeficit / w.targetSystemEnergy
	strength = math.Min(maxStrength, strength*(1.0+deficitRatio))

	// Determine decay factor
	minDecay := w.chemicalConfig.MinDecayFactor
	maxDecay := w.chemicalConfig.MaxDecayFactor
	decayFactor := minDecay + rng.Float64()*(maxDecay-minDecay)

	// Find a random position for the new source
	// Try to keep it away from edges
	margin := w.Width * 0.1
	x := margin + rng.Float64()*(w.Width-2*margin)
	y := margin + rng.Float64()*(w.Height-2*margin)

	// Create and add the new source
	source := types.NewChemicalSource(
		types.Point{X: x, Y: y},
		strength,
		decayFactor,
	)

	// Add to the world
	added := w.AddChemicalSource(source)

	// Update system energy if source was added successfully
	if added {
		w.totalSystemEnergy += source.Energy
	}
}

// GetSystemEnergyInfo returns the current total system energy and target energy
func (w *World) GetSystemEnergyInfo() (float64, float64) {
	w.energyMutex.RLock()
	defer w.energyMutex.RUnlock()

	return w.totalSystemEnergy, w.targetSystemEnergy
}
