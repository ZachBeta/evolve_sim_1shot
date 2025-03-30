# Chemical Energy Balance System Implementation Guide

This guide outlines how to implement a balanced energy system where chemical sources deplete over time and new sources spawn to maintain equilibrium.

## Design Goals

1. **Finite Energy Sources**: Chemical sources should deplete as organisms consume energy from them
2. **Dynamic Environment**: Create shifting hotspots to encourage organism migration and adaptation
3. **Energy Equilibrium**: Maintain relatively constant total energy in the system
4. **Visual Feedback**: Provide clear visual cues for depleting sources

## Implementation Steps

### Phase 1: Modify Chemical Sources

#### Step 1: Add Energy Fields to Chemical Sources

```go
// In pkg/types/chemical_source.go
type ChemicalSource struct {
    Position    Point   // Position in the world
    Strength    float64 // Maximum concentration at the source
    DecayFactor float64 // Controls how quickly concentration drops with distance
    
    // New fields for energy balance
    Energy          float64 // Current energy level of the source
    MaxEnergy       float64 // Maximum energy capacity
    DepletionRate   float64 // Base rate at which the source depletes (per second)
    IsActive        bool    // Whether the source is currently active
}

// Update constructor to initialize new fields
func NewChemicalSource(position Point, strength, decayFactor float64) ChemicalSource {
    maxEnergy := strength * 1000 // Scale max energy with strength
    
    return ChemicalSource{
        Position:       position,
        Strength:       strength,
        DecayFactor:    decayFactor,
        Energy:         maxEnergy,      // Start with full energy
        MaxEnergy:      maxEnergy,
        DepletionRate:  0.2,            // Base depletion rate (adjust as needed)
        IsActive:       true,
    }
}
```

#### Step 2: Modify Concentration Calculation

Update the `GetConcentrationAt` method to factor in the current energy level:

```go
// In pkg/types/chemical_source.go
func (s ChemicalSource) GetConcentrationAt(point Point) float64 {
    // If source is inactive, it produces no concentration
    if !s.IsActive {
        return 0
    }
    
    // Calculate distance to the source
    dx := point.X - s.Position.X
    dy := point.Y - s.Position.Y
    distanceSquared := dx*dx + dy*dy
    
    // Calculate concentration using inverse square law with decay factor
    concentration := s.Strength * math.Exp(-s.DecayFactor*distanceSquared)
    
    // Scale by energy percentage
    energyRatio := s.Energy / s.MaxEnergy
    
    return concentration * energyRatio
}
```

### Phase 2: Implement Energy Depletion

#### Step 1: Add Update Method to Chemical Sources

```go
// In pkg/types/chemical_source.go
func (s *ChemicalSource) Update(deltaTime float64, worldEnergy *float64) {
    // Skip inactive sources
    if !s.IsActive {
        return
    }
    
    // Base depletion (happens regardless of organisms)
    baseDepletion := s.DepletionRate * deltaTime
    
    // Deplete energy
    s.Energy -= baseDepletion
    
    // Track total energy removed from the system
    *worldEnergy -= baseDepletion
    
    // Check if source is depleted
    if s.Energy <= 0 {
        s.Energy = 0
        s.IsActive = false
    }
}
```

#### Step 2: Additional Depletion Based on Organism Consumption

When organisms gain energy from being in a favorable concentration, track which chemical sources contributed to that concentration:

```go
// In pkg/organism/behavior.go - inside the Update function where energy gain happens
if concentrationFit > ENERGY_GAIN_THRESHOLD {
    // Calculate energy gain
    gainFactor := (concentrationFit - ENERGY_GAIN_THRESHOLD) / (1.0 - ENERGY_GAIN_THRESHOLD)
    energyGain := gainFactor * MAX_ENERGY_GAIN * deltaTime
    
    // Add energy to organism
    org.Energy = math.Min(org.Energy+energyGain, org.EnergyCapacity)
    
    // Request energy depletion from nearby sources
    world.DepleteEnergyFromSourcesAt(org.Position, energyGain)
}
```

#### Step 3: Add Method to World for Source Depletion

```go
// In pkg/world/world.go
func (w *World) DepleteEnergyFromSourcesAt(position types.Point, amount float64) {
    w.mutex.Lock()
    defer w.mutex.Unlock()
    
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
            depletionAmount := amount * proportion * 2.0 // Multiplier for energy conversion
            
            // Deplete the source
            w.ChemicalSources[i].Energy -= depletionAmount
            
            // Check for depletion
            if w.ChemicalSources[i].Energy <= 0 {
                w.ChemicalSources[i].Energy = 0
                w.ChemicalSources[i].IsActive = false
            }
        }
    }
}
```

### Phase 3: Implement Source Regeneration

#### Step 1: Track Total System Energy

Add a field to the World struct to track total energy:

```go
// In pkg/world/world.go
type World struct {
    types.World
    config               config.WorldConfig
    mutex                sync.RWMutex
    concentrationGrid    *ConcentrationGrid
    
    // New fields
    totalSystemEnergy    float64
    targetSystemEnergy   float64
}
```

#### Step 2: Calculate Initial System Energy

Initialize the energy tracking in the World constructor:

```go
// In pkg/world/world.go - NewWorld function
func NewWorld(cfg config.SimulationConfig) *World {
    // Create world
    baseWorld := types.NewWorld(cfg.World.Width, cfg.World.Height)
    world := &World{
        World:  baseWorld,
        config: cfg.World,
    }
    
    // Populate the world with organisms and chemical sources
    world.PopulateWorld(cfg)
    
    // Calculate initial system energy
    world.targetSystemEnergy = 0
    for _, source := range world.ChemicalSources {
        world.targetSystemEnergy += source.MaxEnergy
    }
    
    world.totalSystemEnergy = world.targetSystemEnergy
    
    // Initialize the concentration grid for faster lookups
    world.InitializeConcentrationGrid(5.0)
    
    return world
}
```

#### Step 3: Add Source Creation Method

```go
// In pkg/world/world.go
func (w *World) CreateChemicalSource(rng *rand.Rand) {
    w.mutex.Lock()
    defer w.mutex.Unlock()
    
    // Calculate energy deficit in the system
    energyDeficit := w.targetSystemEnergy - w.totalSystemEnergy
    
    // Don't create if the deficit is too small
    if energyDeficit < w.targetSystemEnergy * 0.1 {
        return
    }
    
    // Determine strength of new source based on deficit
    // Make it relatively strong to create interesting new hotspots
    strength := w.config.Chemical.MinStrength + 
        rng.Float64()*(w.config.Chemical.MaxStrength-w.config.Chemical.MinStrength)
    
    // Scale based on deficit
    strength = math.Min(w.config.Chemical.MaxStrength, 
        strength * (1.0 + energyDeficit/w.targetSystemEnergy))
    
    // Determine decay factor
    decayFactor := w.config.Chemical.MinDecayFactor + 
        rng.Float64()*(w.config.Chemical.MaxDecayFactor-w.config.Chemical.MinDecayFactor)
    
    // Random position
    x := rng.Float64() * w.Width
    y := rng.Float64() * w.Height
    
    // Create and add new source
    source := types.NewChemicalSource(types.Point{X: x, Y: y}, strength, decayFactor)
    
    // Update total system energy
    w.totalSystemEnergy += source.MaxEnergy
    
    // Add the source to the world
    w.ChemicalSources = append(w.ChemicalSources, source)
    
    // Invalidate concentration grid
    w.concentrationGrid = nil
}
```

#### Step 4: Update Chemical Sources in World Update Method

```go
// In pkg/world/world.go - add a method to update chemical sources
func (w *World) UpdateChemicalSources(deltaTime float64, rng *rand.Rand) {
    w.mutex.Lock()
    defer w.mutex.Unlock()
    
    // Update each source
    activeSourceCount := 0
    for i := range w.ChemicalSources {
        w.ChemicalSources[i].Update(deltaTime, &w.totalSystemEnergy)
        if w.ChemicalSources[i].IsActive {
            activeSourceCount++
        }
    }
    
    // Check if we need to create a new source
    // Create new sources when: 
    // 1. System energy is below target
    // 2. We have at least one inactive source
    // 3. Random chance (to avoid creating too many at once)
    sourceCreationProbability := deltaTime * 0.2 // Adjust for desired frequency
    
    if w.totalSystemEnergy < w.targetSystemEnergy * 0.8 && 
       activeSourceCount < len(w.ChemicalSources) &&
       rng.Float64() < sourceCreationProbability {
        w.CreateChemicalSource(rng)
    }
}
```

### Phase 4: Integration into Simulation Loop

#### Step 1: Update Simulator's Step Method

```go
// In pkg/simulation/simulator.go - Step method
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
    
    // Process reproduction
    s.World.ProcessReproduction()
    
    // Update simulation time
    s.Time += adjustedTimeStep
}
```

#### Step 2: Add RNG to Simulator

```go
// In pkg/simulation/simulator.go - add a field
type Simulator struct {
    World           *world.World
    Config          config.SimulationConfig
    Time            float64 // Simulation time in seconds
    TimeStep        float64 // Fixed time step in seconds
    IsPaused        bool    // Flag to pause/resume simulation
    SimulationSpeed float64 // Speed multiplier
    rng             *rand.Rand // Random number generator
}

// Update constructor
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
    }
}
```

### Phase 5: Visual Enhancements

#### Step 1: Visual Indication of Source Energy

```go
// In pkg/renderer/renderer.go - update drawChemicalSources
func (r *Renderer) drawChemicalSources(screen *ebiten.Image) {
    sources := r.World.GetChemicalSources()
    
    for _, source := range sources {
        // Only draw active sources or sources that have just become inactive
        if !source.IsActive && source.Energy <= 0 {
            continue
        }
        
        // Convert to screen coordinates
        screenX, screenY := r.worldToScreen(source.Position)
        
        // Calculate size based on strength and energy level
        size := 5.0 + math.Sqrt(source.Strength) * 0.5
        
        // Scale size with energy level
        energyRatio := source.Energy / source.MaxEnergy
        size *= 0.5 + 0.5 * energyRatio
        
        // Determine color based on energy level
        alpha := uint8(255 * math.Max(0.2, energyRatio))
        
        // Pulse effect for low energy
        if energyRatio < 0.3 {
            pulseFreq := 3.0 + (0.3 - energyRatio) * 10.0 // Increase frequency as energy decreases
            pulseAmp := 0.3 * (0.3 - energyRatio) / 0.3   // Amplitude proportional to how low energy is
            pulse := 1.0 + pulseAmp * math.Sin(r.Simulator.Time * pulseFreq)
            size *= pulse
        }
        
        // Draw the source
        ebitenutil.DrawCircle(screen, screenX, screenY, size, color.RGBA{255, 200, 0, alpha})
        
        // Draw concentric rings to visualize the spread
        for i := 1; i <= 3; i++ {
            ringSize := size * float64(i) * 1.5
            ringAlpha := alpha / uint8(i*2)
            ebitenutil.DrawCircle(screen, screenX, screenY, ringSize, color.RGBA{255, 200, 0, ringAlpha})
        }
    }
}
```

#### Step 2: Add Visual Effect for Source Creation

```go
// In pkg/renderer/renderer.go - similar to reproduction events
type SourceCreationEvent struct {
    Position types.Point
    TimeLeft float64
    Strength float64
}

// Add to Renderer struct
type Renderer struct {
    // Existing fields
    sourceCreationEvents []SourceCreationEvent
}

// Add method to create an event
func (r *Renderer) AddSourceCreationEvent(position types.Point, strength float64) {
    r.sourceCreationEvents = append(r.sourceCreationEvents, SourceCreationEvent{
        Position: position,
        TimeLeft: 2.0, // 2 second effect
        Strength: strength,
    })
}

// Update and draw these events in the appropriate Renderer methods
// ...
```

## Testing the Implementation

### Key Tests

1. **Source Depletion Test**: Verify that sources deplete at the expected rate
2. **Organism Consumption Test**: Confirm that organism energy gain depletes source energy
3. **Energy Balance Test**: Check that total system energy stays within acceptable bounds
4. **Source Creation Test**: Ensure new sources are created when appropriate

### User Testing Checklist

- [ ] Chemical sources visibly deplete as organisms consume from them
- [ ] Depleted sources are less attractive to organisms
- [ ] New sources appear at appropriate intervals
- [ ] The simulation maintains interesting dynamics without too many or too few chemical sources
- [ ] Visual effects clearly communicate source energy levels

## Configurable Parameters

Consider adding these parameters to config.json:

```json
"chemical": {
  "count": 5,
  "minStrength": 100,
  "maxStrength": 500,
  "minDecayFactor": 0.001,
  "maxDecayFactor": 0.01,
  "depletionRate": 0.2,
  "regenerationProbability": 0.2,
  "targetSystemEnergy": 10000
}
```

## Expected Behavior

With this implementation, the simulation should exhibit the following behaviors:

1. Organisms will deplete chemical sources they gather energy from
2. Sources will gradually fade away as they're depleted
3. New sources will appear in random locations to maintain energy balance
4. Organisms will need to migrate to new sources as old ones deplete
5. A dynamic ecosystem will develop with organisms adapting to changing conditions 