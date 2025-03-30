# Chemical Energy Balance System: Next Steps Tutorial

This tutorial outlines the remaining steps to complete our chemical energy balance system implementation. We've already completed the core functionality, but there are some issues to fix and enhancements to add.

## Phase 1: Debugging Source Depletion Issues

Our tests indicate that source depletion isn't working correctly in the full simulation. Let's diagnose and fix this issue.

### Step 1: Investigate Source Depletion Logic

The issue appears to be that chemical sources don't deplete properly during simulation. Let's check the following:

1. Add debug logging to monitor source energy during simulation:
   ```go
   // In pkg/world/world.go - inside UpdateChemicalSources method
   for i := range w.ChemicalSources {
       // Log before and after energy values
       fmt.Printf("Source %d before: Energy=%.2f, Active=%v\n", i, w.ChemicalSources[i].Energy, w.ChemicalSources[i].IsActive)
       
       // Update logic here...
       
       fmt.Printf("Source %d after: Energy=%.2f, Active=%v\n", i, w.ChemicalSources[i].Energy, w.ChemicalSources[i].IsActive)
   }
   ```

2. Check call frequency by adding a counter in the UpdateChemicalSources method:
   ```go
   // In pkg/world/world.go
   var updateCounter int

   func (w *World) UpdateChemicalSources(deltaTime float64, rng *rand.Rand) {
       updateCounter++
       fmt.Printf("UpdateChemicalSources called %d times\n", updateCounter)
       
       // Existing logic...
   }
   ```

### Step 2: Fix Source Depletion Rate

The issue may be related to energy consumption by organisms. Let's adjust:

1. Ensure the depletion multiplier is appropriate:
   ```go
   // In pkg/world/world.go - inside DepleteEnergyFromSourcesAt method
   // Adjust depletionAmount calculation:
   depletionAmount := amount * proportion * 5.0 // Increase multiplier for more evident depletion
   ```

2. Override depletion rate in configuration:
   ```go
   // In config.json
   "chemical": {
     "depletionRate": 1.0 // Set higher depletion rate
   }
   ```

## Phase 2: Visual Enhancements

Now that the core functionality works, we need to implement visual feedback to help users understand the system.

### Step 1: Update Chemical Source Visualization

1. Modify renderer to visualize source energy levels:
   ```go
   // In pkg/renderer/renderer.go - update drawChemicalSources method
   func (r *Renderer) drawChemicalSources(screen *ebiten.Image) {
       sources := r.World.GetChemicalSources()
       
       for _, source := range sources {
           // Only draw active sources or sources that have just become inactive
           if !source.IsActive && source.Energy <= 0 {
               continue
           }
           
           // Convert to screen coordinates
           screenX, screenY := r.worldToScreen(source.Position)
           
           // Calculate size based on strength and energy level
           size := 5.0 + math.Sqrt(source.Strength) * 0.5
           
           // Scale size with energy level
           energyRatio := source.Energy / source.MaxEnergy
           size *= 0.5 + 0.5 * energyRatio
           
           // Determine color based on energy level
           alpha := uint8(255 * math.Max(0.2, energyRatio))
           
           // Draw the source
           ebitenutil.DrawCircle(screen, screenX, screenY, size, color.RGBA{255, 200, 0, alpha})
       }
   }
   ```

### Step 2: Add Pulse Effect for Low Energy Sources

1. Implement pulse effect for sources with low energy:
   ```go
   // In pkg/renderer/renderer.go - within drawChemicalSources function
   
   // Pulse effect for low energy
   if energyRatio < 0.3 {
       pulseFreq := 3.0 + (0.3 - energyRatio) * 10.0 // Increase frequency as energy decreases
       pulseAmp := 0.3 * (0.3 - energyRatio) / 0.3   // Amplitude proportional to how low energy is
       pulse := 1.0 + pulseAmp * math.Sin(r.Simulator.Time * pulseFreq)
       size *= pulse
   }
   ```

### Step 3: Add Visual Effect for Source Creation

1. Add a data structure to track source creation events:
   ```go
   // In pkg/renderer/renderer.go
   type SourceCreationEvent struct {
       Position types.Point
       TimeLeft float64
       Strength float64
   }
   
   // Add to Renderer struct
   type Renderer struct {
       // Existing fields
       sourceCreationEvents []SourceCreationEvent
   }
   ```

2. Implement the method to add creation events:
   ```go
   func (r *Renderer) AddSourceCreationEvent(position types.Point, strength float64) {
       r.sourceCreationEvents = append(r.sourceCreationEvents, SourceCreationEvent{
           Position: position,
           TimeLeft: 2.0, // 2 second effect
           Strength: strength,
       })
   }
   ```

3. Draw and update the events in the render loop:
   ```go
   func (r *Renderer) updateAndDrawSourceCreationEvents(screen *ebiten.Image, deltaTime float64) {
       // Update events
       updatedEvents := make([]SourceCreationEvent, 0, len(r.sourceCreationEvents))
       for _, event := range r.sourceCreationEvents {
           // Reduce time left
           event.TimeLeft -= deltaTime
           
           // Keep if not expired
           if event.TimeLeft > 0 {
               updatedEvents = append(updatedEvents, event)
               
               // Draw effect
               screenX, screenY := r.worldToScreen(event.Position)
               
               // Calculate size based on remaining time (starts large, shrinks)
               maxSize := 20.0 + math.Sqrt(event.Strength) * 0.5
               timeRatio := event.TimeLeft / 2.0 // Normalized 0-1
               size := maxSize * (1.0 - timeRatio) // Expands as time progresses
               
               // Calculate alpha (starts transparent, becomes opaque)
               alpha := uint8(255 * (1.0 - timeRatio))
               
               // Draw expanding circle
               ebitenutil.DrawCircle(screen, screenX, screenY, size, 
                   color.RGBA{255, 220, 100, alpha})
           }
       }
       
       // Update the slice
       r.sourceCreationEvents = updatedEvents
   }
   ```

4. Connect to source creation logic in the world:
   ```go
   // In pkg/world/world.go - modify CreateChemicalSource to notify renderer
   func (w *World) CreateChemicalSource(rng *rand.Rand) {
       // Existing creation logic...
       
       // Notify renderer about creation
       if added && w.renderer != nil {
           w.renderer.AddSourceCreationEvent(source.Position, source.Strength)
       }
   }
   ```

## Phase 3: Final Testing and Tuning

### Step 1: Parameter Tuning

Adjust parameters to ensure balanced gameplay:

```json
// In config.json
"chemical": {
  "count": 10,
  "minStrength": 100,
  "maxStrength": 500,
  "minDecayFactor": 0.001,
  "maxDecayFactor": 0.01,
  "depletionRate": 0.5,
  "regenerationProbability": 0.2,
  "targetSystemEnergy": 10000
}
```

### Step 2: Manual Testing

1. Run the simulation and verify:
   - Chemical sources visibly deplete as organisms consume energy
   - Depleted sources are less attractive to organisms
   - New sources appear at appropriate intervals
   - Organisms migrate to new sources when old ones deplete

### Step 3: Performance Profiling

If needed, profile the simulation to ensure the energy balance system doesn't impact performance:

```go
// Add profiling code
import "runtime/pprof"

// In main.go
func main() {
    f, _ := os.Create("cpu.prof")
    pprof.StartCPUProfile(f)
    defer pprof.StopCPUProfile()
    
    // Run simulation...
}
```

## Conclusion

By following these steps, we'll have a complete chemical energy balance system with:
- Depleting sources that encourage organism migration
- Visual feedback for source energy levels
- Effects for source creation and depletion
- Balanced gameplay through parameter tuning

This system will create a more dynamic environment where organisms must adapt to changing conditions, encouraging emergent behaviors and evolution strategies. 