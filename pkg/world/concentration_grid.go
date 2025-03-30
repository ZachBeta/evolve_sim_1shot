package world

import (
	"math"

	"github.com/zachbeta/evolve_sim/pkg/types"
)

// ConcentrationGrid represents a simplified grid of chemical concentration values
// This is a performance-optimized version that doesn't store full concentration data
type ConcentrationGrid struct {
	Width     float64                // Width of the world
	Height    float64                // Height of the world
	CellSize  float64                // Size of each grid cell
	NumCellsX int                    // Number of cells in X direction
	NumCellsY int                    // Number of cells in Y direction
	Sources   []types.ChemicalSource // References to chemical sources
}

// NewConcentrationGrid creates a new concentration grid with the specified dimensions and resolution
func NewConcentrationGrid(width, height, cellSize float64) *ConcentrationGrid {
	numCellsX := int(math.Ceil(width / cellSize))
	numCellsY := int(math.Ceil(height / cellSize))

	return &ConcentrationGrid{
		Width:     width,
		Height:    height,
		CellSize:  cellSize,
		NumCellsX: numCellsX,
		NumCellsY: numCellsY,
		Sources:   make([]types.ChemicalSource, 0),
	}
}

// SetConcentration is a compatibility function that does nothing in the simplified implementation
func (cg *ConcentrationGrid) SetConcentration(x, y int, value float64) {
	// No-op in simplified implementation
}

// SetSources updates the reference to chemical sources
func (cg *ConcentrationGrid) SetSources(sources []types.ChemicalSource) {
	cg.Sources = make([]types.ChemicalSource, len(sources))
	copy(cg.Sources, sources)
}

// GetConcentrationAt returns the concentration value at the specified world coordinates
// This simplified version calculates directly from sources without using a grid
func (cg *ConcentrationGrid) GetConcentrationAt(point types.Point) float64 {
	// Direct calculation from sources
	var totalConcentration float64 = 0

	// Find nearest source as a simple approximation
	minDist := math.MaxFloat64
	var nearestSource *types.ChemicalSource

	for i := range cg.Sources {
		source := &cg.Sources[i]
		if !source.IsActive {
			continue
		}

		dist := source.Position.DistanceTo(point)
		if dist < minDist {
			minDist = dist
			nearestSource = source
		}
	}

	// If we found a nearby source, return its concentration
	if nearestSource != nil && minDist < cg.Width/5 {
		return nearestSource.GetConcentrationAt(point)
	}

	return totalConcentration
}

// GetGradientAt returns the gradient of the concentration field at the specified world coordinates
// This simplified version calculates direction toward nearest chemical source
func (cg *ConcentrationGrid) GetGradientAt(point types.Point) types.Point {
	// Find the nearest active chemical source
	minDist := math.MaxFloat64
	var nearestSource *types.ChemicalSource

	for i := range cg.Sources {
		source := &cg.Sources[i]
		if !source.IsActive {
			continue
		}

		dist := source.Position.DistanceTo(point)
		if dist < minDist {
			minDist = dist
			nearestSource = source
		}
	}

	// If we found a source, return direction toward it
	if nearestSource != nil {
		// Vector from point to source
		dx := nearestSource.Position.X - point.X
		dy := nearestSource.Position.Y - point.Y

		// Normalize
		length := math.Sqrt(dx*dx + dy*dy)
		if length > 1e-9 {
			return types.Point{X: dx / length, Y: dy / length}
		}
	}

	// No nearby sources or at source position
	return types.Point{X: 0, Y: 0}
}

// Note: Contouring functionality has been removed to improve performance.
