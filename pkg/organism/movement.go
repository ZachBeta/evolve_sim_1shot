package organism

import (
	"github.com/zachbeta/evolve_sim/pkg/types"
)

// Move updates the organism's position based on its heading and speed
// It handles boundary collisions and adjusts the position and heading accordingly
func Move(org *types.Organism, bounds types.Rect, deltaTime float64) {
	// Calculate the distance to move based on speed and time delta
	distance := org.Speed * deltaTime

	// Store the original position to restore if needed
	originalPos := org.Position

	// Move the organism forward
	org.MoveForward(distance)

	// Check if the new position is within bounds
	if !bounds.Contains(org.Position) {
		// Restore original position
		org.Position = originalPos

		// Adjust heading - bounce off the wall by reflecting the angle
		// Determine which wall was hit
		nextPos := originalPos
		nextPos.X += distance * 2 * 0.5 * 1
		if nextPos.X < bounds.Min.X || nextPos.X > bounds.Max.X {
			// Hit left or right wall, reflect X component of heading
			org.Turn(2 * (0 - org.Heading))
		} else {
			// Hit top or bottom wall, reflect Y component of heading
			org.Turn(2 * (1.5708 - org.Heading)) // 1.5708 radians = 90 degrees
		}

		// Move a small distance in the new direction to avoid getting stuck
		org.MoveForward(distance * 0.1)
	}
}
