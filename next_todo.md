# Evolutionary Simulator Enhancement Tutorial & Todo

This document serves as both a tutorial and a checklist for implementing the next version of our evolutionary simulator. Each section includes theory, implementation steps, and testing requirements.

## Visual Enhancements

### Phase 1: Chemical Gradient Visualization

#### 1.1: Enhanced Heat Map Implementation

**Theory**: Heat maps use color to represent scalar values (concentration levels). Good color maps should be perceptually uniform, colorblind-friendly, and intuitive.

- [ ] **Upgrade color interpolation**
  - [ ] Study existing `drawChemicalConcentration` method in `pkg/renderer/renderer.go`
  - [ ] Replace direct RGB interpolation with HSL or Lab color space interpolation
  - [ ] Add smoothing function to reduce banding artifacts
  - [ ] Test with extreme concentration values

- [ ] **Implement configurable color schemes**
  - [ ] Create `ColorScheme` struct with predefined gradients
  - [ ] Add at least 3 schemes: viridis (default), magma, and plasma
  - [ ] Create `SetColorScheme` method in `Renderer`
  - [ ] Add keyboard shortcut for cycling schemes (e.g., 'C')

- [ ] **Add concentration legend**
  - [ ] Create `drawLegend` method in `Renderer`
  - [ ] Show gradient bar with min/max labels
  - [ ] Position in bottom-right corner of screen
  - [ ] Add toggle option (e.g., 'L' key)

**Testing**: Verify colors are visually distinguishable and legend accurately represents concentration ranges.

#### 1.2: Contour Line Implementation

**Theory**: Contour lines connect points of equal value, helping viewers understand the shape of a scalar field. The marching squares algorithm is efficient for generating contours from a grid.

- [ ] **Implement marching squares algorithm**
  - [ ] Create `generateContours` method in `ConcentrationGrid`
  - [ ] Define lookup table for 16 possible square configurations
  - [ ] Iterate through grid cells to generate line segments
  - [ ] Combine segments into continuous contour lines

- [ ] **Add contour configuration**
  - [ ] Create `ContourConfig` struct with:
    - [ ] Number of contour levels
    - [ ] Min/max concentration values
    - [ ] Line thickness and color
  - [ ] Add to `RenderConfig` and UI controls

- [ ] **Implement smooth line rendering**
  - [ ] Use Ebiten vectors for anti-aliased lines
  - [ ] Add option for dashed lines for minor contours
  - [ ] Implement line smoothing for better visual quality

- [ ] **Add contour labels**
  - [ ] Calculate good positions for labels (curved text)
  - [ ] Show concentration value on major contours
  - [ ] Avoid overlapping labels

**Testing**: Verify contours correctly represent concentration levels and remain stable during simulation.

### Phase 2: Organism Visualization

#### 2.1: Better Organism Representation

**Theory**: Visual representation should clearly communicate direction, properties, and state of organisms.

- [ ] **Implement direction indicator**
  - [ ] Replace existing circle+line with triangle shape
  - [ ] Size triangle proportionally to organism size
  - [ ] Ensure smooth rotation during movement
  - [ ] Add option to revert to circle style

- [ ] **Add property-based variations**
  - [ ] Scale organism size based on energy or fitness
  - [ ] Vary opacity based on age/generation
  - [ ] Add visual patterns for genetic traits

- [ ] **Add selection highlighting**
  - [ ] Create glow/outline effect for selected organism
  - [ ] Add pulse animation for better visibility
  - [ ] Ensure highlight doesn't obscure organism details

**Testing**: Verify organisms are easily distinguishable and direction is clearly visible.

#### 2.2: Sensor Visualization

**Theory**: Visualizing sensors helps users understand how organisms perceive and interact with their environment.

- [ ] **Add sensor position indicators**
  - [ ] Draw small dots at calculated sensor positions
  - [ ] Use `GetSensorPositions` from `pkg/organism` package
  - [ ] Scale dot size with organism size

- [ ] **Implement sensor connections**
  - [ ] Draw lines from organism center to each sensor
  - [ ] Use semi-transparent lines to reduce visual clutter
  - [ ] Style based on sensor importance/activity

- [ ] **Add sensor reading visualization**
  - [ ] Color-code sensor dots based on concentration reading
  - [ ] Implement size variation based on reading accuracy
  - [ ] Add small indicators for gradient direction

- [ ] **Create sensor detail view**
  - [ ] Show numerical readings when organism selected
  - [ ] Compare readings to organism preference
  - [ ] Display calculated "satisfaction" level

**Testing**: Ensure sensor visualization correctly reflects the actual sensor data and remains performant with many organisms.

### Phase 3: Movement and History Visualization

#### 3.1: Organism Trail System

**Theory**: Trails help visualize movement patterns over time, revealing behavioral patterns and environmental responses.

- [ ] **Implement basic trail effect**
  - [ ] Create `TrailManager` in `pkg/renderer`
  - [ ] Store recent positions for each organism (circular buffer)
  - [ ] Draw connecting lines or dots with decreasing opacity

- [ ] **Add trail configuration**
  - [ ] Configurable trail length (time and segments)
  - [ ] Trail style options (line, dots, both)
  - [ ] Trail visibility toggle (key 'T')

- [ ] **Optimize trail storage**
  - [ ] Implement adaptive sampling based on movement speed
  - [ ] Use pooled memory to reduce allocations
  - [ ] Add level-of-detail for distant trails

- [ ] **Enhance trail visualization**
  - [ ] Color trails based on organism state/activity
  - [ ] Add width variation for speed indication
  - [ ] Implement trail fading options

**Testing**: Verify trails correctly show movement history without impacting performance.

#### 3.2: Animation Improvements

**Theory**: Smooth animations help create a more engaging and understandable visualization.

- [ ] **Implement movement interpolation**
  - [ ] Separate logical position from rendered position
  - [ ] Use linear or cubic interpolation between steps
  - [ ] Adjust interpolation based on frame rate

- [ ] **Add rotation tweening**
  - [ ] Calculate shortest rotation path
  - [ ] Apply smooth rotation transition
  - [ ] Match rotation speed to movement fluidity

- [ ] **Add event-based effects**
  - [ ] Create particle system for significant events
  - [ ] Add effects for reproduction, death, collisions
  - [ ] Implement efficient particle management

**Testing**: Verify animations remain smooth at target frame rates even with many organisms.

### Phase 4: Advanced Visualization

#### 4.1: Chemical Flow Visualization

**Theory**: Visualizing flow fields helps users understand the gradient direction and potential organism movement.

- [ ] **Implement gradient vector field**
  - [ ] Calculate gradient vectors across concentration grid
  - [ ] Create visualization option to show vectors
  - [ ] Add scaling and density controls

- [ ] **Add particle-based flow**
  - [ ] Create particle system that follows gradients
  - [ ] Set particle color based on concentration
  - [ ] Add particle lifetime and spawn rate controls

- [ ] **Implement streamlines**
  - [ ] Generate continuous lines following gradient
  - [ ] Add length and count configuration
  - [ ] Implement line animation for flow direction

**Testing**: Ensure flow visualization correctly represents the chemical gradients.

#### 4.2: Visualization UI Integration

**Theory**: A dedicated UI for visualization settings improves user experience and makes complex features accessible.

- [ ] **Create visualization control panel**
  - [ ] Add collapsible sidebar for controls
  - [ ] Organize controls by category
  - [ ] Create keyboard shortcut (e.g., 'V') to toggle panel

- [ ] **Implement visualization presets**
  - [ ] Create set of predefined visualization styles
  - [ ] Add save/load functionality for custom presets
  - [ ] Include preset descriptions and previews

- [ ] **Add performance controls**
  - [ ] Detail level slider
  - [ ] Visual effects toggle options
  - [ ] FPS target setting

- [ ] **Implement visualization settings persistence**
  - [ ] Save settings to JSON file
  - [ ] Auto-load last used settings
  - [ ] Reset to defaults option

**Testing**: Verify UI is intuitive and all visualization options work correctly.

## Interactive Features

### 1. Chemical Source Placement

- [ ] **Implement click-to-place sources**
  - [ ] Add mouse position tracking in world coordinates
  - [ ] Create method to add sources at clicked position
  - [ ] Add visual feedback during placement

- [ ] **Add source adjustment**
  - [ ] Implement drag to adjust strength
  - [ ] Add right-click to remove sources
  - [ ] Show strength preview during adjustment

- [ ] **Create different source types**
  - [ ] Implement attractant and repellent types
  - [ ] Add visual distinction between types
  - [ ] Create UI for selecting source type

### 2. Organism Inspection

- [ ] **Implement organism selection**
  - [ ] Add hit detection for organisms
  - [ ] Create selection indicator
  - [ ] Handle selection changes (deselect, new select)

- [ ] **Create organism info panel**
  - [ ] Design compact but informative layout
  - [ ] Show key statistics and properties
  - [ ] Update in real-time

- [ ] **Add sensor visualization for selected organism**
  - [ ] Show all sensors and their readings
  - [ ] Visualize concentration at sensor positions
  - [ ] Indicate preferred direction

### 3. Enhanced UI Controls

- [ ] **Implement parameter sliders**
  - [ ] Create general-purpose slider component
  - [ ] Add sliders for simulation speed, mutation rate, etc.
  - [ ] Ensure real-time parameter updates

- [ ] **Add population statistics graph**
  - [ ] Create time-series graph component
  - [ ] Plot organism count, diversity metrics, etc.
  - [ ] Add zoom and pan controls

- [ ] **Implement save/load system**
  - [ ] Create serialization for simulation state
  - [ ] Add file dialog for save/load operations
  - [ ] Include metadata (timestamp, parameters)

## Simulation Mechanics

### 1. Energy System

- [ ] **Design core energy mechanics**
  - [ ] Add energy field to Organism struct
  - [ ] Implement energy gain from optimal concentrations
  - [ ] Add energy cost for movement and actions

- [ ] **Create energy visualization**
  - [ ] Show energy level in organism display
  - [ ] Add visual effects for energy changes
  - [ ] Implement low-energy warning indicators

- [ ] **Implement death mechanic**
  - [ ] Add check for energy depletion
  - [ ] Create death event and animation
  - [ ] Handle organism removal from simulation

### 2. Reproduction System

- [ ] **Implement basic reproduction**
  - [ ] Add reproduction threshold to energy system
  - [ ] Create offspring generation method
  - [ ] Handle position and initial energy allocation

- [ ] **Add genetic variation**
  - [ ] Design genetic representation
  - [ ] Implement mutation system
  - [ ] Create inheritance rules

- [ ] **Implement sexual reproduction option**
  - [ ] Add proximity-based mating
  - [ ] Create genetic recombination
  - [ ] Balance energy costs and benefits

## Implementation Guidelines

### Getting Started

1. **Set up your development environment**
   - Ensure all dependencies are installed
   - Create a branch for the new version
   - Read through existing codebase to understand structure

2. **Pick a manageable starting point**
   - Begin with the Enhanced Heat Map implementation
   - This builds on existing code with clear improvements
   - Provides immediate visual feedback

3. **Follow the iterative process**
   - Implement one feature at a time
   - Write tests before or alongside implementation
   - Commit frequently with descriptive messages

### Performance Considerations

- Always profile before and after optimization
- Consider the impact of new features on large simulations
- Use efficient data structures and algorithms
- Implement level-of-detail rendering for distant objects

### Code Quality Guidelines

- Maintain clear separation between visualization and simulation logic
- Document all new methods and structs
- Consider backwards compatibility
- Follow existing code style and conventions

## Resources and References

- **Marching Squares Algorithm**: [Wikipedia](https://en.wikipedia.org/wiki/Marching_squares)
- **Color Maps**: [ColorBrewer](https://colorbrewer2.org/)
- **Ebiten Documentation**: [Ebiten Website](https://ebiten.org/documents/index.html)
- **Performance Profiling**: [Go Blog](https://go.dev/blog/pprof)

## Progress Tracking

- [ ] Phase 1: Chemical Gradient Visualization
- [ ] Phase 2: Organism Visualization
- [ ] Phase 3: Movement and History Visualization
- [ ] Phase 4: Advanced Visualization
- [ ] Interactive Features Implementation
- [ ] Simulation Mechanics Enhancement 