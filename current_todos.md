# Current TODOs

## Performance Optimization

### Diagnosis
- [ ] Install profiling tools
  - [ ] `go install github.com/google/pprof@latest`
- [ ] Run simulation with CPU profiling
  - [ ] `go run -cpuprofile=cpu.prof main.go`
- [ ] Analyze CPU profile
  - [ ] `go tool pprof cpu.prof`
  - [ ] Identify top CPU consumers
  - [ ] Note functions with highest cumulative time
- [ ] Check for common bottlenecks
  - [ ] Identify excessive logging in hot paths
  - [ ] Find frequent concentration grid recalculations
  - [ ] Look for lock contention in UpdateChemicalSources
  - [ ] Identify expensive per-frame calculations

### Code Optimizations
- [ ] Remove Debug Logging
  - [ ] Remove debug prints from `DepleteEnergyFromSourcesAt`
  - [ ] Remove debug prints from `UpdateChemicalSources`
  - [ ] Remove counter variables (`depletionCounter`, `updateCounter`)
  - [ ] Clean up debug prints from `Update` method
- [ ] Optimize Concentration Grid
  - [ ] Increase grid resolution for better performance
  - [ ] Add energy change threshold before invalidating grid
  - [ ] Implement dirty region tracking for grid
  - [ ] Skip concentration calculations for distant sources
- [ ] Reduce Lock Contention
  - [ ] Replace global mutex with finer-grained locks
  - [ ] Implement read/write locks where appropriate
  - [ ] Reduce critical section sizes in hot functions
- [ ] Optimize Chemical Source Updates
  - [ ] Add update interval parameter to config
  - [ ] Implement spatial partitioning for better lookups
  - [ ] Skip far-away sources in concentration calculations

### Advanced Optimizations (if needed)
- [ ] Implement multi-threading for independent calculations
- [ ] Optimize memory usage and reduce allocations
- [ ] Profile and optimize rendering code

### Testing and Measurement
- [ ] Measure baseline performance
- [ ] Document improvements after each optimization
- [ ] Final comprehensive testing with various configurations

## Visual Enhancements
- [ ] Update renderer to visualize source energy levels
- [ ] Implement pulse effect for low-energy sources
- [ ] Add visual effect for source creation
- [ ] Test visual feedback with various energy levels

## Chemical Source Visual Verification
- [ ] Chemical sources visibly deplete as organisms consume from them
- [ ] Depleted sources are less attractive to organisms
- [ ] New sources appear at appropriate intervals
- [ ] Visual effects clearly communicate source energy levels 