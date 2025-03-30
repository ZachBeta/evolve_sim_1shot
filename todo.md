# Chemical Energy Balance System: Test-Driven Implementation Checklist

## Phase 1: Chemical Source Energy Implementation
- [x] Update `ChemicalSource` struct with new energy fields
  - [x] Add `Energy`, `MaxEnergy`, `DepletionRate`, and `IsActive` fields
  - [x] Update constructor to initialize these fields
- [x] Write test for chemical source initialization
- [x] Modify `GetConcentrationAt` to factor in energy level
- [x] Write test for concentration scaling with energy
- [x] Run tests and verify both tests pass

## Phase 2: Energy Depletion System
- [x] Add `Update` method to `ChemicalSource` to handle basic depletion
- [x] Write test for basic depletion over time
- [x] Implement source deactivation when energy is depleted
- [x] Write test for source deactivation
- [x] Implement organism energy consumption tracking in behavior code
- [x] Add `DepleteEnergyFromSourcesAt` method to the `World` struct
- [x] Write test for organism-based energy depletion
- [x] Run all depletion tests and verify they pass

## Phase 3: System Energy Management
- [x] Add `totalSystemEnergy` and `targetSystemEnergy` fields to `World` struct
- [x] Update `World` constructor to calculate initial system energy
- [x] Write test for system energy tracking
- [x] Implement energy tracking in source depletion functions
- [x] Run energy tracking tests and verify they pass

## Phase 4: Source Regeneration
- [x] Implement `CreateChemicalSource` method for the `World`
- [x] Write test for source creation based on energy deficit
- [x] Add `UpdateChemicalSources` method to handle source updates and creation
- [x] Write test for source creation probability
- [x] Run source regeneration tests and verify they pass

## Phase 5: Integration
- [ ] Update `Simulator` struct to include RNG
- [ ] Update `Step` method to call `UpdateChemicalSources`
- [ ] Write integration test for full simulation cycle
- [ ] Run integration tests and verify energy balance is maintained
- [ ] Manual testing: observe organism migration as sources deplete

## Phase 6: Visual Enhancements
- [ ] Update renderer to visualize source energy levels
- [ ] Implement pulse effect for low-energy sources
- [ ] Add visual effect for source creation
- [ ] Test visual feedback with various energy levels

## Final Verification
- [ ] Chemical sources visibly deplete as organisms consume from them
- [ ] Depleted sources are less attractive to organisms
- [ ] New sources appear at appropriate intervals
- [ ] The simulation maintains interesting dynamics
- [ ] Visual effects clearly communicate source energy levels

## Configuration
- [x] Add configurable parameters to config.json
  - [x] `depletionRate`
  - [x] `regenerationProbability`
  - [x] `targetSystemEnergy` 