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

// NewRenderer creates a new renderer with the specified world and config
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

	// Create renderer
	renderer := &Renderer{
		World:               world,
		Simulator:           simulator,
		Config:              config,
		WindowWidth:         config.Render.WindowWidth,
		WindowHeight:        config.Render.WindowHeight,
		ShowGrid:            config.Render.ShowGrid,
		ShowSensors:         config.Render.ShowSensors,
		ShowLegend:          config.Render.ShowLegend,
		ShowTrails:          false, // Default to off
		FPS:                 0.0,
		keyStates:           make(map[ebiten.Key]bool),
		CurrentColorScheme:  colorSchemes[0],
		ColorSchemes:        colorSchemes,
		CurrentSchemeIndex:  0,
		interpolationFactor: 0.5, // Default interpolation for animations
		reproductionEvents:  make([]ReproductionEvent, 0),
		previousOrgCount:    initialCount,
	}

	// Create triangle image for optimized drawing
	renderer.triangleImage = ebiten.NewImage(16, 16)
	renderer.triangleOpts = ebiten.DrawImageOptions{}

	// Register with the simulator to receive reproduction events
	simulator.SetReproductionHandler(renderer.AddReproductionEvent)

	return renderer
}

// isKeyJustPressed checks if a key was just pressed this frame
func (r *Renderer) isKeyJustPressed(key ebiten.Key) bool {
	wasPressed := r.keyStates[key]
	isPressed := ebiten.IsKeyPressed(key)
	r.keyStates[key] = isPressed
	return isPressed && !wasPressed
}

// Update handles user input and updates animation states
func (r *Renderer) Update() error {
	// Process user input first
	// Space: Pause/Resume
	if r.isKeyJustPressed(ebiten.KeySpace) {
		r.Simulator.SetPaused(!r.Simulator.IsPaused)
	}

	// G: Toggle grid
	if r.isKeyJustPressed(ebiten.KeyG) {
		r.ShowGrid = !r.ShowGrid
	}

	// S: Toggle sensor visualization
	if r.isKeyJustPressed(ebiten.KeyS) {
		r.ShowSensors = !r.ShowSensors
	}

	// L: Toggle legend
	if r.isKeyJustPressed(ebiten.KeyL) {
		r.ShowLegend = !r.ShowLegend
	}

	// T: Toggle trails
	if r.isKeyJustPressed(ebiten.KeyT) {
		r.ShowTrails = !r.ShowTrails
	}

	// M: Cycle color schemes
	if r.isKeyJustPressed(ebiten.KeyM) {
		r.CurrentSchemeIndex = (r.CurrentSchemeIndex + 1) % len(r.ColorSchemes)
		r.CurrentColorScheme = r.ColorSchemes[r.CurrentSchemeIndex]
	}

	// R: Reset simulation
	if r.isKeyJustPressed(ebiten.KeyR) {
		r.Simulator.Reset()
	}

	// +: Increase simulation speed
	if r.isKeyJustPressed(ebiten.KeyEqual) {
		r.Simulator.SetSimulationSpeed(r.Simulator.SimulationSpeed * 1.5)
	}

	// -: Decrease simulation speed
	if r.isKeyJustPressed(ebiten.KeyMinus) {
		r.Simulator.SetSimulationSpeed(r.Simulator.SimulationSpeed / 1.5)
	}

	// Step the simulation
	r.Simulator.Step()

	// Update FPS counter
	r.FPS = ebiten.CurrentFPS()

	// Update reproduction events
	r.updateReproductionEvents(r.Simulator.TimeStep * r.Simulator.SimulationSpeed)

	// Update statistics
	stats := simulation.CalculateStatistics(r.World, r.Simulator.Time)
	r.Stats = stats

	return nil
}

// Draw renders the current state of the simulation
func (r *Renderer) Draw(screen *ebiten.Image) {
	// Clear the screen with a dark background
	screen.Fill(color.RGBA{20, 20, 25, 255})

	// Draw concentration grid if available
	r.drawChemicalConcentration(screen)

	// Draw grid for visual reference if enabled
	if r.ShowGrid {
		r.drawGrid(screen)
	}

	// Draw chemical sources
	r.drawChemicalSources(screen)

	// Draw organisms
	r.drawOrganisms(screen)

	// Draw reproduction events
	r.drawReproductionEvents(screen)

	// Draw legend if enabled
	if r.ShowLegend {
		r.drawLegend(screen)
	}

	// Draw statistics
	r.drawStats(screen)
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
	currentTime := r.Simulator.Time // Get current simulation time for animations

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

		// Critical energy effect (pulsing when below 20%)
		var pulseEffect float64 = 1.0
		if energyRatio < 0.2 {
			// Create a pulsing effect based on time
			pulseFrequency := 5.0                                                 // pulses per second
			pulseAmount := 0.5 + 0.5*math.Sin(currentTime*pulseFrequency*math.Pi) // 0.5-1.5 range

			// Make pulse more intense as energy decreases
			pulseIntensity := 1.0 - (energyRatio / 0.2) // 0-1 range as energy drops from 20% to 0%
			pulseEffect = 1.0 + (pulseAmount-1.0)*pulseIntensity

			// Apply pulse to color intensity
			energyRatio = math.Min(1.0, energyRatio*pulseEffect)
		}

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

		// Add pulsing effect for critically low energy
		if energyRatio < 0.2 && pulseEffect > 1.0 {
			sizeMultiplier *= pulseEffect * 0.8 // Pulsing size, slightly subdued
		}

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

		// Draw energy bar
		// Always draw the energy bar, enhanced version
		barWidth := 12.0
		barHeight := 2.5
		barX := screenX - barWidth/2
		barY := screenY - size*2.5 // Position higher above organism

		// Background (empty) bar with border
		bgAlpha := uint8(80 + 120*energyRatio) // More visible when energy is higher
		ebitenutil.DrawRect(screen, barX-0.5, barY-0.5, barWidth+1, barHeight+1, color.RGBA{30, 30, 30, bgAlpha})
		ebitenutil.DrawRect(screen, barX, barY, barWidth, barHeight, color.RGBA{50, 50, 50, bgAlpha})

		// Filled portion based on energy
		fillWidth := barWidth * energyRatio

		// Color changes from red (low) to yellow (medium) to green (high)
		barRed := uint8(255)
		barGreen := uint8(0)

		if energyRatio > 0.5 {
			// Green increases as energy goes from 50% to 100%
			barGreen = uint8(255 * (energyRatio - 0.5) * 2)
		} else {
			// Red stays at max, green increases as energy goes from 0% to 50%
			barGreen = uint8(255 * energyRatio * 2)
		}

		// Make bar pulse for critical energy
		if energyRatio < 0.2 && pulseEffect > 1.0 {
			// Make bar flash more intensely when critically low
			barRed = uint8(math.Min(255, float64(barRed)*pulseEffect))
		}

		// Draw the energy bar with anti-aliasing by drawing multiple rects with varying alpha
		aaOffset := 0.5
		ebitenutil.DrawRect(screen, barX-aaOffset, barY-aaOffset, fillWidth+aaOffset*2, barHeight+aaOffset*2,
			color.RGBA{barRed / 2, barGreen / 2, 0, 128})
		ebitenutil.DrawRect(screen, barX, barY, fillWidth, barHeight,
			color.RGBA{barRed, barGreen, 0, 230})

		// Add glow effect for organisms gaining energy
		// Detect if organism is in optimal environment and gaining energy
		concentration := r.World.GetConcentrationAt(org.Position)
		similarityFactor := 1.0 - math.Min(math.Abs(concentration-org.ChemPreference)/org.ChemPreference, 1.0)

		// If in optimal environment (similarity > 70%), show energy gain glow
		if similarityFactor > 0.7 && energyRatio < 0.99 {
			// Glow intensity based on how optimal the environment is
			glowIntensity := (similarityFactor - 0.7) / 0.3 // 0-1 range

			// Create a pulsing glow effect
			glowFrequency := 2.0
			glowPulse := 0.6 + 0.4*math.Sin(currentTime*glowFrequency*math.Pi*2) // 0.6-1.0 range

			// Glow color matches energy bar but more transparent
			glowRed := barRed / 2
			glowGreen := barGreen / 2
			glowAlpha := uint8(100 * glowIntensity * glowPulse)

			// Create a glow around the energy bar
			ebitenutil.DrawRect(screen, barX-2, barY-2, fillWidth+4, barHeight+4,
				color.RGBA{glowRed, glowGreen, 0, glowAlpha})
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

		// Draw generation number above energy bar if multi-generation simulation is running
		if org.Generation > 1 {
			// Only draw for non-first generation organisms
			genText := fmt.Sprintf("Gen %d", org.Generation)

			// Calculate text position above energy bar
			textX := int(barX)
			textY := int(barY - 10)

			ebitenutil.DebugPrintAt(screen, genText, textX, textY)
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

// drawLegend shows a legend explaining the colors and symbols used in the simulation
func (r *Renderer) drawLegend(screen *ebiten.Image) {
	// Position and size of the legend
	margin := 20
	legendWidth := 220
	lineHeight := 18
	x := r.WindowWidth - legendWidth - margin
	y := margin

	// Background for the legend
	for ly := y - 5; ly < y+200; ly++ {
		for lx := x - 5; lx < x+legendWidth; lx++ {
			if lx >= 0 && lx < r.WindowWidth && ly >= 0 && ly < r.WindowHeight {
				screen.Set(lx, ly, color.RGBA{0, 0, 0, 150})
			}
		}
	}

	// Header
	ebitenutil.DebugPrintAt(screen, "LEGEND", x, y)
	y += lineHeight + 5

	// Organism color explanation
	ebitenutil.DebugPrintAt(screen, "Organisms:", x, y)
	y += lineHeight

	// Preference colors
	prefBoxSize := 10
	ebitenutil.DebugPrintAt(screen, "Preference:", x, y)

	// Low preference (blue)
	for py := y - prefBoxSize + 2; py < y+2; py++ {
		for px := x + 100; px < x+100+prefBoxSize; px++ {
			screen.Set(px, py, color.RGBA{0, 0, 255, 255})
		}
	}
	ebitenutil.DebugPrintAt(screen, "Low", x+115, y)

	// Medium preference (green)
	for py := y - prefBoxSize + 2; py < y+2; py++ {
		for px := x + 150; px < x+150+prefBoxSize; px++ {
			screen.Set(px, py, color.RGBA{0, 255, 0, 255})
		}
	}
	ebitenutil.DebugPrintAt(screen, "Mid", x+165, y)

	// High preference (red)
	for py := y - prefBoxSize + 2; py < y+2; py++ {
		for px := x + 195; px < x+195+prefBoxSize; px++ {
			screen.Set(px, py, color.RGBA{255, 0, 0, 255})
		}
	}
	ebitenutil.DebugPrintAt(screen, "High", x+210, y)

	y += lineHeight

	// Energy level
	ebitenutil.DebugPrintAt(screen, "Energy:", x, y)

	// Full energy
	for py := y - prefBoxSize + 2; py < y+2; py++ {
		for px := x + 100; px < x+100+prefBoxSize; px++ {
			screen.Set(px, py, color.RGBA{200, 200, 200, 255})
		}
	}
	ebitenutil.DebugPrintAt(screen, "Full", x+115, y)

	// Low energy
	for py := y - prefBoxSize + 2; py < y+2; py++ {
		for px := x + 150; px < x+150+prefBoxSize; px++ {
			screen.Set(px, py, color.RGBA{100, 100, 100, 255})
		}
	}
	ebitenutil.DebugPrintAt(screen, "Low", x+165, y)

	// Critical (pulsing)
	for py := y - prefBoxSize + 2; py < y+2; py++ {
		for px := x + 195; px < x+195+prefBoxSize; px++ {
			screen.Set(px, py, color.RGBA{255, 0, 0, 200})
		}
	}
	ebitenutil.DebugPrintAt(screen, "Critical", x+210, y)

	y += lineHeight + 5

	// Chemical sources
	ebitenutil.DebugPrintAt(screen, "Chemical Sources:", x, y)
	y += lineHeight

	// Draw a small representation of a chemical source
	sourceX := float64(x + 30)
	sourceY := float64(y + 5)
	sourceRadius := 8.0

	// Draw the source circle
	for cy := int(sourceY) - int(sourceRadius); cy <= int(sourceY)+int(sourceRadius); cy++ {
		for cx := int(sourceX) - int(sourceRadius); cx <= int(sourceX)+int(sourceRadius); cx++ {
			dx := float64(cx) - sourceX
			dy := float64(cy) - sourceY
			if dx*dx+dy*dy <= sourceRadius*sourceRadius {
				screen.Set(cx, cy, GetColorFromScheme(r.CurrentColorScheme, 0.5))
			}
		}
	}

	// Draw the source outline
	outlineColor := color.RGBA{255, 255, 255, 200}
	for angle := 0.0; angle < 2*math.Pi; angle += 0.1 {
		cx := int(sourceX + math.Cos(angle)*sourceRadius)
		cy := int(sourceY + math.Sin(angle)*sourceRadius)
		screen.Set(cx, cy, outlineColor)
	}

	// Source description
	ebitenutil.DebugPrintAt(screen, "Size indicates energy", x+45, y)
	y += lineHeight
	ebitenutil.DebugPrintAt(screen, "Color indicates decay rate", x+45, y)

	y += lineHeight + 5

	// Reproduction events
	ebitenutil.DebugPrintAt(screen, "Reproduction:", x, y)
	y += lineHeight

	// Draw a small representation of a reproduction event
	reproX := float64(x + 30)
	reproY := float64(y)
	reproRadius := 6.0

	// Draw the ripple effect
	for radius := 3.0; radius <= reproRadius; radius += 1.0 {
		alpha := uint8(255 * (reproRadius - radius) / reproRadius)
		rippleColor := color.RGBA{255, 200, 0, alpha}
		for angle := 0.0; angle < 2*math.Pi; angle += 0.1 {
			cx := int(reproX + math.Cos(angle)*radius)
			cy := int(reproY + math.Sin(angle)*radius)
			screen.Set(cx, cy, rippleColor)
		}
	}

	ebitenutil.DebugPrintAt(screen, "Yellow ripple effect", x+45, y)

	y += lineHeight + 5

	// Controls
	ebitenutil.DebugPrintAt(screen, "CONTROLS", x, y)
	y += lineHeight
	ebitenutil.DebugPrintAt(screen, "Space: Pause/Resume", x, y)
	y += lineHeight
	ebitenutil.DebugPrintAt(screen, "G: Toggle Grid", x, y)
	y += lineHeight
	ebitenutil.DebugPrintAt(screen, "S: Toggle Sensors", x, y)
	y += lineHeight
	ebitenutil.DebugPrintAt(screen, "L: Toggle Legend", x, y)
	y += lineHeight
	ebitenutil.DebugPrintAt(screen, "T: Toggle Trails", x, y)
	y += lineHeight
	ebitenutil.DebugPrintAt(screen, "R: Reset Simulation", x, y)
}
