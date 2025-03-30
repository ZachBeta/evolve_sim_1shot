# Concentration Removal Todo List

## Background
- [x] The concentration grid functionality is slow and not essential
- [x] Removing it will improve performance and simplify the codebase

## Understanding Concentration Usage
- [x] Review how ConcentrationGrid is used in the World
- [x] Identify which organism behaviors depend on concentration
- [x] Note which visualization features use concentration data

## Modifying Core Functionality
- [x] Keep the ConcentrationGrid struct but simplify its implementation
- [x] Modify GetConcentrationAt to return minimal data
- [x] Update GetGradientAt to return sensible default values
- [x] Update organism movement to not rely on concentration gradients
- [x] Simplify sensor functionality to work without detailed concentration data
- [x] Ensure organisms can still navigate toward chemical sources directly

## Renderer Updates
- [x] Disable drawChemicalConcentration method
- [x] Remove concentration visualization code
- [x] Remove concentration legend rendering

## Configuration Changes
- [x] Set ShowConcentration to false by default
- [x] Consider removing the option entirely if no longer needed

## UI Cleanup
- [x] Remove concentration toggle key (C)
- [x] Update help text and controls documentation
- [x] Remove legend display related to concentration

## Testing
- [x] Run the simulation to verify it still works without concentration grid
- [x] Check that performance has improved
- [x] Ensure organisms still behave in a meaningful way

## Expected Benefits
- [x] Significantly improved performance
- [x] Reduced memory usage
- [x] Simplified codebase
- [x] More focused simulation behavior

## Implementation Notes
- [x] Use simplified placeholder if concentration is deeply integrated
- [x] Prioritize making organisms function without detailed gradient information
- [x] Keep a minimal API surface to avoid breaking many dependencies 