# Contouring Removal Todo List

## Background
- [x] The contouring functionality is slow and not useful
- [x] Removing it will improve performance and simplify the codebase

## Code Removal Tasks

### Concentration Grid
- [x] Remove `ContourLine` struct in `pkg/world/concentration_grid.go`
- [x] Remove `Direction` type
- [x] Remove `Cell` type
- [x] Remove `Segment` type
- [x] Remove `GenerateContourLines` function
- [x] Remove `marchingSquares` function
- [x] Remove `segmentsToContours` function
- [x] Remove any other contouring-specific helpers

### Renderer
- [x] Remove contouring visualization code in `pkg/renderer/renderer.go`
- [x] Remove code that draws contour lines
- [x] Disable any contouring-specific configuration options

### Configuration
- [x] Remove contouring-related configuration in `config.json`
- [x] Set any remaining contouring options to disabled

## Documentation
- [x] Update comments/documentation that reference contouring
- [x] Add notes in class/function comments indicating feature removal

## Testing
- [x] Run simulation to verify it works after removal
- [x] Verify performance improvements
- [x] Check for any remaining visual artifacts

## Expected Benefits
- [x] Improved performance
- [x] Simplified codebase
- [x] Reduced memory usage
- [x] Cleaner visualization

## Notes for Future Reference
- [x] Consider creating a branch before removal if restoration might be needed
- [x] If reimplementing later, use a more efficient algorithm 