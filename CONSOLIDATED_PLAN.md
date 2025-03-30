# Evolutionary Simulator Project Plan

## Project Status
- **Core Implementation:** ✅ Complete (Phases 1-3)
- **Visualization System:** ✅ Complete (Phase 4)
- **Current Focus:** 🔄 Organism Selection & UI Refinement
- **Next Steps:** Planning for performance optimization and future enhancements

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

## Current Priority: Organism Selection

Implementation checklist:
- [ ] Add click detection for organisms
- [ ] Highlight selected organism
- [ ] Display organism properties (preference, speed, etc.)
- [ ] Track selected organism even when moving
- [ ] Add visual indicators for selected organism's sensors

Rationale: This allows users to inspect individual organisms and understand their behavior in relation to the environment.

### Implementation Plan for Organism Selection

1. **Mouse Input Detection**
   - Add mouse position tracking in renderer
   - Convert screen coordinates to world coordinates
   - Update the renderer to process mouse clicks

2. **Organism Collision Detection**
   - Implement function to find organism at clicked position
   - Consider organism size for accurate selection
   - Handle multiple organisms at the same position (select closest)

3. **Selection UI**
   - Add visual highlight for selected organism (e.g., glowing outline)
   - Store reference to selected organism in renderer
   - Add method to clear selection

4. **Properties Panel**
   - Create UI panel to display organism properties
   - Include preference, speed, sensor readings
   - Position panel in a non-intrusive location

5. **Sensor Visualization Enhancement**
   - Add detailed sensor visualization for selected organism
   - Show actual concentration values at sensor positions
   - Use color coding to indicate reading vs. preference match

6. **Testing**
   - Test selection in various organism density scenarios
   - Verify selection persists during simulation steps
   - Ensure UI elements don't interfere with simulation

## Known Issues

### Contour Line Logging
- Problem: Excessive "Generated 1941 contour lines across 6 levels" logging
- Solution: Remove or limit debug logging in contour generation code
- Priority: High (clutters console output) ✅ FIXED

## Future Enhancements (Prioritized)

### Phase 1: UI Refinement
1. **Complete Organism Selection** (Current Task)
2. **Fix Control Mappings**
   - Ensure all keys work correctly
   - Add confirmation messages for toggles

### Phase 2: Performance Optimization
1. **Profile Application**
   - Add profiling instrumentation
   - Identify bottlenecks
2. **Optimize Chemical Calculations**
   - Add spatial partitioning
   - Implement caching for gradients
3. **Optimize Rendering**
   - Batch rendering operations
   - Add level-of-detail rendering

### Phase 3: Interactive Features
1. **Click-to-Place Chemical Sources** (Deprioritized)
   - Add mouse position tracking
   - Create method to add sources at clicked position
   - Add visual feedback during placement

### Phase 4: Simulation Mechanics
1. **Energy System**
   - Add energy field to Organism
   - Implement energy gain and costs
2. **Reproduction System**
   - Add reproduction threshold
   - Create offspring generation
   - Handle mutation

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