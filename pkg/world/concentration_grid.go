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

// ContourLine represents a line connecting points of equal concentration
type ContourLine struct {
	Level  float64       // The concentration level
	Points []types.Point // Points along the contour
}

// Direction represents an edge direction in the marching squares algorithm
type Direction int

const (
	Bottom Direction = iota
	Right
	Top
	Left
)

// Cell defines the concentration values at the four corners of a grid cell
type Cell struct {
	X, Y   int
	Values [4]float64 // Concentration at corners (bottom-left, bottom-right, top-right, top-left)
}

// Segment represents a line segment between two points
type Segment struct {
	Start, End types.Point
}

// GenerateContourLines generates contour lines at specified concentration levels
// Returns a map of level -> list of contour lines
func (cg *ConcentrationGrid) GenerateContourLines(levels []float64) map[float64][]ContourLine {
	result := make(map[float64][]ContourLine)

	// Initialize result map with empty slices for each level
	for _, level := range levels {
		result[level] = []ContourLine{}
	}

	// Process each grid cell
	for x := 0; x < cg.NumCellsX-1; x++ {
		for y := 0; y < cg.NumCellsY-1; y++ {
			cell := Cell{
				X: x,
				Y: y,
				Values: [4]float64{
					cg.Grid[x][y],     // Bottom-left
					cg.Grid[x+1][y],   // Bottom-right
					cg.Grid[x+1][y+1], // Top-right
					cg.Grid[x][y+1],   // Top-left
				},
			}

			// Generate contour segments for each level
			for _, level := range levels {
				segments := cg.marchingSquares(cell, level)

				if len(segments) > 0 {
					// Convert segments to contour lines
					contours := cg.segmentsToContours(segments, level)
					result[level] = append(result[level], contours...)
				}
			}
		}
	}

	return result
}

// marchingSquares implements the marching squares algorithm for a single cell and contour level
func (cg *ConcentrationGrid) marchingSquares(cell Cell, level float64) []Segment {
	// Determine which corners are above the threshold
	caseIndex := 0
	for i := 0; i < 4; i++ {
		if cell.Values[i] >= level {
			caseIndex |= (1 << i)
		}
	}

	// Early exit for cases with no contour lines
	if caseIndex == 0 || caseIndex == 15 {
		return nil
	}

	// Calculate edge intersection points using linear interpolation
	edges := make(map[Direction]types.Point)

	// Check bottom edge (between corners 0 and 1)
	if (caseIndex & 1) != ((caseIndex & 2) >> 1) {
		t := (level - cell.Values[0]) / (cell.Values[1] - cell.Values[0])
		edges[Bottom] = types.Point{
			X: (float64(cell.X) + t) * cg.CellSize,
			Y: float64(cell.Y) * cg.CellSize,
		}
	}

	// Check right edge (between corners 1 and 2)
	if ((caseIndex & 2) >> 1) != ((caseIndex & 4) >> 2) {
		t := (level - cell.Values[1]) / (cell.Values[2] - cell.Values[1])
		edges[Right] = types.Point{
			X: float64(cell.X+1) * cg.CellSize,
			Y: (float64(cell.Y) + t) * cg.CellSize,
		}
	}

	// Check top edge (between corners 3 and 2)
	if ((caseIndex & 8) >> 3) != ((caseIndex & 4) >> 2) {
		t := (level - cell.Values[3]) / (cell.Values[2] - cell.Values[3])
		edges[Top] = types.Point{
			X: (float64(cell.X) + 1.0 - t) * cg.CellSize,
			Y: float64(cell.Y+1) * cg.CellSize,
		}
	}

	// Check left edge (between corners 0 and 3)
	if (caseIndex & 1) != ((caseIndex & 8) >> 3) {
		t := (level - cell.Values[0]) / (cell.Values[3] - cell.Values[0])
		edges[Left] = types.Point{
			X: float64(cell.X) * cg.CellSize,
			Y: (float64(cell.Y) + t) * cg.CellSize,
		}
	}

	// Connect the edges to form line segments based on the case
	segments := make([]Segment, 0, 2) // At most 2 segments per cell

	// Handle ambiguous cases (saddle points)
	if caseIndex == 5 || caseIndex == 10 {
		// Calculate the center value as the average of the four corners
		centerValue := (cell.Values[0] + cell.Values[1] + cell.Values[2] + cell.Values[3]) / 4.0

		// Disambiguate based on the center value
		if (caseIndex == 5 && centerValue >= level) || (caseIndex == 10 && centerValue < level) {
			// Connect the edges differently to resolve ambiguity
			if caseIndex == 5 {
				segments = append(segments, Segment{edges[Left], edges[Bottom]})
				segments = append(segments, Segment{edges[Right], edges[Top]})
			} else { // caseIndex == 10
				segments = append(segments, Segment{edges[Bottom], edges[Right]})
				segments = append(segments, Segment{edges[Left], edges[Top]})
			}
			return segments
		}
	}

	// Standard cases with lookup table approach
	switch caseIndex {
	case 1, 14:
		segments = append(segments, Segment{edges[Bottom], edges[Left]})
	case 2, 13:
		segments = append(segments, Segment{edges[Bottom], edges[Right]})
	case 3, 12:
		segments = append(segments, Segment{edges[Left], edges[Right]})
	case 4, 11:
		segments = append(segments, Segment{edges[Right], edges[Top]})
	case 5:
		segments = append(segments, Segment{edges[Bottom], edges[Right]})
		segments = append(segments, Segment{edges[Left], edges[Top]})
	case 6, 9:
		segments = append(segments, Segment{edges[Bottom], edges[Top]})
	case 7, 8:
		segments = append(segments, Segment{edges[Left], edges[Top]})
	case 10:
		segments = append(segments, Segment{edges[Left], edges[Bottom]})
		segments = append(segments, Segment{edges[Right], edges[Top]})
	}

	return segments
}

// segmentsToContours converts line segments to contour lines
func (cg *ConcentrationGrid) segmentsToContours(segments []Segment, level float64) []ContourLine {
	if len(segments) == 0 {
		return nil
	}

	// Create a map to track processed segments
	processed := make(map[int]bool)

	// List to hold resulting contour lines
	contours := make([]ContourLine, 0)

	// Process all segments
	for i := 0; i < len(segments); i++ {
		if processed[i] {
			continue
		}

		// Start a new contour
		contour := ContourLine{
			Level:  level,
			Points: make([]types.Point, 0),
		}

		// Add the first segment
		segment := segments[i]
		contour.Points = append(contour.Points, segment.Start, segment.End)
		processed[i] = true

		// Try to extend the contour by finding connected segments
		// Look for a segment where Start or End matches our End
		endPoint := segment.End

		// Keep extending the contour until no more connected segments are found
		for {
			foundConnection := false

			for j := 0; j < len(segments); j++ {
				if processed[j] {
					continue
				}

				// Check if this segment connects to our end point
				if pointsAreClose(segments[j].Start, endPoint) {
					contour.Points = append(contour.Points, segments[j].End)
					endPoint = segments[j].End
					processed[j] = true
					foundConnection = true
					break
				} else if pointsAreClose(segments[j].End, endPoint) {
					contour.Points = append(contour.Points, segments[j].Start)
					endPoint = segments[j].Start
					processed[j] = true
					foundConnection = true
					break
				}
			}

			if !foundConnection {
				break
			}
		}

		// If the contour has at least 2 points, add it to the result
		if len(contour.Points) >= 2 {
			contours = append(contours, contour)
		}
	}

	return contours
}

// pointsAreClose checks if two points are within a small epsilon of each other
func pointsAreClose(p1, p2 types.Point) bool {
	const epsilon = 1e-6
	dx := p1.X - p2.X
	dy := p1.Y - p2.Y
	return (dx*dx + dy*dy) < epsilon
}
