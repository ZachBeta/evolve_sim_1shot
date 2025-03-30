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
	World             *world.World
	Simulator         *simulation.Simulator
	Config            config.SimulationConfig
	WindowWidth       int
	WindowHeight      int
	ShowGrid          bool
	ShowConcentration bool
	ShowSensors       bool
	Stats             simulation.SimulationStats
	FPS               float64
}

// NewRenderer creates a new renderer with the given configuration
func NewRenderer(world *world.World, simulator *simulation.Simulator, config config.SimulationConfig) *Renderer {
	return &Renderer{
		World:             world,
		Simulator:         simulator,
		Config:            config,
		WindowWidth:       config.Render.WindowWidth,
		WindowHeight:      config.Render.WindowHeight,
		ShowGrid:          config.Render.ShowGrid,
		ShowConcentration: config.Render.ShowConcentration,
		ShowSensors:       config.Render.ShowSensors,
		FPS:               0,
	}
}

// Update is called each frame by Ebiten
func (r *Renderer) Update() error {
	// Update simulator
	r.Simulator.Step()

	// Update stats
	r.Stats = r.Simulator.CollectStats()

	// Update FPS
	r.FPS = ebiten.ActualFPS()

	// Handle keyboard input
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		r.Simulator.SetPaused(!r.Simulator.IsPaused)
	}

	if ebiten.IsKeyPressed(ebiten.KeyR) {
		r.Simulator.Reset()
	}

	if ebiten.IsKeyPressed(ebiten.KeyG) {
		r.ShowGrid = !r.ShowGrid
	}

	if ebiten.IsKeyPressed(ebiten.KeyC) {
		r.ShowConcentration = !r.ShowConcentration
	}

	if ebiten.IsKeyPressed(ebiten.KeyS) {
		r.ShowSensors = !r.ShowSensors
	}

	// Speed control
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

	// Draw chemical concentration if enabled
	if r.ShowConcentration {
		r.drawChemicalConcentration(screen)
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

// Draw a placeholder for chemical concentration
func (r *Renderer) drawChemicalConcentration(screen *ebiten.Image) {
	// This is a placeholder - will be implemented in future
	ebitenutil.DebugPrint(screen, "Chemical concentration visualization not yet implemented")
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

		// Draw a small circle for the organism
		radius := 3.0

		// Draw a filled circle
		for y := int(screenY) - int(radius); y <= int(screenY)+int(radius); y++ {
			for x := int(screenX) - int(radius); x <= int(screenX)+int(radius); x++ {
				dx := float64(x) - screenX
				dy := float64(y) - screenY
				if dx*dx+dy*dy <= radius*radius {
					screen.Set(x, y, color.RGBA{red, green, blue, 255})
				}
			}
		}

		// Draw heading indicator
		headingX := screenX + math.Cos(org.Heading)*8
		headingY := screenY + math.Sin(org.Heading)*8
		ebitenutil.DrawLine(screen, screenX, screenY, headingX, headingY, color.RGBA{255, 255, 255, 200})

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
		"S: Toggle Sensors",
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
