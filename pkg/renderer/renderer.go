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

// Renderer is responsible for visualizing the simulation
type Renderer struct {
	World               *world.World
	Simulator           *simulation.Simulator
	Config              config.SimulationConfig
	WindowWidth         int
	WindowHeight        int
	ShowGrid            bool
	ShowConcentration   bool
	ShowSensors         bool
	ShowLegend          bool
	ShowContours        bool
	Stats               simulation.SimulationStats
	FPS                 float64
	keyStates           map[ebiten.Key]bool
	CurrentColorScheme  ColorScheme
	ColorSchemes        []ColorScheme
	CurrentSchemeIndex  int
	ContourLevels       []float64
	contourCache        map[float64][]ContourLine
	lastContourUpdate   float64
	interpolationFactor float64 // For smooth animations between frames
	triangleImage       *ebiten.Image
	triangleOpts        ebiten.DrawImageOptions
	selectedOrganism    *types.Organism // For future organism selection feature
}

// ContourLine is a local representation of the world's ContourLine
type ContourLine struct {
	Level  float64
	Points []types.Point
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

	// Define default contour levels
	contourLevels := []float64{10, 20, 50, 100, 200, 500}

	// Create a new renderer
	r := &Renderer{
		World:               world,
		Simulator:           simulator,
		Config:              config,
		WindowWidth:         config.Render.WindowWidth,
		WindowHeight:        config.Render.WindowHeight,
		ShowGrid:            config.Render.ShowGrid,
		ShowConcentration:   config.Render.ShowConcentration,
		ShowSensors:         config.Render.ShowSensors,
		ShowLegend:          config.Render.ShowLegend,
		ShowContours:        config.Render.ShowContours,
		FPS:                 0,
		keyStates:           make(map[ebiten.Key]bool),
		CurrentColorScheme:  ViridisScheme,
		ColorSchemes:        colorSchemes,
		CurrentSchemeIndex:  0,
		ContourLevels:       contourLevels,
		contourCache:        make(map[float64][]ContourLine),
		lastContourUpdate:   0,
		interpolationFactor: 1.0, // Default to full interpolation
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

	// Update contour lines every 0.5 seconds if showing contours
	if r.ShowContours && r.Simulator.Time-r.lastContourUpdate > 0.5 {
		r.updateContourLines()
		r.lastContourUpdate = r.Simulator.Time
	}

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

	if r.isKeyJustPressed(ebiten.KeyC) {
		r.ShowConcentration = !r.ShowConcentration
	}

	if r.isKeyJustPressed(ebiten.KeyS) {
		r.ShowSensors = !r.ShowSensors
	}

	if r.isKeyJustPressed(ebiten.KeyL) {
		r.ShowLegend = !r.ShowLegend
	}

	if r.isKeyJustPressed(ebiten.KeyO) {
		r.ShowContours = !r.ShowContours
		if r.ShowContours {
			r.updateContourLines()
		}
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
	// Force contour update on first draw if enabled
	if r.ShowContours && len(r.contourCache) == 0 {
		r.updateContourLines()
	}

	// Clear screen
	screen.Fill(color.RGBA{20, 20, 30, 255})

	// Draw grid if enabled
	if r.ShowGrid {
		r.drawGrid(screen)
	}

	// Draw chemical concentration if enabled
	if r.ShowConcentration {
		r.drawChemicalConcentration(screen)
	}

	// Draw contour lines if enabled
	if r.ShowContours {
		r.drawContourLines(screen)
	}

	// Draw chemical sources
	r.drawChemicalSources(screen)

	// Draw organisms
	r.drawOrganisms(screen)

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

// Draw a visualization of chemical concentration
func (r *Renderer) drawChemicalConcentration(screen *ebiten.Image) {
	bounds := r.World.GetBounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	// Define grid resolution for visualization (lower = higher performance)
	cellSizeX := float64(r.WindowWidth) / 80 // Increased resolution
	cellSizeY := float64(r.WindowHeight) / 80

	// Get concentration stats for color scaling
	maxConcentration := r.Stats.Chemicals.MaxConcentration
	if maxConcentration <= 0 {
		maxConcentration = 1.0 // Prevent division by zero
	}

	// Draw concentration grid
	for screenY := 0; screenY < r.WindowHeight; screenY += int(cellSizeY) {
		for screenX := 0; screenX < r.WindowWidth; screenX += int(cellSizeX) {
			// Convert screen coordinates to world coordinates
			normalizedX := float64(screenX) / float64(r.WindowWidth)
			normalizedY := float64(screenY) / float64(r.WindowHeight)

			worldX := bounds.Min.X + normalizedX*width
			worldY := bounds.Min.Y + normalizedY*height

			// Get concentration at this point
			point := types.Point{X: worldX, Y: worldY}
			concentration := r.World.GetConcentrationAt(point)

			// Normalize concentration for color mapping (0.0 to 1.0)
			normalizedConc := math.Min(1.0, concentration/maxConcentration)

			// Get color from current scheme
			cellColor := GetColorFromScheme(r.CurrentColorScheme, normalizedConc)

			// Apply transparency
			cellColor.A = 130 // Semi-transparent

			// Apply a small amount of smoothing to reduce banding
			if screenX > 0 && screenY > 0 && screenX < r.WindowWidth-int(cellSizeX) && screenY < r.WindowHeight-int(cellSizeY) {
				// Simple box blur-like effect by slightly blending with neighbors
				// This is a simplified version - a full implementation would be more complex
				smoothingFactor := 0.2
				if normalizedConc > 0 && normalizedConc < 1 {
					// Only apply smoothing to non-extreme values
					jitter := (math.Sin(float64(screenX)*0.1) + math.Cos(float64(screenY)*0.1)) * smoothingFactor
					normalizedConc += jitter * 0.01 // Subtle dithering effect
				}
			}

			// Draw a rectangle for this grid cell
			for y := 0; y < int(cellSizeY); y++ {
				for x := 0; x < int(cellSizeX); x++ {
					if screenX+x < r.WindowWidth && screenY+y < r.WindowHeight {
						screen.Set(screenX+x, screenY+y, cellColor)
					}
				}
			}
		}
	}

	// Draw legend if enabled
	if r.ShowLegend {
		r.drawConcentrationLegend(screen)
	}
}

// Draw a legend for the concentration visualization
func (r *Renderer) drawConcentrationLegend(screen *ebiten.Image) {
	// Position in bottom-right corner
	legendWidth := 150
	legendHeight := 20
	padding := 10
	x := r.WindowWidth - legendWidth - padding
	y := r.WindowHeight - legendHeight - padding

	// Draw legend background
	for ly := 0; ly < legendHeight; ly++ {
		for lx := 0; lx < legendWidth; lx++ {
			position := float64(lx) / float64(legendWidth)
			color := GetColorFromScheme(r.CurrentColorScheme, position)
			color.A = 200 // More opaque for the legend
			screen.Set(x+lx, y+ly, color)
		}
	}

	// Draw border
	for lx := 0; lx < legendWidth; lx++ {
		screen.Set(x+lx, y-1, color.RGBA{200, 200, 200, 255})
		screen.Set(x+lx, y+legendHeight, color.RGBA{200, 200, 200, 255})
	}
	for ly := -1; ly <= legendHeight; ly++ {
		screen.Set(x-1, y+ly, color.RGBA{200, 200, 200, 255})
		screen.Set(x+legendWidth, y+ly, color.RGBA{200, 200, 200, 255})
	}

	// Draw min/max labels
	minLabel := "0.0"
	maxLabel := fmt.Sprintf("%.1f", r.Stats.Chemicals.MaxConcentration)

	ebitenutil.DebugPrintAt(screen, minLabel, x, y-15)
	ebitenutil.DebugPrintAt(screen, maxLabel, x+legendWidth-30, y-15)

	// Draw scheme name
	schemeName := r.CurrentColorScheme.Name
	ebitenutil.DebugPrintAt(screen, schemeName, x+legendWidth/2-20, y+legendHeight+5)
}

// Draw chemical sources
func (r *Renderer) drawChemicalSources(screen *ebiten.Image) {
	sources := r.World.GetChemicalSources()

	for _, source := range sources {
		// Convert world coordinates to screen coordinates
		screenX, screenY := r.worldToScreen(source.Position)

		// Draw a circle at the source position
		radius := 5.0 + 10.0*(source.Strength/r.Config.Chemical.MaxStrength)

		// Draw a filled circle
		for y := int(screenY) - int(radius); y <= int(screenY)+int(radius); y++ {
			for x := int(screenX) - int(radius); x <= int(screenX)+int(radius); x++ {
				dx := float64(x) - screenX
				dy := float64(y) - screenY
				if dx*dx+dy*dy <= radius*radius {
					screen.Set(x, y, color.RGBA{200, 100, 0, 255})
				}
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

		// Determine color based on chemical preference
		// Map preference to a blue-to-red gradient
		prefRange := r.Config.Organism.PreferenceDistributionMean * 3
		normalizedPref := org.ChemPreference / prefRange

		red := uint8(normalizedPref * 255)
		blue := uint8((1 - normalizedPref) * 255)
		green := uint8(128 - math.Abs(float64(normalizedPref*255-128)))

		// Calculate the visual heading with interpolation for smooth rotation
		visualHeading := org.PreviousHeading + (org.Heading-org.PreviousHeading)*r.interpolationFactor

		// Define triangle size (can be adjusted based on organism properties)
		size := 4.0

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
			color.RGBA{red, green, blue, 255})

		// Add a border for better visibility
		ebitenutil.DrawLine(screen, frontX, frontY, leftX, leftY, color.RGBA{255, 255, 255, 200})
		ebitenutil.DrawLine(screen, leftX, leftY, rightX, rightY, color.RGBA{255, 255, 255, 200})
		ebitenutil.DrawLine(screen, rightX, rightY, frontX, frontY, color.RGBA{255, 255, 255, 200})

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
		fmt.Sprintf("Color Scheme: %s", r.CurrentColorScheme.Name),
		fmt.Sprintf("Grid: %v", r.ShowGrid),
		fmt.Sprintf("Contours: %v", r.ShowContours),
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
		"C: Toggle Concentration",
		"O: Toggle Contours",
		"S: Toggle Sensors",
		"L: Toggle Legend",
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

// Update contour lines
func (r *Renderer) updateContourLines() {
	// Get the concentration grid from the world
	grid := r.World.GetConcentrationGrid()
	if grid == nil {
		fmt.Println("Warning: concentration grid is nil")
		return
	}

	// Generate contours for the current levels
	worldContours := grid.GenerateContourLines(r.ContourLevels)

	// Count total contours for debugging
	totalContours := 0
	for _, contours := range worldContours {
		totalContours += len(contours)
	}

	fmt.Printf("Generated %d contour lines across %d levels\n", totalContours, len(worldContours))

	// Convert to local representation
	r.contourCache = make(map[float64][]ContourLine)
	for level, contours := range worldContours {
		r.contourCache[level] = make([]ContourLine, len(contours))
		for i, contour := range contours {
			r.contourCache[level][i] = ContourLine{
				Level:  contour.Level,
				Points: contour.Points,
			}
		}
	}
}

// Draw contour lines
func (r *Renderer) drawContourLines(screen *ebiten.Image) {
	if len(r.contourCache) == 0 {
		return
	}

	// Iterate through each contour level
	for level, contours := range r.contourCache {
		// Normalize level for color selection
		maxConcentration := r.Stats.Chemicals.MaxConcentration
		if maxConcentration <= 0 {
			maxConcentration = 1.0
		}

		normalizedLevel := math.Min(1.0, level/maxConcentration)

		// Get color for this contour level
		levelColor := GetColorFromScheme(r.CurrentColorScheme, normalizedLevel)
		// Make lines more visible
		levelColor.A = 255 // Fully opaque

		// Draw each contour line
		for _, contour := range contours {
			// Skip contours with too few points
			if len(contour.Points) < 2 {
				continue
			}

			// Draw the contour as connected line segments with increased thickness
			for i := 0; i < len(contour.Points)-1; i++ {
				// Convert world coordinates to screen coordinates
				x1, y1 := r.worldToScreen(contour.Points[i])
				x2, y2 := r.worldToScreen(contour.Points[i+1])

				// Draw thicker line by drawing multiple lines with slight offsets
				ebitenutil.DrawLine(screen, x1, y1, x2, y2, levelColor)

				// Draw additional lines for thickness
				offset := 0.5
				ebitenutil.DrawLine(screen, x1+offset, y1, x2+offset, y2, levelColor)
				ebitenutil.DrawLine(screen, x1-offset, y1, x2-offset, y2, levelColor)
				ebitenutil.DrawLine(screen, x1, y1+offset, x2, y2+offset, levelColor)
				ebitenutil.DrawLine(screen, x1, y1-offset, x2, y2-offset, levelColor)
			}

			// Optionally, draw the contour level value at the middle of the contour
			if len(contour.Points) > 5 && math.Mod(level, 50) < 0.1 { // Only label major contours
				midIndex := len(contour.Points) / 2
				midX, midY := r.worldToScreen(contour.Points[midIndex])

				// Draw label background for better visibility
				labelStr := fmt.Sprintf("%.0f", level)

				// Draw a small background box
				textWidth := 7 * len(labelStr) // Estimate text width
				boxPadding := 2
				for y := int(midY) - 8 - boxPadding; y <= int(midY)+boxPadding; y++ {
					for x := int(midX) - (textWidth / 2) - boxPadding; x <= int(midX)+(textWidth/2)+boxPadding; x++ {
						if x >= 0 && x < r.WindowWidth && y >= 0 && y < r.WindowHeight {
							screen.Set(x, y, color.RGBA{20, 20, 30, 255}) // Fully opaque
						}
					}
				}

				// Draw the text
				ebitenutil.DebugPrintAt(screen, labelStr, int(midX)-textWidth/2, int(midY)-8)
			}
		}
	}
}

// drawTriangle draws a filled triangle with the specified points and color
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
