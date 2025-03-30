package world

import (
	"math"

	"github.com/zachbeta/evolve_sim/pkg/types"
)

// ConcentrationGrid represents a discrete grid of chemical concentration values
type ConcentrationGrid struct {
	Width     float64     // Width of the world
	Height    float64     // Height of the world
	CellSize  float64     // Size of each grid cell
	NumCellsX int         // Number of cells in X direction
	NumCellsY int         // Number of cells in Y direction
	Grid      [][]float64 // 2D grid of concentration values
}

// NewConcentrationGrid creates a new concentration grid with the specified dimensions and resolution
func NewConcentrationGrid(width, height, cellSize float64) *ConcentrationGrid {
	numCellsX := int(math.Ceil(width / cellSize))
	numCellsY := int(math.Ceil(height / cellSize))

	// Initialize the 2D grid
	grid := make([][]float64, numCellsX)
	for i := range grid {
		grid[i] = make([]float64, numCellsY)
	}

	return &ConcentrationGrid{
		Width:     width,
		Height:    height,
		CellSize:  cellSize,
		NumCellsX: numCellsX,
		NumCellsY: numCellsY,
		Grid:      grid,
	}
}

// SetConcentration sets the concentration value at the specified grid coordinates
func (cg *ConcentrationGrid) SetConcentration(x, y int, value float64) {
	if x >= 0 && x < cg.NumCellsX && y >= 0 && y < cg.NumCellsY {
		cg.Grid[x][y] = value
	}
}

// GetConcentrationAt returns the interpolated concentration value at the specified world coordinates
func (cg *ConcentrationGrid) GetConcentrationAt(point types.Point) float64 {
	// Convert world coordinates to grid coordinates
	gridX := point.X / cg.CellSize
	gridY := point.Y / cg.CellSize

	// Get integer grid coordinates
	x0 := int(math.Floor(gridX))
	y0 := int(math.Floor(gridY))
	x1 := x0 + 1
	y1 := y0 + 1

	// Calculate fractional parts for interpolation
	fx := gridX - float64(x0)
	fy := gridY - float64(y0)

	// Ensure we're within grid bounds
	if x0 < 0 || y0 < 0 || x1 >= cg.NumCellsX || y1 >= cg.NumCellsY {
		// Outside the grid, return zero
		return 0
	}

	// Handle edge case: Use nearest neighbor at the edge of the grid
	if x1 >= cg.NumCellsX {
		x1 = x0
	}
	if y1 >= cg.NumCellsY {
		y1 = y0
	}

	// Bilinear interpolation
	c00 := cg.Grid[x0][y0]
	c10 := cg.Grid[x1][y0]
	c01 := cg.Grid[x0][y1]
	c11 := cg.Grid[x1][y1]

	// Interpolate in x direction
	cx0 := c00*(1-fx) + c10*fx
	cx1 := c01*(1-fx) + c11*fx

	// Interpolate in y direction
	c := cx0*(1-fy) + cx1*fy

	return c
}

// GetGradientAt returns the gradient of the concentration field at the specified world coordinates
// The gradient points in the direction of increasing concentration
func (cg *ConcentrationGrid) GetGradientAt(point types.Point) types.Point {
	// Use central difference method to calculate gradient
	const delta = 0.5 // Small distance for finite difference

	// Get concentrations at points slightly offset from the original
	cCenter := cg.GetConcentrationAt(point)
	cRight := cg.GetConcentrationAt(types.Point{X: point.X + delta, Y: point.Y})
	cLeft := cg.GetConcentrationAt(types.Point{X: point.X - delta, Y: point.Y})
	cUp := cg.GetConcentrationAt(types.Point{X: point.X, Y: point.Y + delta})
	cDown := cg.GetConcentrationAt(types.Point{X: point.X, Y: point.Y - delta})

	// Calculate partial derivatives using central difference
	dCdx := (cRight - cLeft) / (2 * delta)
	dCdy := (cUp - cDown) / (2 * delta)

	// If we're at the edge where central difference isn't available, fall back to forward/backward difference
	if math.Abs(cRight-cCenter) > 1e-9 && math.Abs(cLeft-cCenter) < 1e-9 {
		dCdx = (cRight - cCenter) / delta
	} else if math.Abs(cRight-cCenter) < 1e-9 && math.Abs(cLeft-cCenter) > 1e-9 {
		dCdx = (cCenter - cLeft) / delta
	}

	if math.Abs(cUp-cCenter) > 1e-9 && math.Abs(cDown-cCenter) < 1e-9 {
		dCdy = (cUp - cCenter) / delta
	} else if math.Abs(cUp-cCenter) < 1e-9 && math.Abs(cDown-cCenter) > 1e-9 {
		dCdy = (cCenter - cDown) / delta
	}

	// Return the gradient vector
	gradient := types.Point{X: dCdx, Y: dCdy}

	// Normalize if not zero
	length := math.Sqrt(gradient.X*gradient.X + gradient.Y*gradient.Y)
	if length > 1e-9 {
		gradient.X /= length
		gradient.Y /= length
	}

	return gradient
}

// Note: Contouring functionality (ContourLine, Direction, Cell, Segment types and
// related functions like GenerateContourLines, marchingSquares, and segmentsToContours)
// has been removed to improve performance.
