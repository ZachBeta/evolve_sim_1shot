package renderer

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/zachbeta/evolve_sim/pkg/config"
	"github.com/zachbeta/evolve_sim/pkg/simulation"
	"github.com/zachbeta/evolve_sim/pkg/types"
	"github.com/zachbeta/evolve_sim/pkg/world"
)

// ReproductionEvent tracks visual effects for organism reproduction
type ReproductionEvent struct {
	Position types.Point // Position of reproduction
	TimeLeft float64     // Time left for this effect (in seconds)
}

// Renderer is responsible for visualizing the simulation
type Renderer struct {
	World               *world.World
	Simulator           *simulation.Simulator
	Config              config.SimulationConfig
	WindowWidth         int
	WindowHeight        int
	ShowGrid            bool
	ShowSensors         bool
	ShowLegend          bool
	ShowTrails          bool
	Stats               simulation.SimulationStats
	FPS                 float64
	keyStates           map[ebiten.Key]bool
	CurrentColorScheme  ColorScheme
	ColorSchemes        []ColorScheme
	CurrentSchemeIndex  int
	interpolationFactor float64 // For smooth animations between frames
	triangleImage       *ebiten.Image
	triangleOpts        ebiten.DrawImageOptions
	selectedOrganism    *types.Organism     // For future organism selection feature
	reproductionEvents  []ReproductionEvent // Track reproduction visual effects
	previousOrgCount    int                 // To detect reproduction events
}

// NewRenderer creates a new renderer with the given configuration
func NewRenderer(world *world.World, simulator *simulation.Simulator, config config.SimulationConfig) *Renderer {
	// Initialize available color schemes
	colorSchemes := []ColorScheme{
		ViridisScheme, // Default
		MagmaScheme,
		PlasmaScheme,
		ClassicScheme,
	}

	// Get initial organism count
	initialCount, _ := world.GetPopulationInfo()

	// Create a new renderer
	r := &Renderer{
		World:               world,
		Simulator:           simulator,
		Config:              config,
		WindowWidth:         config.Render.WindowWidth,
		WindowHeight:        config.Render.WindowHeight,
		ShowGrid:            config.Render.ShowGrid,
		ShowSensors:         config.Render.ShowSensors,
		ShowLegend:          config.Render.ShowLegend,
		ShowTrails:          false, // Default to false, can be toggled
		FPS:                 0,
		keyStates:           make(map[ebiten.Key]bool),
		CurrentColorScheme:  ViridisScheme,
		ColorSchemes:        colorSchemes,
		CurrentSchemeIndex:  0,
		interpolationFactor: 1.0, // Default to full interpolation
		reproductionEvents:  make([]ReproductionEvent, 0),
		previousOrgCount:    initialCount,
	}

	// Create a triangle image for optimized drawing
	r.triangleImage = ebiten.NewImage(16, 16)
	r.triangleOpts = ebiten.DrawImageOptions{}

	return r
}

// isKeyJustPressed checks if a key was just pressed this frame
func (r *Renderer) isKeyJustPressed(key ebiten.Key) bool {
	wasPressed := r.keyStates[key]
	isPressed := ebiten.IsKeyPressed(key)
	r.keyStates[key] = isPressed
	return isPressed && !wasPressed
}

// Update is called each frame by Ebiten
func (r *Renderer) Update() error {
	// Update simulator
	r.Simulator.Step()

	// Update stats
	r.Stats = r.Simulator.CollectStats()

	// Update interpolation factor (will be 1.0 most of the time, could be adjusted for smoother animations)
	r.interpolationFactor = 1.0

	// Update FPS
	r.FPS = ebiten.ActualFPS()

	// Update reproduction events
	r.updateReproductionEvents(r.Simulator.TimeStep * r.Simulator.SimulationSpeed)

	// Handle keyboard input - only respond to key presses, not holds
	if r.isKeyJustPressed(ebiten.KeySpace) {
		r.Simulator.SetPaused(!r.Simulator.IsPaused)
	}

	if r.isKeyJustPressed(ebiten.KeyR) {
		r.Simulator.Reset()
	}

	if r.isKeyJustPressed(ebiten.KeyG) {
		r.ShowGrid = !r.ShowGrid
	}

	if r.isKeyJustPressed(ebiten.KeyS) {
		r.ShowSensors = !r.ShowSensors
	}

	if r.isKeyJustPressed(ebiten.KeyL) {
		r.ShowLegend = !r.ShowLegend
	}

	if r.isKeyJustPressed(ebiten.KeyT) {
		r.ShowTrails = !r.ShowTrails
	}

	// Cycle through color schemes
	if r.isKeyJustPressed(ebiten.KeyM) {
		r.CurrentSchemeIndex = (r.CurrentSchemeIndex + 1) % len(r.ColorSchemes)
		r.CurrentColorScheme = r.ColorSchemes[r.CurrentSchemeIndex]
	}

	// Speed control - these can respond continuously
	if ebiten.IsKeyPressed(ebiten.KeyEqual) {
		r.Simulator.SetSimulationSpeed(r.Simulator.SimulationSpeed * 1.1)
	}

	if ebiten.IsKeyPressed(ebiten.KeyMinus) {
		r.Simulator.SetSimulationSpeed(r.Simulator.SimulationSpeed * 0.9)
	}

	return nil
}

// Draw renders the current state to the screen
func (r *Renderer) Draw(screen *ebiten.Image) {
	// Clear screen
	screen.Fill(color.RGBA{20, 20, 30, 255})

	// Draw grid if enabled
	if r.ShowGrid {
		r.drawGrid(screen)
	}

	// Chemical concentration visualization has been disabled for performance reasons

	// Draw chemical sources
	r.drawChemicalSources(screen)

	// Draw organisms
	r.drawOrganisms(screen)

	// Draw reproduction events
	r.drawReproductionEvents(screen)

	// Draw statistics
	r.drawStats(screen)

	// Draw concentration legend if enabled - disabled for performance
	// if r.ShowLegend && r.ShowConcentration {
	//	r.drawConcentrationLegend(screen)
	// }
}

// Layout returns the logical screen dimensions
func (r *Renderer) Layout(outsideWidth, outsideHeight int) (int, int) {
	return r.WindowWidth, r.WindowHeight
}

// Helper method to convert world coordinates to screen coordinates
func (r *Renderer) worldToScreen(point types.Point) (float64, float64) {
	bounds := r.World.GetBounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	// Convert world coordinates to normalized coordinates (0-1)
	normalizedX := (point.X - bounds.Min.X) / width
	normalizedY := (point.Y - bounds.Min.Y) / height

	// Convert normalized coordinates to screen coordinates
	screenX := normalizedX * float64(r.WindowWidth)
	screenY := normalizedY * float64(r.WindowHeight)

	return screenX, screenY
}

// Draw a visualization of chemical concentration - removed for performance
func (r *Renderer) drawChemicalConcentration(screen *ebiten.Image) {
	// This method is kept for compatibility but its functionality has been disabled
	// for performance reasons
}

// Draw chemical sources
func (r *Renderer) drawChemicalSources(screen *ebiten.Image) {
	// Get chemical sources
	sources := r.World.GetChemicalSources()

	// Draw each chemical source
	for _, source := range sources {
		// Skip inactive sources
		if !source.IsActive {
			continue
		}

		// Convert world coordinates to screen coordinates
		x, y := r.worldToScreen(source.Position)

		// Calculate size based on source strength
		radius := math.Sqrt(source.Strength) * 0.3
		radius = math.Max(5, math.Min(30, radius)) // Clamp between 5 and 30 pixels

		// Size indicates strength
		// Scale the size based on energy level
		energyRatio := source.Energy / source.MaxEnergy
		sizeModifier := 0.5 + 0.5*energyRatio // 50% - 100% of original size
		radius *= sizeModifier

		// Get color from scheme based on decay factor
		// Higher decay = faster falloff = "hotter" color
		relativeDecay := (source.DecayFactor - 0.001) / (0.01 - 0.001) // Normalized between 0-1
		sourceColor := GetColorFromScheme(r.CurrentColorScheme, 1.0-relativeDecay)

		// Make source more visible by increasing opacity with energy
		sourceColor.A = uint8(200 * energyRatio)

		// Draw filled circle
		for cy := int(y) - int(radius); cy <= int(y)+int(radius); cy++ {
			for cx := int(x) - int(radius); cx <= int(x)+int(radius); cx++ {
				dx := float64(cx) - x
				dy := float64(cy) - y
				if dx*dx+dy*dy <= radius*radius {
					if cx >= 0 && cx < r.WindowWidth && cy >= 0 && cy < r.WindowHeight {
						screen.Set(cx, cy, sourceColor)
					}
				}
			}
		}

		// Draw outline
		outlineColor := color.RGBA{255, 255, 255, 200}
		for angle := 0.0; angle < 2*math.Pi; angle += 0.01 {
			cx := int(x + math.Cos(angle)*radius)
			cy := int(y + math.Sin(angle)*radius)
			if cx >= 0 && cx < r.WindowWidth && cy >= 0 && cy < r.WindowHeight {
				screen.Set(cx, cy, outlineColor)
			}
		}
	}
}

// Draw organisms
func (r *Renderer) drawOrganisms(screen *ebiten.Image) {
	organisms := r.World.GetOrganisms()

	for _, org := range organisms {
		// Convert world coordinates to screen coordinates
		screenX, screenY := r.worldToScreen(org.Position)

		// Determine base color based on chemical preference
		// Map preference to a blue-to-red gradient
		prefRange := r.Config.Organism.PreferenceDistributionMean * 3
		normalizedPref := org.ChemPreference / prefRange

		baseRed := uint8(normalizedPref * 255)
		baseBlue := uint8((1 - normalizedPref) * 255)
		baseGreen := uint8(128 - math.Abs(float64(normalizedPref*255-128)))

		// Modify color based on energy level
		// Low energy organisms appear darker/more transparent
		energyRatio := org.Energy / org.EnergyCapacity
		red := uint8(float64(baseRed) * math.Sqrt(energyRatio))
		green := uint8(float64(baseGreen) * math.Sqrt(energyRatio))
		blue := uint8(float64(baseBlue) * math.Sqrt(energyRatio))

		// Full alpha for the organism itself
		alpha := uint8(255)

		// Draw trail if enabled
		if r.ShowTrails && len(org.PositionHistory) > 1 {
			// Draw a line connecting all positions in history
			trailColor := color.RGBA{red, green, blue, 100} // Semi-transparent

			// Draw lines between consecutive points
			for i := 0; i < len(org.PositionHistory)-1; i++ {
				// Convert world coordinates to screen coordinates for both points
				x1, y1 := r.worldToScreen(org.PositionHistory[i])
				x2, y2 := r.worldToScreen(org.PositionHistory[i+1])

				// Fade the trail as it gets older
				trailAlpha := uint8(40 + (160 * i / len(org.PositionHistory)))
				fadedColor := color.RGBA{red, green, blue, trailAlpha}

				// Draw the line
				ebitenutil.DrawLine(screen, x1, y1, x2, y2, fadedColor)
			}

			// Connect the last history point to current position
			if len(org.PositionHistory) > 0 {
				lastX, lastY := r.worldToScreen(org.PositionHistory[len(org.PositionHistory)-1])
				ebitenutil.DrawLine(screen, lastX, lastY, screenX, screenY, trailColor)
			}
		}

		// Calculate the visual heading with interpolation for smooth rotation
		visualHeading := org.PreviousHeading + (org.Heading-org.PreviousHeading)*r.interpolationFactor

		// Define triangle size (can be adjusted based on organism properties)
		// Scale size slightly with energy level for visual feedback
		sizeMultiplier := 0.8 + 0.4*energyRatio // Size reduced by up to 20% when low energy
		size := 4.0 * sizeMultiplier

		// Calculate triangle vertices
		// The triangle should point in the direction of heading
		// First point: front of the triangle (in heading direction)
		frontX := screenX + math.Cos(visualHeading)*size*1.5
		frontY := screenY + math.Sin(visualHeading)*size*1.5

		// Calculate the back corners (perpendicular to heading)
		backOffsetX := math.Cos(visualHeading+math.Pi/2) * size
		backOffsetY := math.Sin(visualHeading+math.Pi/2) * size

		// Left back corner
		leftX := screenX - math.Cos(visualHeading)*size/2 - backOffsetX
		leftY := screenY - math.Sin(visualHeading)*size/2 - backOffsetY

		// Right back corner
		rightX := screenX - math.Cos(visualHeading)*size/2 + backOffsetX
		rightY := screenY - math.Sin(visualHeading)*size/2 + backOffsetY

		// Draw the triangle
		r.drawTriangle(screen, frontX, frontY, leftX, leftY, rightX, rightY,
			color.RGBA{red, green, blue, alpha})

		// Add a border for better visibility
		borderAlpha := uint8(150 + 50*energyRatio) // Border fades a bit when low energy
		ebitenutil.DrawLine(screen, frontX, frontY, leftX, leftY, color.RGBA{255, 255, 255, borderAlpha})
		ebitenutil.DrawLine(screen, leftX, leftY, rightX, rightY, color.RGBA{255, 255, 255, borderAlpha})
		ebitenutil.DrawLine(screen, rightX, rightY, frontX, frontY, color.RGBA{255, 255, 255, borderAlpha})

		// Draw energy indicator (small bar above organism)
		if energyRatio < 0.99 { // Only draw when not full
			barWidth := 8.0
			barHeight := 1.5
			barX := screenX - barWidth/2
			barY := screenY - size*2

			// Background (empty) bar
			ebitenutil.DrawRect(screen, barX, barY, barWidth, barHeight, color.RGBA{40, 40, 40, 200})

			// Filled portion based on energy
			fillWidth := barWidth * energyRatio

			// Color goes from red (low) to green (high)
			barRed := uint8(255 * (1 - energyRatio))
			barGreen := uint8(255 * energyRatio)
			ebitenutil.DrawRect(screen, barX, barY, fillWidth, barHeight, color.RGBA{barRed, barGreen, 0, 230})
		}

		// Draw sensors if enabled
		if r.ShowSensors {
			sensorPositions := org.GetSensorPositions(r.Config.Organism.SensorDistance)

			// Draw lines to sensors
			for _, sensorPos := range sensorPositions {
				sensorX, sensorY := r.worldToScreen(sensorPos)
				ebitenutil.DrawLine(screen, screenX, screenY, sensorX, sensorY, color.RGBA{200, 200, 200, 128})
			}
		}
	}
}

// Draw statistics on screen
func (r *Renderer) drawStats(screen *ebiten.Image) {
	stats := []string{
		fmt.Sprintf("FPS: %.1f", r.FPS),
		fmt.Sprintf("Time: %.2f", r.Simulator.Time),
		fmt.Sprintf("Organisms: %d", r.Stats.Organisms.Count),
		fmt.Sprintf("Speed: %.1fx", r.Simulator.SimulationSpeed),
		fmt.Sprintf("Paused: %v", r.Simulator.IsPaused),
		fmt.Sprintf("Avg Preference: %.1f", r.Stats.Organisms.AveragePreference),
		fmt.Sprintf("Avg Energy: %.1f (%.0f%%)",
			r.Stats.Organisms.AverageEnergy,
			r.Stats.Organisms.EnergyRatio*100),
		fmt.Sprintf("Grid: %v", r.ShowGrid),
		fmt.Sprintf("Trails: %v", r.ShowTrails),
	}

	// Draw stats in the top-left corner
	for i, stat := range stats {
		ebitenutil.DebugPrintAt(screen, stat, 10, 20+i*20)
	}

	// Draw controls help
	controls := []string{
		"Space: Pause/Resume",
		"R: Reset",
		"G: Toggle Grid",
		"S: Toggle Sensors",
		"L: Toggle Legend",
		"T: Toggle Trails",
		"M: Cycle Color Schemes",
		"+/-: Adjust Speed",
	}

	// Draw controls in the bottom-left corner
	for i, control := range controls {
		ebitenutil.DebugPrintAt(
			screen,
			control,
			10,
			r.WindowHeight-20*len(controls)+i*20,
		)
	}
}

// Draw a grid for visual reference
func (r *Renderer) drawGrid(screen *ebiten.Image) {
	bounds := r.World.GetBounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	// Define grid cell size in world coordinates
	gridCellSize := 50.0 // World units per grid cell

	// Calculate number of grid lines
	numLinesX := int(width/gridCellSize) + 1
	numLinesY := int(height/gridCellSize) + 1

	// Draw vertical grid lines
	for i := 0; i < numLinesX; i++ {
		worldX := bounds.Min.X + float64(i)*gridCellSize
		startX, startY := r.worldToScreen(types.Point{X: worldX, Y: bounds.Min.Y})
		endX, endY := r.worldToScreen(types.Point{X: worldX, Y: bounds.Max.Y})
		ebitenutil.DrawLine(screen, startX, startY, endX, endY, color.RGBA{120, 120, 140, 180}) // Brighter and more opaque
	}

	// Draw horizontal grid lines
	for i := 0; i < numLinesY; i++ {
		worldY := bounds.Min.Y + float64(i)*gridCellSize
		startX, startY := r.worldToScreen(types.Point{X: bounds.Min.X, Y: worldY})
		endX, endY := r.worldToScreen(types.Point{X: bounds.Max.X, Y: worldY})
		ebitenutil.DrawLine(screen, startX, startY, endX, endY, color.RGBA{120, 120, 140, 180}) // Brighter and more opaque
	}
}

// Draw a triangle with the specified points and color
func (r *Renderer) drawTriangle(screen *ebiten.Image, x1, y1, x2, y2, x3, y3 float64, clr color.Color) {
	// Find the bounding box of the triangle
	minX := math.Min(x1, math.Min(x2, x3))
	maxX := math.Max(x1, math.Max(x2, x3))
	minY := math.Min(y1, math.Min(y2, y3))
	maxY := math.Max(y1, math.Max(y2, y3))

	// Iterate over each pixel in the bounding box
	for y := int(minY); y <= int(maxY); y++ {
		for x := int(minX); x <= int(maxX); x++ {
			// Check if the point is inside the triangle
			if pointInTriangle(float64(x), float64(y), x1, y1, x2, y2, x3, y3) {
				screen.Set(x, y, clr)
			}
		}
	}
}

// pointInTriangle determines if a point is inside a triangle using barycentric coordinates
func pointInTriangle(px, py, x1, y1, x2, y2, x3, y3 float64) bool {
	// Calculate area of the full triangle
	area := 0.5 * math.Abs((x2-x1)*(y3-y1)-(x3-x1)*(y2-y1))
	if area < 0.00001 {
		return false // Degenerate triangle
	}

	// Calculate barycentric coordinates
	alpha := 0.5 * math.Abs((x2-x3)*(py-y3)-(y2-y3)*(px-x3)) / area
	beta := 0.5 * math.Abs((x3-x1)*(py-y1)-(y3-y1)*(px-x1)) / area
	gamma := 1.0 - alpha - beta

	// Point is in triangle if all coordinates are between 0 and 1
	return alpha >= 0 && beta >= 0 && gamma >= 0 && alpha <= 1 && beta <= 1 && gamma <= 1
}

// Add a reproduction event at the specified position
func (r *Renderer) AddReproductionEvent(position types.Point) {
	r.reproductionEvents = append(r.reproductionEvents, ReproductionEvent{
		Position: position,
		TimeLeft: 1.0, // 1 second duration
	})
}

// Update reproduction events (fade out over time)
func (r *Renderer) updateReproductionEvents(deltaTime float64) {
	// If we have too many events, trim the list to prevent memory issues
	if len(r.reproductionEvents) > 100 {
		r.reproductionEvents = r.reproductionEvents[len(r.reproductionEvents)-100:]
	}

	// Update existing events
	updatedEvents := make([]ReproductionEvent, 0, len(r.reproductionEvents))
	for _, event := range r.reproductionEvents {
		event.TimeLeft -= deltaTime
		if event.TimeLeft > 0 {
			updatedEvents = append(updatedEvents, event)
		}
	}
	r.reproductionEvents = updatedEvents

	// Check for new reproduction events by comparing organism count
	currentCount, _ := r.World.GetPopulationInfo()
	if currentCount > r.previousOrgCount {
		// Get the newest organisms for visual effects
		organisms := r.World.GetOrganisms()
		if len(organisms) > 0 {
			// Just add an effect at the newest organism position (the last in the list)
			// In a more sophisticated implementation, we'd track exact reproduction events
			r.AddReproductionEvent(organisms[len(organisms)-1].Position)
		}
	}
	r.previousOrgCount = currentCount
}

// Draw reproduction events as expanding circles
func (r *Renderer) drawReproductionEvents(screen *ebiten.Image) {
	for _, event := range r.reproductionEvents {
		// Convert world coordinates to screen coordinates
		screenX, screenY := r.worldToScreen(event.Position)

		// Calculate radius based on time left (grows then shrinks)
		timeProgress := 1.0 - event.TimeLeft
		radius := 10.0 * math.Sin(timeProgress*math.Pi) // Sine wave for smooth animation

		// Calculate alpha (fades out)
		alpha := uint8(255 * event.TimeLeft)

		// Draw a series of concentric circles with decreasing alpha
		for i := 0; i < 3; i++ {
			innerRadius := radius * float64(i+1) * 0.5
			innerAlpha := alpha / uint8(i+1)

			// Yellow-orange glow for reproduction
			glowColor := color.RGBA{255, 200, 50, innerAlpha}

			// Draw the circle approximately using line segments
			const segments = 12
			for j := 0; j < segments; j++ {
				angle1 := float64(j) * 2 * math.Pi / segments
				angle2 := float64(j+1) * 2 * math.Pi / segments

				x1 := screenX + math.Cos(angle1)*innerRadius
				y1 := screenY + math.Sin(angle1)*innerRadius
				x2 := screenX + math.Cos(angle2)*innerRadius
				y2 := screenY + math.Sin(angle2)*innerRadius

				ebitenutil.DrawLine(screen, x1, y1, x2, y2, glowColor)
			}
		}
	}
}
