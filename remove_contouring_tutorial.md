# Remove Contouring Functionality - Tutorial Checklist

## Overview
This tutorial will guide you through removing the contouring functionality from the simulation as it's currently slow and not useful. By eliminating this feature, we'll improve performance and simplify the codebase.

## Steps

### 1. Identify Contouring Code
- [ ] Find the `ContourLine` struct in `pkg/world/concentration_grid.go`
- [ ] Locate the `GenerateContourLines` function
- [ ] Identify other contouring-related functions/methods (like `marchingSquares`, `segmentsToContours`)

### 2. Remove Contouring Code from Concentration Grid
- [ ] Comment out or remove the `ContourLine` struct definition
- [ ] Remove the `Direction`, `Cell`, and `Segment` types (if they're only used for contouring)
- [ ] Comment out or remove the `GenerateContourLines` function
- [ ] Remove the `marchingSquares` function
- [ ] Remove the `segmentsToContours` function and any other contouring-specific helpers

### 3. Update Renderer Code
- [ ] Find contouring visualization code in `pkg/renderer/renderer.go`
- [ ] Comment out or remove code that draws contour lines
- [ ] If there are contouring-specific configuration options, disable them by default

### 4. Clean Up Configuration
- [ ] Check for contouring-related configuration in the `config.json` file
- [ ] Remove or set contouring configuration options to disabled

### 5. Update Documentation
- [ ] Update any comments/documentation that reference contouring functionality
- [ ] Add a note to class or function comments indicating the feature has been removed

### 6. Testing
- [ ] Run the simulation to verify it still works after removing contouring
- [ ] Check that performance has improved
- [ ] Ensure no visual artifacts remain from the removed contouring

## Expected Benefits
- ✅ Improved performance
- ✅ Simplified codebase
- ✅ Reduced memory usage
- ✅ Cleaner visualization

## Note
If you need to restore contouring functionality in the future, consider:
1. Creating a separate branch before removing the code
2. Using conditional compilation or feature flags if you want to keep the code but disable it
3. Implementing a more efficient contouring algorithm if needed later 