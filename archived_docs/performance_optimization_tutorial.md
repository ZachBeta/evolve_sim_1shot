# Performance Optimization Tutorial: Improving FPS in the Simulation

This tutorial outlines steps to diagnose and fix performance issues in the evolve_sim application to improve the frames per second (FPS) rate.

## 1. Diagnosis Phase

### 1.1 Identifying Performance Bottlenecks

```bash
# Install profiling tools if needed
go install github.com/google/pprof@latest

# Run the simulation with CPU profiling
go run -cpuprofile=cpu.prof main.go

# Analyze the profile
go tool pprof cpu.prof
```

### 1.2 Common Bottlenecks to Look For

- **Excessive Logging**: Debug statements (especially in hot paths)
- **Frequent Grid Invalidation**: Recalculating the concentration grid
- **Lock Contention**: Mutex usage in UpdateChemicalSources and other methods
- **Per-frame Calculations**: Expensive operations performed each frame

## 2. Optimization Steps

### 2.1 Remove Debug Logging

First, disable the extensive logging we added for debugging:

```go
// Remove or comment out debug print statements in DepleteEnergyFromSourcesAt
// Remove or comment out debug print statements in UpdateChemicalSources
// Remove or comment out debug print statements in ChemicalSource.Update
```

### 2.2 Optimize Concentration Grid Usage

```go
// Reduce grid resolution to improve performance
world.InitializeConcentrationGrid(10.0) // Increase from 5.0 to 10.0

// Only invalidate grid when necessary
// Add a threshold for invalidation based on energy change
```

### 2.3 Reduce Lock Contention

```go
// Use finer-grained locks where possible
// Consider read/write separation for hot paths
// Batch updates where possible
```

### 2.4 Optimize Chemical Source Updates

```go
// Adjust update frequency - not every source needs updating every frame
// Use spatial partitioning for collision detection and concentration calculations
// Implement early-exit optimizations for far-away sources
```

### 2.5 Lazy Evaluation Strategies

```go
// Only calculate concentration when needed
// Use dirty flags to track which areas need recalculation
// Skip inactive sources in calculations
```

## 3. Implementation Plan

### 3.1 First Pass: Remove Debug Logging

Remove all debug print statements from these files:
- pkg/world/world.go
- pkg/types/chemical_source.go

### 3.2 Second Pass: Optimize Grid and Source Updates

Modify the concentration grid:
- Increase grid cell size
- Implement lazy grid updates
- Add threshold for grid invalidation

### 3.3 Third Pass: Advanced Optimizations

If necessary:
- Implement spatial partitioning
- Add multi-threading for independent calculations
- Optimize memory usage and reduce allocations

## 4. Testing & Measuring Improvements

After each change:

```bash
# Measure FPS before and after changes
go run main.go

# Run profiling again to verify improvements
go run -cpuprofile=cpu_after.prof main.go
go tool pprof -http=:8080 cpu_after.prof
```

## 5. Configuration Tweaks

If code optimizations aren't sufficient:

```json
{
  "chemical": {
    "count": 5,          // Consider reducing if necessary
    "depletionRate": 10.0, // Balance between simulation accuracy and performance
    "regenerationProbability": 0.5
  },
  "render": {
    "showConcentration": true,  // Consider disabling for performance
    "showGrid": true            // Consider disabling for performance
  }
}
```

## Conclusion

Performance optimization is an iterative process. Start with the most impactful changes (like removing debug logging) and gradually implement more complex optimizations as needed while measuring the impact of each change.

Remember that some trade-offs between simulation accuracy and performance may be necessary for real-time visualization. 