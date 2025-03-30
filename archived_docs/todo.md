# Phase 1: Simulation Mechanics - Implementation Guide

This guide provides a step-by-step approach to implementing the energy and reproduction systems for our evolutionary simulator.

## 1. Energy System

### Step 1: Extend the Organism Structure
- Add an `energy` field to the Organism struct
- Add an `energyCapacity` field to set maximum energy
- Add constants for initial energy values

```go
type Organism struct {
    // ... existing fields ...
    energy        float64
    energyCapacity float64
}
```

### Step 2: Implement Energy Consumption
- Modify the Update method to consume energy when moving
- The amount of energy consumed should be proportional to:
  - Distance moved
  - Speed of the organism
- Example formula: `energyCost = distanceMoved * speedFactor * ENERGY_COST_MULTIPLIER`

### Step 3: Implement Energy Gain
- Add logic to gain energy when in preferred concentration
- Calculate energy gain based on how close the current concentration is to the organism's preference
- Example code approach:

```go
// Inside organism Update method
func (o *Organism) Update(world *World, dt float64) {
    // ... existing movement code ...
    
    // Energy consumption from movement
    distanceMoved := math.Sqrt(dx*dx + dy*dy)
    energyCost := distanceMoved * o.speed * MOVEMENT_ENERGY_COST
    o.energy -= energyCost
    
    // Energy gain from being in preferred environment
    currentConcentration := world.GetConcentrationAt(o.position.X, o.position.Y)
    concentrationDiff := math.Abs(currentConcentration - o.preferredConcentration)
    concentrationFit := 1.0 - (concentrationDiff / MAX_CONCENTRATION)
    
    if concentrationFit > ENERGY_GAIN_THRESHOLD {
        energyGain := concentrationFit * MAX_ENERGY_GAIN * dt
        o.energy = math.Min(o.energy+energyGain, o.energyCapacity)
    }
}
```

### Step 4: Implement Death Mechanic
- Add logic to remove organisms when energy is depleted
- Create a method in the World struct to handle organism removal
- Example approach:

```go
// In World update method
func (w *World) Update(dt float64) {
    // ... existing code ...
    
    // Check for dead organisms
    var organismsToRemove []int
    for i, organism := range w.organisms {
        if organism.energy <= 0 {
            organismsToRemove = append(organismsToRemove, i)
        }
    }
    
    // Remove dead organisms (in reverse order to maintain indexes)
    for i := len(organismsToRemove) - 1; i >= 0; i-- {
        w.RemoveOrganism(organismsToRemove[i])
    }
}
```

### Step 5: Add Energy Visualization
- Modify the renderer to visualize organism energy levels
- Options include:
  - Color gradient based on energy percentage
  - Size variation based on energy levels
  - Energy bar above organisms

```go
// In renderer
func (r *Renderer) drawOrganism(screen *ebiten.Image, organism *Organism) {
    // ... existing drawing code ...
    
    // Add energy visualization - for example, using opacity
    energyRatio := organism.energy / organism.energyCapacity
    color := colorForEnergyLevel(energyRatio)
    
    // Use the color when drawing the organism
}

func colorForEnergyLevel(ratio float64) color.RGBA {
    // Create a color that goes from red (low energy) to green (high energy)
    r := uint8(255 * (1 - ratio))
    g := uint8(255 * ratio)
    return color.RGBA{r, g, 0, 255}
}
```

## 2. Reproduction System

### Step 1: Define Reproduction Threshold
- Add a constant for the reproduction energy threshold
- This threshold determines when an organism has enough energy to reproduce
- Add a cooldown period between reproduction events

```go
const (
    REPRODUCTION_THRESHOLD = 0.75    // Percentage of max energy
    REPRODUCTION_COOLDOWN = 5.0      // Seconds between reproduction attempts
)
```

### Step 2: Implement Reproduction Logic
- Add a method to check if an organism can reproduce
- Add a method to create offspring with inherited traits
- Implement reproduction in the organism's Update method

```go
func (o *Organism) CanReproduce() bool {
    return o.energy >= o.energyCapacity * REPRODUCTION_THRESHOLD &&
           o.timeSinceReproduction >= REPRODUCTION_COOLDOWN
}

func (o *Organism) Reproduce(world *World) *Organism {
    // Create offspring
    offspring := &Organism{
        // Copy parent properties with small mutations
        preferredConcentration: o.preferredConcentration + randomMutation(-0.1, 0.1),
        speed: o.speed + randomMutation(-0.2, 0.2),
        sensorDistance: o.sensorDistance + randomMutation(-5, 5),
        // Position slightly offset from parent
        position: vector.Vector2D{
            X: o.position.X + randomMutation(-10, 10),
            Y: o.position.Y + randomMutation(-10, 10),
        },
        // Start with a percentage of parent's energy
        energy: o.energy * 0.3,
        energyCapacity: o.energyCapacity,
    }
    
    // Reduce parent's energy
    o.energy *= 0.7
    o.timeSinceReproduction = 0
    
    return offspring
}
```

### Step 3: Add Population Management
- Modify the World struct to handle new organisms
- Add methods to track population statistics
- Set limits to prevent unbounded growth

```go
func (w *World) AddOffspring(parent *Organism) {
    if len(w.organisms) < MAX_ORGANISMS {
        offspring := parent.Reproduce(w)
        w.organisms = append(w.organisms, offspring)
    }
}

// In World update
func (w *World) Update(dt float64) {
    // ... existing code ...
    
    // Check for reproduction
    for _, organism := range w.organisms {
        if organism.CanReproduce() {
            w.AddOffspring(organism)
        }
    }
}
```

### Step 4: Add Visual Effects for Reproduction
- Add visual feedback when reproduction occurs
- Options include:
  - Momentary flash or glow
  - Animation connecting parent to offspring
  - Particle effect at reproduction location

```go
// In renderer, add a struct to track reproduction events
type ReproductionEvent struct {
    position vector.Vector2D
    timeLeft float64
}

// Add a slice to Renderer to track active events
var reproductionEvents []ReproductionEvent

// When reproduction happens, add an event
func (r *Renderer) AddReproductionEvent(position vector.Vector2D) {
    r.reproductionEvents = append(r.reproductionEvents, ReproductionEvent{
        position: position,
        timeLeft: 1.0, // 1 second display time
    })
}

// In the Draw method, render all active reproduction events
func (r *Renderer) drawReproductionEffects(screen *ebiten.Image) {
    for _, event := range r.reproductionEvents {
        // Draw a glowing circle or other effect
        // Size can be based on timeLeft for fade-out effect
    }
}
```

## 3. Testing Strategy

### Energy System Tests
- Test energy consumption with different movement patterns
- Test energy gain in different concentration environments
- Test death mechanism when energy depletes
- Verify energy visualization correctly reflects energy levels

### Reproduction System Tests
- Test reproduction threshold logic
- Test trait inheritance and mutation
- Test energy distribution between parent and offspring
- Verify population control mechanisms

## 4. Implementation Order

1. Add energy fields to Organism struct
2. Implement basic energy consumption from movement
3. Add energy gain from optimal environments
4. Implement death when energy depletes
5. Add energy visualization
6. Add reproduction threshold and methods
7. Implement offspring creation with mutations
8. Add population management and control
9. Implement visual effects for reproduction events

## 5. Checkpoint Verification

After completing each major component, verify the system behaves as expected:

### Energy System Checkpoint
- Organisms should lose energy when moving
- Organisms should gain energy in preferred environments
- Organisms should die when energy is depleted
- Energy levels should be visually apparent

### Reproduction System Checkpoint
- Organisms should reproduce when they have sufficient energy
- Offspring should have slightly different traits from parents
- Population should stabilize at a sustainable level
- Reproduction events should have visual feedback 