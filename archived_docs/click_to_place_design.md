# Click-to-Place Chemical Sources: Implementation Challenges

## Overview

The ability to place chemical sources by clicking on the simulation is a highly desirable feature that would significantly enhance user interaction. However, we're deprioritizing this feature in the current development cycle due to several technical challenges. This document outlines those challenges and proposes future approaches.

## Current Implementation Challenges

### 1. Input Handling in Ebiten

Ebiten's input handling system works well for keyboard events but has some quirks when dealing with mouse events in a simulation context:

- **Coordinate Transformation**: Converting screen coordinates to simulation world coordinates requires precise transformation logic, especially when the window and world sizes differ.
- **Event Timing**: Mouse clicks are processed in the `Update()` method, but world state modifications should ideally be synchronized with the simulation step.
- **Concurrent Access**: Adding chemical sources via mouse click while the simulation is running may cause race conditions if not properly synchronized with the simulation loop.

### 2. World State Management

The current architecture doesn't easily support runtime modification of chemical sources:

- **Concentration Grid Invalidation**: Adding new chemical sources requires rebuilding the concentration grid, which is computationally expensive.
- **Contour Line Regeneration**: New sources would require regenerating all contour lines, which could cause visual stuttering.
- **State Persistence**: There's no current mechanism to persist user-placed sources across simulation resets.

### 3. User Experience Issues

Several UX issues make implementation challenging:

- **Feedback Mechanism**: Users need clear visual feedback when placing sources, including preview effects.
- **Validation Rules**: Need logic to prevent placing sources in invalid locations (e.g., too close together).
- **Undo/Redo Support**: Good UX would include the ability to undo mistaken placements.

## Technical Implementation Approach (Future)

For future implementation, we recommend the following approach:

### 1. Mouse Input Handling

```go
// Add to Renderer struct
type Renderer struct {
    // ... existing fields
    mouseX, mouseY int
    lastClickX, lastClickY int
    placementMode bool
    placementStrength float64
}

// Add to Update() method
func (r *Renderer) Update() error {
    // ... existing code
    
    // Track mouse position
    r.mouseX, r.mouseY = ebiten.CursorPosition()
    
    // Handle mouse clicks
    if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && !r.mousePressed {
        r.mousePressed = true
        r.handleMouseClick(r.mouseX, r.mouseY)
    } else if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
        r.mousePressed = false
    }
    
    return nil
}

// New method to handle clicks
func (r *Renderer) handleMouseClick(screenX, screenY int) {
    // Convert screen coordinates to world coordinates
    worldX, worldY := r.screenToWorld(float64(screenX), float64(screenY))
    worldPoint := types.Point{X: worldX, Y: worldY}
    
    // Create new chemical source
    newSource := types.ChemicalSource{
        Position: worldPoint,
        Strength: r.placementStrength,
        DecayFactor: r.config.Chemical.DefaultDecayFactor,
    }
    
    // Add to world (needs to be thread-safe)
    r.World.AddChemicalSource(newSource)
}
```

### 2. World Modifications

The `World.AddChemicalSource` method would need to be enhanced:

```go
func (w *World) AddChemicalSource(source types.ChemicalSource) bool {
    w.mutex.Lock()
    defer w.mutex.Unlock()
    
    // Add to sources list
    w.ChemicalSources = append(w.ChemicalSources, source)
    
    // Mark concentration grid for rebuild
    w.concentrationGridNeedsRebuild = true
    
    return true
}

// Update in simulation loop
func (w *World) Update(dt float64) {
    // ... existing code
    
    if w.concentrationGridNeedsRebuild {
        w.RebuildConcentrationGrid()
        w.concentrationGridNeedsRebuild = false
    }
    
    // ... existing code
}
```

### 3. User Interface Enhancements

A proper implementation would also include:

1. **Source Placement Mode**: Toggle with a key (e.g., 'P')
2. **Strength Adjustment**: Use mouse wheel or keys to adjust new source strength
3. **Preview Visualization**: Show potential source location before clicking
4. **Source Management Panel**: List of placed sources with ability to select/delete

## Deprioritization Rationale

Given the complexity and potential bugs in the current implementation, we're deprioritizing this feature for the following reasons:

1. **Focus on Core Simulation**: We want to ensure the basic simulation mechanics are solid before adding interactive features.
2. **Architecture Requirements**: Proper implementation requires changes to several core components.
3. **Performance Concerns**: Adding frequent grid rebuilds could impact simulation performance.
4. **Testing Complexity**: Interactive features are harder to test systematically.

## Recommended Next Steps

Instead of implementing click-to-place now, we recommend:

1. **Triangle-based Organism Representation**: Improve visual clarity first
2. **Organism Trails**: Add movement pattern visualization
3. **Refactor World Update Logic**: Prepare for eventual support of dynamic sources
4. **Add Keyboard-driven Chemical Source Creation**: As an interim solution

## Conclusion

Click-to-place chemical sources remains a desirable feature, but it's being deprioritized in favor of more straightforward enhancements that will improve the simulation's visual clarity and usability. We'll revisit this feature in a future development cycle when the core architecture better supports dynamic world modifications. 