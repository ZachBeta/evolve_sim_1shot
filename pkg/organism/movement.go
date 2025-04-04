package organism

import (
	"math"

	"github.com/zachbeta/evolve_sim/pkg/types"
)

// Move updates the organism's position based on its heading and speed
// It handles boundary collisions and adjusts the position and heading accordingly
func Move(org *types.Organism, bounds types.Rect, deltaTime float64) {
	// Store previous heading before updating
	org.PreviousHeading = org.Heading

	// Calculate the distance to move based on speed and time delta
	distance := org.Speed * deltaTime

	// Store the original position to restore if needed
	originalPos := org.Position

	// Move the organism forward based on heading and speed
	dx := math.Cos(org.Heading) * distance
	dy := math.Sin(org.Heading) * distance
	newPos := types.Point{X: originalPos.X + dx, Y: originalPos.Y + dy}

	// Calculate actual distance moved for energy consumption
	distanceMoved := math.Sqrt(dx*dx + dy*dy)

	// Consume energy based on distance moved and speed
	// Faster organisms use more energy per unit distance
	// Use the organism's MovementCost parameter modified by EnergyEfficiency
	energyCost := distanceMoved * org.MovementCost * org.EnergyEfficiency * (1.0 + org.Speed*0.05)
	org.Energy -= energyCost

	// If energy is depleted, reduce speed proportionally
	if org.Energy <= 0 {
		org.Energy = 0
		distance = 0 // Stop movement when out of energy
		newPos = originalPos
	} else if org.Energy < org.EnergyCapacity*0.1 {
		// Reduce speed when low on energy (less than 10% of capacity)
		energyRatio := org.Energy / (org.EnergyCapacity * 0.1)
		distance *= energyRatio
		dx = math.Cos(org.Heading) * distance
		dy = math.Sin(org.Heading) * distance
		newPos = types.Point{X: originalPos.X + dx, Y: originalPos.Y + dy}
	}

	// Check if the new position is within bounds
	if newPos.X < bounds.Min.X || newPos.X >= bounds.Max.X ||
		newPos.Y < bounds.Min.Y || newPos.Y >= bounds.Max.Y {
		// Calculate new heading based on which boundary was hit
		newHeading := org.Heading

		// Check for horizontal boundary collision
		if newPos.X < bounds.Min.X || newPos.X >= bounds.Max.X {
			// Hit left or right wall, reflect horizontally
			newHeading = math.Pi - org.Heading
			if newHeading < 0 {
				newHeading += 2 * math.Pi
			}
		}

		// Check for vertical boundary collision
		if newPos.Y < bounds.Min.Y || newPos.Y >= bounds.Max.Y {
			// Hit top or bottom wall, reflect vertically
			newHeading = 2*math.Pi - org.Heading
		}

		// Update the heading
		org.Heading = newHeading

		// Keep organism within bounds
		boundedX := math.Max(bounds.Min.X, math.Min(newPos.X, bounds.Max.X-0.001))
		boundedY := math.Max(bounds.Min.Y, math.Min(newPos.Y, bounds.Max.Y-0.001))
		org.Position = types.Point{X: boundedX, Y: boundedY}
	} else {
		// No collision, update position normally
		org.Position = newPos
	}

	// Update the organism's trail
	org.UpdateTrail()

	// Ensure we take the shortest path for rotation (for smooth animation)
	for org.Heading-org.PreviousHeading > math.Pi {
		org.PreviousHeading += 2 * math.Pi
	}
	for org.PreviousHeading-org.Heading > math.Pi {
		org.PreviousHeading -= 2 * math.Pi
	}

	// Update reproduction timer
	org.TimeSinceReproduction += deltaTime
}
