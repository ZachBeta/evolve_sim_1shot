# Organism Trails Feature Guide

This guide explains the new organism trails feature that visualizes the movement patterns of organisms in the evolutionary simulator.

## Overview

The organism trails feature allows you to see the recent path an organism has traveled, making it easier to understand:
- Movement patterns in response to chemical gradients
- Turning behavior and navigation strategies
- Differences in behavior between organisms with different preferences

## How to Use

1. **Toggle Trails**: Press the `T` key to toggle organism trails on and off.
2. **Visualization**: When enabled, each organism will display a fading trail showing its recent movement history.
3. **Color Coding**: Trails match the color of the organism, which is determined by its chemical preference. This makes it easy to track specific organisms and compare behaviors.

## Implementation Details

Organism trails are implemented with the following components:

1. **Position History**:
   - Each organism stores a list of its recent positions in `PositionHistory`.
   - The history is capped at a maximum length (currently 30 positions) to limit memory usage.
   - Positions are recorded every few updates rather than every frame to create an appropriately spaced trail.

2. **Trail Rendering**:
   - Trails are drawn as connected line segments between historical positions.
   - Opacity increases from the oldest to the newest positions, creating a fading effect.
   - The most recent segment connects the last recorded position to the current position.

3. **Performance Optimization**:
   - The update frequency is controlled to balance detail against performance.
   - Only a limited number of historical positions are stored to minimize memory usage.

## Tips for Using Trails

- **Observe Chemical Seeking**: With both contour lines and trails enabled, you can observe how organisms navigate towards their preferred chemical concentration.
- **Compare Different Organisms**: Look for differences in trail patterns between organisms with different preferences (red vs. blue).
- **Analyze Turning Behavior**: Trails make it easy to see how organisms adjust their direction, especially when encountering boundaries or chemical gradients.
- **Study Emergent Patterns**: Over time, you may notice emergent movement patterns as organisms explore the environment.

## Next Steps

Now that organism trails are implemented, we plan to add:

1. **Organism Selection**: The ability to click on organisms to select them and view detailed information.
2. **Enhanced Analytics**: Statistical analysis of movement patterns and behavior.
3. **Customizable Trails**: Options to adjust trail length, appearance, and recording frequency.

## Key Commands

| Key | Function |
|-----|----------|
| T | Toggle organism trails on/off |
| G | Toggle grid on/off |
| C | Toggle chemical concentration visualization on/off |
| O | Toggle contour lines on/off |
| S | Toggle sensor visualization on/off |
| Space | Pause/resume simulation |
| R | Reset simulation |
| M | Cycle through color schemes |
| +/- | Adjust simulation speed |

## Feedback

We welcome your feedback on the trails feature. Is the trail length appropriate? Are the colors visible enough? Let us know how this feature enhances your understanding of the simulation! 