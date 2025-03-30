# Evolutionary Simulator Implementation Plan

This document outlines a detailed, iterative approach to building the evolutionary simulator as specified in `spec.md`. The plan is broken down into small, manageable steps, each with a corresponding prompt for a code-generation LLM.

## Overall Implementation Strategy

We'll use a test-driven development approach with these general phases:
1. Core data structures
2. Basic simulation logic
3. Chemical gradient system
4. Organism behavior
5. Visualization
6. Integration and optimization

Each step builds incrementally on previous steps, ensuring continuous integration and testability.

## Detailed Steps with LLM Prompts

### Phase 1: Project Setup and Core Data Structures

#### Step 1.1: Project Initialization

```
Create a new Go project for an evolutionary simulator. Set up the basic directory structure following Go best practices, with separate packages for simulation, rendering, and testing. Initialize a go.mod file, and create a minimal main.go that prints "Evolutionary Simulator" to verify the setup works. Include appropriate README.md and .gitignore files.
```

#### Step 1.2: Basic Data Types

```
Implement the core data types for our evolutionary simulator as specified:

1. Create a package "types" with:
   - Point struct (X, Y float64)
   - Rect struct for boundaries (X, Y, Width, Height float64)
   - ChemicalSource struct (Position, Strength, DecayFactor)
   - Organism struct (Position, Heading, ChemPreference, Speed, SensorAngles)
   - World struct (Width, Height, Organisms, ChemicalSources, Boundaries)

2. Add basic constructor functions for each type that initialize with sensible defaults.

3. Write unit tests for each constructor and any helper methods.

Ensure all fields have appropriate comments explaining their purpose and expected values/units.
```

#### Step 1.3: Configuration System

```
Create a configuration system that can be used to initialize and customize the simulation:

1. Create a "config" package with a SimulationConfig struct containing:
   - WorldConfig (size, boundaries)
   - OrganismConfig (count, speed range, preference distribution parameters)
   - ChemicalConfig (sources count, strength range, decay factor range)
   - RenderConfig (window size, frame rate)

2. Implement a function to load config from a JSON file.
3. Create a default configuration that can be used if no custom config is provided.
4. Write unit tests for the configuration loading.

Ensure the config system is flexible enough to accommodate future extensions while maintaining backward compatibility.
```

### Phase 2: World and Chemical Gradient System

#### Step 2.1: World Initialization

```
Implement the core World system that will contain our simulation:

1. Create a "world" package that uses our previously defined types.
2. Implement a NewWorld function that takes a configuration and initializes:
   - The world boundaries
   - Empty slices for organisms and chemical sources
3. Add methods to the World struct:
   - AddOrganism(organism Organism)
   - AddChemicalSource(source ChemicalSource)
   - GetWorldBounds() Rect
4. Write comprehensive tests for world initialization and the methods.

Make sure the World implementaton is thread-safe if you plan to support concurrent updates in the future.
```

#### Step 2.2: Chemical Gradient Calculation

```
Implement the chemical gradient calculation system:

1. Add to the "world" package:
   - A method to calculate chemical concentration at a given point
   - For each chemical source, use an inverse square law: concentration = strength / (1 + distance^2 * decayFactor)
   - Sum the contributions from all sources

2. Add a method to calculate the gradient (direction of increasing concentration) at a point.

3. Create helper functions for distance calculations between points.

4. Write unit tests for:
   - Concentration calculation at different points
   - Gradient calculation
   - Edge cases (no sources, point at same location as source, etc.)

Use efficient algorithms that could be optimized or cached in later phases.
```

#### Step 2.3: Chemical Concentration Grid

```
To support visualization and optimize calculations, implement a grid-based approach for chemical concentrations:

1. Add a ConcentrationGrid struct to the "world" package:
   - 2D grid of concentration values
   - Grid resolution parameter
   - Methods to initialize and update the grid

2. Update the World struct to maintain a concentration grid.

3. Add methods to:
   - Rebuild the entire grid (for initialization or when sources change)
   - Get concentration at a point by grid interpolation
   - Calculate contour lines at specified concentration levels

4. Write unit tests for grid creation, interpolation, and contour generation.

Ensure the implementation is efficient enough for real-time updates with multiple chemical sources.
```

### Phase 3: Organism Behavior

#### Step 3.1: Basic Organism Movement

```
Implement basic organism movement without sensing:

1. Create an "organism" package that imports our types.
2. Implement a Move function that updates an organism's position based on:
   - Current position
   - Heading angle
   - Speed
   - Time delta

3. Add boundary collision detection that makes organisms:
   - Stop at the boundary (position adjustment)
   - Change direction (heading adjustment) when hitting boundaries

4. Write unit tests for:
   - Basic movement in different directions
   - Boundary collision handling
   - Movement over multiple time steps

Ensure the movement calculations are accurate and numerically stable for different time steps.
```

#### Step 3.2: Organism Sensing

```
Implement the organism's ability to sense chemical concentrations:

1. Add to the "organism" package:
   - A GetSensorPositions function that calculates the positions of the three sensors (front, left, right) based on the organism's position, heading, and sensor angles
   - A ReadSensors function that takes a World and returns the chemical concentration at each sensor position

2. Write unit tests for:
   - Sensor position calculations at different headings
   - Sensor readings in various chemical gradient scenarios
   - Edge cases (sensors outside boundaries, etc.)

Make sure sensor positions are calculated correctly for all possible heading angles.
```

#### Step 3.3: Organism Decision Making

```
Implement the greedy decision-making algorithm for organisms based on sensor readings:

1. Add to the "organism" package:
   - A DecideDirection function that:
     - Takes sensor readings and the organism's chemical preference
     - Calculates how close each reading is to the preferred concentration
     - Returns a desired direction change (turn left, right, or continue straight)

2. Implement an Update function that combines movement and decision-making:
   - Read sensors
   - Decide direction
   - Update heading based on decision
   - Move forward

3. Write unit tests for:
   - Decision making in various gradient scenarios
   - Complete organism update cycle
   - Behavior when preference exactly matches a reading

Ensure the decision logic correctly implements the greedy approach of turning toward the sensor with the closest reading to the preferred concentration.
```

### Phase 4: Simulation Engine

#### Step 4.1: Basic Simulation Loop

```
Implement a basic simulation engine to update the world state:

1. Create a "simulation" package with a Simulator struct:
   - Reference to a World
   - Current simulation time
   - Time step configuration
   
2. Implement methods:
   - NewSimulator(world *World, config SimulationConfig)
   - Step() - advance simulation by one time step
   - Reset() - reset simulation to initial state

3. In the Step method:
   - Update each organism (sensing, decision making, movement)
   - Handle any new events or entities
   - Update simulation time

4. Write unit tests for:
   - Single simulation step
   - Multiple steps
   - Reset functionality

Ensure the simulation maintains consistent behavior regardless of the time step size.
```

#### Step 4.2: World Initialization with Organisms

```
Expand the world initialization to include generating organisms with a normal distribution of preferences:

1. Add to the "world" package:
   - A function to generate organisms with:
     - Evenly distributed positions
     - Random headings
     - Chemical preferences following a normal distribution

2. Update the NewWorld function to populate the initial organisms based on configuration.

3. Write unit tests for:
   - Organism generation
   - Chemical preference distribution
   - Position distribution

Ensure the generation of organisms is deterministic if given a seed, for reproducible tests.
```

#### Step 4.3: Simulation Statistics and Analytics

```
Add statistics and analytics capabilities to the simulation:

1. Create a "stats" package with:
   - Collection of metrics (organism counts, chemical concentrations, etc.)
   - Methods to calculate statistics (averages, histograms, etc.)
   - Export functionality (CSV, JSON)

2. Integrate stats collection into the simulation loop.

3. Add methods to get current statistics at any point.

4. Write unit tests for:
   - Statistics calculation
   - Data export
   - Integration with simulation

This will help with analyzing the behavior of the simulation and verifying it's working correctly.
```

### Phase 5: Visualization

#### Step 5.1: Initialize Graphics System with Ebiten

```
Set up the Ebiten graphics library for our simulation:

1. Add Ebiten to your go.mod.
2. Create a "renderer" package with:
   - A Renderer struct holding a reference to the World
   - Setup for an Ebiten game loop
   - Basic window initialization

3. Implement the required Ebiten interface methods:
   - Update()
   - Draw(screen *ebiten.Image)
   - Layout(width, height int)

4. Create a simple main function that initializes the renderer and starts the game loop.

5. Test the setup with a blank screen that runs without errors.

This step just focuses on getting the graphics system working, not on actual visualization yet.
```

#### Step 5.2: Organism Visualization

```
Implement visualization of organisms:

1. Expand the "renderer" package to:
   - Draw organisms as circles with their direction indicated
   - Color-code organisms based on their chemical preference
   - Scale visualization based on world size and window size

2. Create helper functions for:
   - Converting world coordinates to screen coordinates
   - Drawing circles and direction indicators
   - Generating colors from preference values

3. Test the renderer with a small number of static organisms.

Focus on clean, readable code that separates rendering logic from simulation logic.
```

#### Step 5.3: Chemical Gradient Visualization

```
Implement visualization of chemical gradients using contour lines:

1. Add to the "renderer" package:
   - A function to render chemical concentration contours
   - Visual indicators for chemical sources
   - Optional toggle for gradient visualization

2. Use the concentration grid and contour line generation from earlier steps.

3. Ensure efficient rendering that doesn't slow down the simulation.

4. Test with various chemical source configurations.

The contour visualization should make it easy to understand the chemical landscape that organisms are navigating.
```

### Phase 6: Integration and Refinement

#### Step 6.1: Full Integration

```
Integrate all components into a complete simulation:

1. Update the main package to:
   - Initialize configuration
   - Create the world with chemical sources and organisms
   - Set up the simulation engine
   - Initialize the renderer
   - Start the main loop

2. Handle keyboard/mouse input for:
   - Pausing/resuming simulation
   - Adjusting simulation speed
   - Toggling visualization options

3. Create a smooth loop that:
   - Updates the simulation at a fixed time step
   - Renders at the display refresh rate
   - Maintains separation between simulation and rendering

4. Test the full integration with various configurations.

Ensure the integration maintains the clean separation of concerns between components.
```

#### Step 6.2: Performance Optimization

```
Optimize the simulation for better performance:

1. Profile the application to identify bottlenecks:
   - Concentration calculations
   - Organism updates
   - Rendering

2. Implement optimizations:
   - Spatial partitioning for faster queries
   - Caching frequently calculated values
   - Parallel processing where appropriate
   - More efficient rendering techniques

3. Add configuration options for performance vs. accuracy tradeoffs.

4. Benchmark before and after optimizations to verify improvements.

Focus on optimizations that maintain correctness and don't overly complicate the code.
```

#### Step 6.3: User Interface and Controls

```
Add a basic user interface and controls:

1. Implement UI elements:
   - Statistics display (FPS, organism count, etc.)
   - Control panel for adjusting parameters
   - Visual feedback for system state

2. Add interactive features:
   - Adding chemical sources with mouse clicks
   - Selecting and inspecting organisms
   - Parameter adjustment sliders
   
3. Create keyboard shortcuts for common actions.

4. Ensure UI elements don't interfere with the simulation visualization.

The UI should be minimal and non-intrusive while still providing useful information and controls.
```

#### Step 6.4: Final Testing and Documentation

```
Finalize the project with comprehensive testing and documentation:

1. Implement integration tests for the entire system.
2. Create benchmark tests for performance-critical components.
3. Document:
   - Code (with proper Go comments)
   - System architecture
   - User guide
   - Example scenarios
   - Performance considerations

4. Clean up any technical debt:
   - Refactor complex functions
   - Standardize error handling
   - Improve naming consistency

5. Create example configurations and scenarios that demonstrate the system's capabilities.

This final step ensures the project is maintainable, understandable, and ready for extension.
```

## Development Process Recommendations

1. Complete each step fully before moving to the next
2. Write tests before implementing functionality
3. Commit code after each step
4. Review performance implications regularly
5. Maintain separations of concerns between packages
6. Consider future extensibility in all design decisions

## Potential Extensions (Post-MVP)

- Energy systems
- Reproduction and genetics
- Multiple chemical types
- Dynamic chemical sources
- Environmental obstacles
- Predator-prey relationships
- Evolution visualization (statistics over time)
- Save/load simulation states

These can be added after the core functionality is working correctly and efficiently. 