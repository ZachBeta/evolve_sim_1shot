# Energy and Reproduction Implementation TODO

This document outlines the specific tasks needed to implement the energy and reproduction systems as detailed in the tutorial.

## Phase 1: Energy System

### Core Energy Implementation

- [x] Update `Organism` struct with energy-related fields:
  - [x] `Energy` (current level)
  - [x] `MaxEnergy` (capacity)
  - [x] `MetabolicRate` (base consumption)
  - [x] `MovementCost` (movement energy cost)
  - [x] `SensingCost` (sensing energy cost)
  - [x] `OptimalGain` (max energy gain)
  - [x] `EnergyEfficiency` (metabolism multiplier)

- [x] Modify organism initialization to set energy values:
  - [x] Initialize with default energy level
  - [x] Add randomization for initial energy parameters
  - [x] Create appropriate constructors/factory functions

- [x] Implement energy consumption logic:
  - [x] Create `UpdateEnergy` method for `Organism`
  - [x] Add base metabolic cost calculation
  - [x] Add movement cost based on velocity
  - [x] Add sensing cost

- [x] Implement energy gain from environment:
  - [x] Calculate concentration at position
  - [x] Compare to preferred concentration
  - [x] Calculate similarity factor
  - [x] Apply energy gain based on similarity

- [x] Implement energy-based death:
  - [x] Check for energy depletion
  - [x] Mark organism for removal when energy <= 0
  - [x] Update organism removal logic in simulation

### Energy Configuration

- [x] Create `EnergyConfig` struct:
  - [x] Add all required parameters
  - [x] Set reasonable default values
  - [x] Add documentation for each parameter

- [x] Integrate energy config with main config system:
  - [x] Add to existing config structure
  - [x] Update config loading/saving
  - [x] Add config validation

### Energy Visualization

- [ ] Update organism rendering based on energy:
  - [ ] Modify color intensity based on energy level
  - [ ] Add visual indicator for critically low energy

- [ ] Add optional energy bar visualization:
  - [ ] Design energy bar UI element
  - [ ] Position relative to organism
  - [ ] Make togglable via config/key

- [ ] Add special effects for energy events:
  - [ ] Visual effect for energy gain
  - [ ] Visual effect for energy depletion
  - [ ] Death animation/effect

### Integration with Simulation

- [x] Update simulation loop to handle energy system:
  - [x] Call organism energy update during simulation step
  - [ ] Track energy statistics (avg, min, max)
  - [x] Remove dead organisms based on energy

- [ ] Add energy debugging tools:
  - [ ] Logging for energy levels
  - [ ] Statistics display
  - [ ] Toggle for energy visualization

## Phase 2: Basic Reproduction System

### Reproduction Mechanics

- [x] Implement reproduction threshold logic:
  - [x] Create `CheckReproduction` method
  - [x] Implement energy threshold check
  - [x] Add reproduction cooldown/timing

- [x] Implement offspring creation:
  - [x] Create `Reproduce` method
  - [x] Set offspring initial position near parent
  - [x] Handle energy division between parent and offspring

- [x] Create offspring placement logic:
  - [x] Determine safe position for offspring
  - [x] Handle edge cases (boundaries, obstacles)
  - [x] Add randomization to placement

- [x] Update simulation to handle new organisms:
  - [x] Add offspring to world during simulation step
  - [ ] Track reproduction events
  - [x] Add population controls and limits

### Reproduction Configuration

- [x] Create `ReproductionConfig` struct:
  - [x] Add threshold, energy transfer, and distance parameters
  - [x] Set reasonable default values
  - [x] Add documentation

- [x] Integrate with main config system:
  - [x] Add to existing config structure
  - [x] Update config loading/saving
  - [ ] Add validation for reproduction parameters

## Phase 3: Genetic Inheritance and Mutation

### Trait Inheritance System

- [x] Implement basic inheritance logic:
  - [x] Copy parent traits to offspring
  - [x] Handle trait initialization
  - [x] Create trait inheritance helper functions

- [x] Create mutation system:
  - [x] Implement `mutateValue` function
  - [x] Set up mutation probability and magnitude
  - [x] Add bounds checking for mutated values

- [x] Implement sensor mutation:
  - [x] Create `mutateSensors` function
  - [x] Add randomization to sensor angles and distances
  - [x] Ensure sensors remain within valid ranges

### Visualization of Mutations

- [ ] Update organism appearance based on traits:
  - [ ] Vary organism color based on preferences
  - [ ] Adjust size/shape based on energy capacity
  - [ ] Visualize sensor differences

- [ ] Add reproduction visual effects:
  - [ ] Momentary glow/pulse when reproducing
  - [ ] Connection line between parent and offspring
  - [ ] Color inheritance with slight variation

## Phase 4: Testing and Balancing

### Testing Tools

- [ ] Create debug visualization modes:
  - [ ] Energy level heat map
  - [ ] Reproduction events tracking
  - [ ] Preferred concentration visualization

- [ ] Add statistics collection:
  - [ ] Track population size over time
  - [ ] Monitor average energy levels
  - [ ] Record reproduction rates
  - [ ] Calculate mutation statistics

- [ ] Implement parameter adjustment UI:
  - [ ] Real-time energy parameter adjustment
  - [ ] Reproduction parameter tuning
  - [ ] Mutation rate/magnitude controls

### Balance and Optimization

- [ ] Test and tune energy parameters:
  - [ ] Balance energy gain vs consumption rates
  - [ ] Adjust movement and sensing costs
  - [ ] Validate death mechanics

- [ ] Optimize reproduction system:
  - [ ] Fine-tune reproduction threshold
  - [ ] Adjust offspring distance parameters
  - [ ] Balance energy transfer ratio

- [ ] Test mutation system:
  - [ ] Validate mutation rates produce sufficient variation
  - [ ] Test bounds on mutated values
  - [ ] Ensure trait distributions remain reasonable

## Phase 5: Advanced Features (Optional)

### Population Dynamics

- [x] Implement population controls:
  - [x] Add maximum population limit
  - [ ] Create population density effects
  - [ ] Add carrying capacity mechanics

- [x] Track lineages:
  - [x] Add generation counter
  - [x] Implement organism IDs
  - [ ] Add optional family tree visualization

### Adaptive Behaviors

- [ ] Add behavioral adaptations:
  - [ ] Allow mutation of decision-making parameters
  - [ ] Create emergent group behaviors
  - [ ] Add memory/learning mechanisms

- [ ] Implement environmental challenges:
  - [ ] Add resource scarcity events
  - [ ] Create seasonal/cyclical changes
  - [ ] Add environmental hazards

## Integration Points

### Simulation Engine Integration

- [x] Update main simulation loop:
  - [x] Add energy processing
  - [x] Handle reproduction
  - [x] Process mutations
  - [ ] Track statistics

### Renderer Integration 

- [ ] Modify renderer to handle new visuals:
  - [ ] Energy level visualization
  - [ ] Reproduction events
  - [ ] Genetic variation visualization

### UI Integration

- [ ] Add energy and reproduction UI elements:
  - [ ] Energy statistics display
  - [ ] Reproduction rate counter
  - [ ] Population graph
  - [ ] Parameter adjustment controls

### Configuration Integration

- [x] Update configuration system:
  - [x] Add energy parameters
  - [x] Add reproduction parameters
  - [x] Add mutation parameters
  - [ ] Create preset configurations 