# Energy and Reproduction Systems Tutorial

This document provides a detailed plan for implementing energy and reproduction systems in our evolutionary simulator, which will enable emergent behaviors through natural selection.

## Energy System

### Conceptual Overview

The energy system adds a metabolic dimension to organisms, creating resource constraints that drive selection pressure:

- Each organism has a finite energy reserve
- Energy is consumed by activities (movement, sensing, existing)
- Energy is gained by finding optimal environmental conditions
- When energy depletes completely, the organism dies
- Energy surplus enables reproduction

### Implementation Details

#### 1. Organism Structure Modifications

```go
type Organism struct {
    // Existing fields
    Position      vector.Vector2D
    Velocity      vector.Vector2D
    Heading       float64
    PreferredConc float64
    // ...
    
    // New energy-related fields
    Energy            float64   // Current energy level
    MaxEnergy         float64   // Maximum energy capacity
    MetabolicRate     float64   // Base energy consumption per tick
    MovementCost      float64   // Energy cost per unit of movement
    SensingCost       float64   // Energy cost for sensing operations
    OptimalGain       float64   // Maximum energy gain in optimal conditions
    EnergyEfficiency  float64   // Multiplier affecting energy consumption (mutable trait)
}
```

#### 2. Energy Consumption

Energy is consumed in several ways:

```go
func (o *Organism) UpdateEnergy(world *World) {
    // Base metabolic cost (just existing)
    o.Energy -= o.MetabolicRate * o.EnergyEfficiency
    
    // Movement cost based on velocity magnitude
    moveCost := o.Velocity.Magnitude() * o.MovementCost * o.EnergyEfficiency
    o.Energy -= moveCost
    
    // Sensing cost (fixed per operation)
    o.Energy -= o.SensingCost * o.EnergyEfficiency
    
    // Energy gain from environment
    concAtPosition := world.GetConcentrationAt(o.Position)
    similarityFactor := 1.0 - math.Abs(concAtPosition-o.PreferredConc)/o.PreferredConc
    
    // Max gain when concentration matches preference exactly
    if similarityFactor > 0 {
        energyGain := o.OptimalGain * similarityFactor
        o.Energy += energyGain
        
        // Cap at maximum energy
        if o.Energy > o.MaxEnergy {
            o.Energy = o.MaxEnergy
        }
    }
    
    // Check for death condition
    if o.Energy <= 0 {
        o.MarkForRemoval = true
    }
}
```

#### 3. Energy Visualization

- Color intensity proportional to energy level (bright = full energy, dim = low energy)
- Optional energy bar above organisms
- Special effects for energy gain/consumption events

#### 4. Configuration Parameters

```go
type EnergyConfig struct {
    InitialEnergy         float64
    MaximumEnergy         float64
    BaseMetabolicRate     float64
    MovementCostFactor    float64
    SensingCostBase       float64
    OptimalEnergyGainRate float64
    StarvationThreshold   float64
    EnergyEfficiencyRange [2]float64  // Min/max for random initialization
}
```

## Reproduction System

### Conceptual Overview

Reproduction creates genetic diversity and enables natural selection:

- Organisms reproduce when energy exceeds threshold
- Offspring inherit parent traits with mutations
- Energy is divided between parent and offspring
- Mutation allows organisms to adapt to environment over generations
- Selection pressure favors organisms with more efficient traits

### Implementation Details

#### 1. Reproduction Mechanics

```go
// Threshold-based reproduction check during update
func (o *Organism) CheckReproduction(world *World, config ReproductionConfig) *Organism {
    if o.Energy >= config.ReproductionThreshold {
        // Create offspring
        offspring := o.Reproduce(config)
        
        // Energy division
        energyForOffspring := o.Energy * config.EnergyTransferRatio
        o.Energy -= energyForOffspring
        offspring.Energy = energyForOffspring
        
        return offspring
    }
    return nil
}
```

#### 2. Genetic Inheritance with Mutations

```go
func (o *Organism) Reproduce(config ReproductionConfig) *Organism {
    offspring := &Organism{
        // Base position near parent
        Position: o.Position.Add(vector.RandomUnitVector().Scale(config.OffspringDistance)),
        Heading:  rand.Float64() * 2 * math.Pi, // Random initial heading
        
        // Inherit traits with mutation
        PreferredConc:    o.mutateValue(o.PreferredConc, config.MutationRate, config.MutationMagnitude),
        MaxEnergy:        o.mutateValue(o.MaxEnergy, config.MutationRate, config.MutationMagnitude),
        MetabolicRate:    o.mutateValue(o.MetabolicRate, config.MutationRate, config.MutationMagnitude),
        MovementCost:     o.mutateValue(o.MovementCost, config.MutationRate, config.MutationMagnitude),
        SensingCost:      o.mutateValue(o.SensingCost, config.MutationRate, config.MutationMagnitude),
        OptimalGain:      o.mutateValue(o.OptimalGain, config.MutationRate, config.MutationMagnitude),
        EnergyEfficiency: o.mutateValue(o.EnergyEfficiency, config.MutationRate, config.MutationMagnitude),
        
        // Add random variance to sensor placement/distance
        Sensors: o.mutateSensors(config.MutationRate, config.MutationMagnitude),
    }
    
    return offspring
}

func (o *Organism) mutateValue(value float64, rate float64, magnitude float64) float64 {
    if rand.Float64() < rate {
        // Apply mutation with random magnitude
        mutationFactor := 1.0 + ((rand.Float64()*2)-1)*magnitude
        return value * mutationFactor
    }
    return value
}
```

#### 3. Visualization of Reproduction

- Momentary visual effect (glow/pulse) when reproduction occurs
- Color inheritance with slight variation between parent and offspring
- Optional family tree tracking visualization (advanced feature)

#### 4. Configuration Parameters

```go
type ReproductionConfig struct {
    ReproductionThreshold float64   // Energy required to reproduce
    EnergyTransferRatio   float64   // Portion of energy given to offspring
    OffspringDistance     float64   // How far offspring spawns from parent
    MutationRate          float64   // Probability of trait mutation 
    MutationMagnitude     float64   // Maximum percent change when mutation occurs
    MaxPopulation         int       // Optional cap on total population
}
```

## Expected Emergent Behaviors

Implementing these systems should lead to several emergent phenomena:

1. **Adaptive Radiation** - Organisms will diversify to occupy different niches based on concentration preferences

2. **Efficiency Evolution** - Over time, organisms will evolve more efficient:
   - Movement patterns (less random wandering)
   - Energy consumption rates
   - Reproduction strategies (timing, energy allocation)

3. **Boom and Bust Cycles** - Population dynamics responding to resource availability

4. **Spatial Distribution Patterns** - Organisms will cluster around optimal resources or spread out to minimize competition

5. **Evolutionary Arms Races** - Competitive adaptations between different organism lineages

## Implementation Approach

1. **Start with Energy System**
   - Implement basic energy mechanics without reproduction
   - Test and balance energy gain/loss rates
   - Verify death mechanics work properly

2. **Add Basic Reproduction**
   - Implement threshold-based reproduction
   - Add simple trait inheritance
   - Test population stability

3. **Add Mutations and Selection**
   - Implement trait mutation system
   - Test over many generations
   - Observe adaptations to environment

4. **Balance and Tune**
   - Adjust parameters for interesting but stable simulations
   - Ensure system leads to emergent behaviors without crashing
   - Document interesting parameter combinations

## Technical Integration

The energy and reproduction systems should be integrated with:

1. **Simulation Engine** - Update loop must process energy and reproduction
2. **Renderer** - New visualization elements for energy and reproduction events
3. **Configuration** - New parameters must be exposed and documented
4. **UI** - Information about energy levels and reproduction events should be visible

## Future Extensions

Once the base systems are working, consider these extensions:

1. **Age and Lifespan** - Organisms die after certain time/energy cycles
2. **Sexual Reproduction** - Requiring two organisms to reproduce
3. **Specialization** - Different reproduction/energy strategies (r/K selection)
4. **Predation** - Organisms gaining energy by consuming others
5. **Resource Cycles** - Seasonal or periodic changes in resource availability 