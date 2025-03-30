# Evolutionary Simulator Specification

## Overview
A 2D simulation of single-cell organisms responding to chemical gradients in their environment. The simulation focuses on emergent behavior through individual organism preferences and directional sensing.

## Core Components

### 1. Environment
- 2D bounded space with firm walls (no wrap-around)
- Static chemical sources that create concentration gradients
- Concentration gradients will be visualized using contour lines
- Chemical diffusion model should be efficient and support future extensions

### 2. Organisms
#### Properties
- Position (x, y coordinates)
- Direction (heading angle)
- Chemical preference (normally distributed across population)
- Constant movement speed
- Three sensory organs (front, left, right)

#### Behavior
- Directional sensing in three directions
- Greedy movement algorithm: always turn toward the direction with readings closest to preferred concentration
- Collision detection with environment boundaries
- Initial even distribution across the environment

### 3. Visualization
- Organisms rendered as circles with visible sensory organ indicators
- Color coding based on chemical preferences
- Chemical sources clearly marked
- Contour lines showing chemical concentration gradients
- Real-time updating of positions and states

## Technical Implementation

### Technology Stack
- Primary Language: Go
- Graphics Library: TBD (Ebiten, SDL2, OpenGL, or Raylib-go)
- Testing Framework: Go's built-in testing package

### Data Structures

```go
type Point struct {
    X, Y float64
}

type ChemicalSource struct {
    Position    Point
    Strength    float64
    DecayFactor float64
}

type Organism struct {
    Position          Point
    Heading          float64
    ChemPreference   float64
    Speed            float64
    SensorAngles     [3]float64  // Front, Left, Right
}

type World struct {
    Width, Height    float64
    Organisms        []Organism
    ChemicalSources  []ChemicalSource
    Boundaries       Rect
}
```

### Core Systems

#### 1. Simulation Engine
- Fixed time step updates
- Separate update and render loops
- Configurable simulation speed

#### 2. Chemical System
- Efficient gradient calculation
- Caching of concentration values where appropriate
- Contour line generation algorithm

#### 3. Organism System
- Collision detection with boundaries
- Sensor reading calculations
- Movement and rotation updates
- Population management

## Error Handling
- Boundary checks for organism movement
- Validation of configuration parameters
- Graceful degradation under performance constraints
- Logging of unusual states or behaviors

## Testing Strategy

### Unit Tests
- Chemical gradient calculations
- Organism movement and sensing
- Boundary collision detection
- Configuration validation

### Integration Tests
- Full simulation cycle
- Multiple organism interactions
- Performance under load

### Performance Tests
- Gradient calculation efficiency
- Rendering performance
- Memory usage monitoring

## Performance Requirements
- Target 60 FPS with 1000+ organisms
- Efficient chemical gradient calculations
- Minimal garbage collection impact

## Future Extensibility
Architecture should support future additions of:
- Energy systems
- Reproduction mechanics
- Environmental changes
- Additional sensory capabilities
- Different organism types
- Dynamic chemical sources

## Configuration
All major parameters should be configurable:
- World size
- Number of organisms
- Chemical source properties
- Organism properties (speed, sensor angles)
- Visualization options

## Development Phases

### Phase 1: Core Implementation
1. Basic world setup
2. Chemical gradient system
3. Organism movement and sensing
4. Simple visualization

### Phase 2: Optimization
1. Performance profiling
2. Gradient calculation optimization
3. Rendering optimization
4. Memory usage optimization

### Phase 3: Testing & Refinement
1. Comprehensive test suite
2. Edge case handling
3. Configuration system
4. Documentation

## Documentation Requirements
- Code documentation following Go standards
- System architecture documentation
- Configuration guide
- Performance tuning guide
- Example scenarios 