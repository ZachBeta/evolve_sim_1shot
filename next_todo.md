# Evolutionary Simulator Enhancement Tutorial & Todo

This document serves as both a tutorial and a checklist for implementing the next version of our evolutionary simulator. Each section includes theory, implementation steps, and testing requirements.

## Visual Enhancements

### Phase 1: Chemical Gradient Visualization

#### 1.1: Enhanced Heat Map Implementation ✅

**Theory**: Heat maps use color to represent scalar values (concentration levels). Good color maps should be perceptually uniform, colorblind-friendly, and intuitive.

- [x] **Upgrade color interpolation**
  - [x] Study existing `drawChemicalConcentration` method in `pkg/renderer/renderer.go`
  - [x] Replace direct RGB interpolation with HSL or Lab color space interpolation
  - [x] Add smoothing function to reduce banding artifacts
  - [x] Test with extreme concentration values

- [x] **Implement configurable color schemes**
  - [x] Create `ColorScheme` struct with predefined gradients
  - [x] Add at least 3 schemes: viridis (default), magma, and plasma
  - [x] Create `SetColorScheme` method in `Renderer`
  - [x] Add keyboard shortcut for cycling schemes (e.g., 'M')

- [x] **Add concentration legend**
  - [x] Create `drawLegend` method in `Renderer`
  - [x] Show gradient bar with min/max labels
  - [x] Position in bottom-right corner of screen
  - [x] Add toggle option (e.g., 'L' key)

**Testing**: ✅ Colors are visually distinguishable and legend accurately represents concentration ranges.

#### 1.2: Contour Line Implementation ✅

**Theory**: Contour lines connect points of equal value, helping viewers understand the shape of a scalar field. The marching squares algorithm is efficient for generating contours from a grid.

- [x] **Implement marching squares algorithm**
  - [x] Create `generateContours` method in `ConcentrationGrid`
  - [x] Define lookup table for 16 possible square configurations
  - [x] Iterate through grid cells to generate line segments
  - [x] Combine segments into continuous contour lines

- [x] **Add contour configuration**
  - [x] Create data structures to store contour data
  - [x] Implement sensible default contour levels
  - [x] Add toggle option (e.g., 'O' key)

- [x] **Implement smooth line rendering**
  - [x] Draw lines with Ebiten
  - [x] Add thickness for better visibility
  - [x] Use color scheme for consistent visuals

- [x] **Add contour labels**
  - [x] Show concentration value on major contours
  - [x] Add background for better visibility
  - [x] Position at middle of contour lines

**Testing**: ✅ Contours correctly represent concentration levels and remain stable during simulation.

### Current Priorities

#### 1. Organism Visualization (High Priority)

- [ ] **Implement triangle-based direction indicator**
  - [ ] Replace current circle+line with triangle shape
  - [ ] Size proportionally to organism
  - [ ] Add smooth rotation during turns
  - [ ] Ensure visibility over background

**Rationale**: This will make organisms easier to track and understand their movement, which is a core part of the simulation.

#### 2. Interactive Features (High Priority)

- [ ] **Implement click-to-place chemical sources**
  - [ ] Add mouse position tracking in world coordinates
  - [ ] Create method to add sources at clicked position
  - [ ] Add visual feedback during placement

**Rationale**: This adds significant user interaction and makes the simulation more engaging.

#### 3. Organism Trails (Medium Priority)

- [ ] **Implement basic trail effect**
  - [ ] Store recent positions for organisms
  - [ ] Draw fading trails behind moving organisms
  - [ ] Add toggle key for trails

**Rationale**: Trails help visualize movement patterns over time, revealing behavior patterns.

#### 4. Organism Selection (Medium Priority)

- [ ] **Implement organism selection**
  - [ ] Add click detection for organisms
  - [ ] Highlight selected organism
  - [ ] Display organism properties (preference, etc.)

**Rationale**: This allows users to inspect individual organisms and understand their behavior.

### Lower Priority Tasks

#### Sensor Visualization

- [ ] **Add detailed sensor visualization**
  - [ ] Draw sensor positions as small dots
  - [ ] Color-code based on sensor readings
  - [ ] Show numerical readings

#### Chemical Flow Visualization

- [ ] **Implement gradient vector field**
  - [ ] Calculate gradient vectors
  - [ ] Show direction arrows
  - [ ] Add directional cues

#### Control Fixes and UI Refinement

- [ ] **Fix control mappings**
  - [ ] Ensure all keys work correctly
  - [ ] Add confirmation messages for toggles
  - [ ] Handle key conflicts

#### Advanced Visualization UI

- [ ] **Create visualization control panel**
  - [ ] Add settings UI
  - [ ] Create presets
  - [ ] Save/load configurations

## Simulation Mechanics (Future Enhancements)

### 1. Energy System

- [ ] **Design core energy mechanics**
  - [ ] Add energy field to Organism struct
  - [ ] Implement energy gain from optimal concentrations
  - [ ] Add energy cost for movement and actions

### 2. Reproduction System

- [ ] **Implement basic reproduction**
  - [ ] Add reproduction threshold to energy system
  - [ ] Create offspring generation method
  - [ ] Handle position and initial energy allocation

## Implementation Guidelines

### Getting Started

1. Select one high-priority task from the current priorities section
2. Implement it completely before moving to the next feature
3. Test thoroughly before marking as complete
4. Update this document to track progress

### Performance Considerations

- Always profile before and after optimization
- Consider the impact of new features on large simulations
- Use efficient data structures and algorithms
- Implement level-of-detail rendering for distant objects

## Progress Tracking

- [x] Enhanced Heat Map Implementation
- [x] Contour Line Implementation
- [ ] Triangle-based Organism Representation
- [ ] Click-to-place Chemical Sources
- [ ] Organism Trails
- [ ] Organism Selection
- [ ] Sensor Visualization
- [ ] Energy System Design 