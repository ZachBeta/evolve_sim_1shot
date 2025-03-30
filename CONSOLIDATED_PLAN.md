# Evolutionary Simulator Project Plan

## Project Status
- **Core Implementation:** ✅ Complete (Phases 1-3)
- **Visualization System:** ✅ Complete (Phase 4)
- **Current Focus:** 🔄 Simulation Mechanics
- **Next Steps:** UI Refinement and Interactive Features

## Completed Features
- ✅ Core data structures and world system
- ✅ Chemical gradient calculation and grid system
- ✅ Organism movement, sensing, and decision-making
- ✅ Simulation engine with time steps
- ✅ Basic visualization with Ebiten
- ✅ Heat map visualization with configurable color schemes
- ✅ Contour line visualization of concentration levels
- ✅ Triangle-based organism representation
- ✅ Organism trails for movement visualization
- ✅ Basic UI with keyboard controls

## Future Enhancements (Prioritized)

### Phase 1: Simulation Mechanics
1. **Energy System**
   - Add energy field to Organism
   - Implement energy gain from optimal concentrations
   - Add energy cost for movement and actions
   - Design visualization for energy levels

2. **Reproduction System**
   - Add reproduction threshold to energy system
   - Create offspring generation method
   - Handle position and initial energy allocation
   - Implement mutation of preferences and properties

3. **Implementation Plan for Energy System**
   - Extend Organism struct to include energy field
   - Modify Update function to consume energy based on movement
   - Add energy gain when organism is in preferred concentration
   - Add death mechanism when energy is depleted
   - Add visualization for organism energy levels

4. **Implementation Plan for Reproduction**
   - Define reproduction threshold
   - Create offspring generation algorithm
   - Implement genetic inheritance with mutations
   - Handle world population management
   - Add visual effects for reproduction events

### Phase 2: UI Refinement and Organism Selection

1. **Organism Selection** (Moved from Current Priority)
   - Add click detection for organisms
   - Highlight selected organism
   - Display organism properties (preference, speed, etc.)
   - Track selected organism even when moving
   - Add visual indicators for selected organism's sensors

2. **Mouse Input Detection**
   - Add mouse position tracking in renderer
   - Convert screen coordinates to world coordinates
   - Update the renderer to process mouse clicks

3. **Selection UI**
   - Add visual highlight for selected organism (e.g., glowing outline)
   - Store reference to selected organism in renderer
   - Add method to clear selection

4. **Properties Panel**
   - Create UI panel to display organism properties
   - Include preference, speed, sensor readings
   - Position panel in a non-intrusive location

5. **Control Mappings**
   - Ensure all keys work correctly
   - Add confirmation messages for toggles
   - Handle key conflicts

### Phase 3: Performance Optimization
1. **Profile Application**
   - Add profiling instrumentation
   - Identify bottlenecks
2. **Optimize Chemical Calculations**
   - Add spatial partitioning
   - Implement caching for gradients
3. **Optimize Rendering**
   - Batch rendering operations
   - Add level-of-detail rendering

### Phase 4: Interactive Features
1. **Click-to-Place Chemical Sources**
   - Add mouse position tracking
   - Create method to add sources at clicked position
   - Add visual feedback during placement

## Known Issues

### Contour Line Logging
- Problem: Excessive "Generated 1941 contour lines across 6 levels" logging
- Solution: Remove or limit debug logging in contour generation code
- Priority: High (clutters console output) ✅ FIXED

## Technical Documentation

### Project Structure
- `cmd/evolve_sim`: Main application entry point
- `pkg/types`: Core data structures
- `pkg/config`: Configuration system
- `pkg/world`: World and chemical gradient system
- `pkg/organism`: Organism behavior and movement
- `pkg/simulation`: Simulation engine
- `pkg/renderer`: Visualization system

### Key Controls
- **Space**: Pause/resume simulation
- **+/-**: Adjust simulation speed
- **C**: Toggle contour visualization
- **H**: Toggle heat map visualization
- **T**: Toggle organism trails
- **M**: Cycle color schemes
- **L**: Toggle legend display
- **Click**: Select organism (pending implementation)

## Implementation Guidelines

### Performance Considerations
- Always profile before and after optimization
- Consider impact of new features on large simulations
- Use efficient data structures and algorithms
- Implement level-of-detail rendering for distant objects

### Code Standards
- Maintain comprehensive test coverage
- Document all public APIs
- Follow Go best practices
- Keep performance-critical code paths optimized 