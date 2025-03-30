# Evolutionary Simulator

A 2D simulation of single-cell organisms responding to chemical gradients in their environment. The simulation focuses on emergent behavior through individual organism preferences and directional sensing.

## Features

- Single-cell organisms with directional sensing (front, left, right)
- Static chemical sources creating concentration gradients
- Organisms have individual chemical preferences (normally distributed)
- Greedy movement algorithm toward preferred concentration
- Visualization with contour lines showing chemical gradients

## Building and Running

```bash
# Build the project
go build -o evolve_sim ./cmd/evolve_sim

# Run the simulator
./evolve_sim
```

## Project Structure

- `cmd/evolve_sim`: Main application entry point
- `pkg/types`: Core data structures
- `pkg/config`: Configuration system
- `pkg/world`: World and chemical gradient system
- `pkg/organism`: Organism behavior and movement
- `pkg/simulation`: Simulation engine
- `pkg/renderer`: Visualization system