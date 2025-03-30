package simulation

import (
	"os"
	"testing"

	"github.com/zachbeta/evolve_sim/pkg/types"
)

// mockWorld implements a simple world that returns predefined concentrations
type mockWorld struct {
	// Define a function that returns concentration at given point
	concentrationFn func(types.Point) float64
}

func (m mockWorld) GetConcentrationAt(p types.Point) float64 {
	return m.concentrationFn(p)
}

// TestCalculateOrganismStats tests the organism statistics calculation
func TestCalculateOrganismStats(t *testing.T) {
	// Create mock world
	mockWorld := mockWorld{
		concentrationFn: func(p types.Point) float64 {
			// Simple linear concentration: equals X coordinate
			return p.X
		},
	}

	// Create test organisms
	organisms := []types.Organism{
		// Organism at low concentration with low preference
		types.NewOrganism(
			types.Point{X: 10, Y: 50},
			0,
			15.0, // preference
			1.0,
			types.DefaultSensorAngles(),
		),
		// Organism at medium concentration with medium preference
		types.NewOrganism(
			types.Point{X: 50, Y: 50},
			0,
			50.0, // preference
			1.0,
			types.DefaultSensorAngles(),
		),
		// Organism at high concentration with high preference
		types.NewOrganism(
			types.Point{X: 90, Y: 50},
			0,
			85.0, // preference
			1.0,
			types.DefaultSensorAngles(),
		),
	}

	// Calculate stats
	stats := calculateOrganismStats(organisms, mockWorld)

	// Verify statistics
	if stats.Count != 3 {
		t.Errorf("Expected 3 organisms, got %d", stats.Count)
	}

	// Check min/max
	if stats.MinPreference != 15.0 {
		t.Errorf("Expected min preference 15.0, got %f", stats.MinPreference)
	}
	if stats.MaxPreference != 85.0 {
		t.Errorf("Expected max preference 85.0, got %f", stats.MaxPreference)
	}

	// Check average (should be (15+50+85)/3 = 50)
	expectedAvg := 50.0
	if stats.AveragePreference < expectedAvg-0.1 || stats.AveragePreference > expectedAvg+0.1 {
		t.Errorf("Expected average preference around %f, got %f", expectedAvg, stats.AveragePreference)
	}

	// Check histogram buckets existence
	buckets := []string{"15", "50", "85"}
	for _, bucket := range buckets {
		if stats.PreferenceHistogram[bucket] != 1 {
			t.Errorf("Expected bucket %s to have count 1, got %d", bucket, stats.PreferenceHistogram[bucket])
		}
	}

	// Test with empty organisms list
	emptyStats := calculateOrganismStats([]types.Organism{}, mockWorld)
	if emptyStats.Count != 0 {
		t.Errorf("Expected 0 organisms, got %d", emptyStats.Count)
	}
}

// TestCalculateChemicalStats tests the chemical statistics calculation
func TestCalculateChemicalStats(t *testing.T) {
	// Create a bounded test area
	bounds := types.Rect{
		Min: types.Point{X: 0, Y: 0},
		Max: types.Point{X: 100, Y: 100},
	}

	// Create mock world with a single chemical source
	sources := []types.ChemicalSource{
		{
			Position:    types.Point{X: 50, Y: 50},
			Strength:    100.0,
			DecayFactor: 0.01,
		},
	}

	// Create a world with inverse square law concentration
	mockWorld := mockWorld{
		concentrationFn: func(p types.Point) float64 {
			// Simple distance-based concentration from the center
			dx := p.X - 50
			dy := p.Y - 50
			distSq := dx*dx + dy*dy
			if distSq < 1 {
				distSq = 1 // Avoid division by zero
			}
			return 100.0 / distSq
		},
	}

	// Calculate stats
	stats := calculateChemicalStats(sources, mockWorld, bounds)

	// Verify source count
	if stats.SourceCount != 1 {
		t.Errorf("Expected 1 source, got %d", stats.SourceCount)
	}

	// Center should have highest concentration
	if stats.MaxConcentration < 5.0 {
		t.Errorf("Expected max concentration > 5.0, got %f", stats.MaxConcentration)
	}

	// Corners should have lowest concentration
	if stats.MinConcentration > 1.0 {
		t.Errorf("Expected min concentration < 1.0, got %f", stats.MinConcentration)
	}

	// Verify histogram has entries
	if len(stats.ConcentrationHistogram) == 0 {
		t.Errorf("Expected non-empty concentration histogram")
	}
}

// TestExportStatsCSV tests CSV export functionality
func TestExportStatsCSV(t *testing.T) {
	// Create test statistics
	stats := []SimulationStats{
		{
			Time: 0.0,
			Organisms: OrganismStats{
				Count:             10,
				AveragePreference: 25.0,
				PreferenceStdDev:  5.0,
			},
			Chemicals: ChemicalStats{
				SourceCount:          3,
				AverageConcentration: 30.0,
				MaxConcentration:     100.0,
			},
		},
		{
			Time: 10.0,
			Organisms: OrganismStats{
				Count:             10,
				AveragePreference: 25.5,
				PreferenceStdDev:  4.8,
			},
			Chemicals: ChemicalStats{
				SourceCount:          3,
				AverageConcentration: 30.0,
				MaxConcentration:     100.0,
			},
		},
	}

	// Create temporary file
	tempFile, err := os.CreateTemp("", "stats_test_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	tempFile.Close()
	defer os.Remove(tempFile.Name())

	// Export stats
	err = ExportStatsCSV(stats, tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to export stats: %v", err)
	}

	// Check if file exists and has content
	info, err := os.Stat(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to stat output file: %v", err)
	}
	if info.Size() == 0 {
		t.Errorf("Expected non-empty CSV file")
	}
}

// TestExportStatsJSON tests JSON export functionality
func TestExportStatsJSON(t *testing.T) {
	// Create test statistics
	stats := []SimulationStats{
		{
			Time: 0.0,
			Organisms: OrganismStats{
				Count:             10,
				AveragePreference: 25.0,
				PreferenceStdDev:  5.0,
				PreferenceHistogram: map[string]int{
					"20": 5,
					"30": 5,
				},
			},
			Chemicals: ChemicalStats{
				SourceCount:          3,
				AverageConcentration: 30.0,
				MaxConcentration:     100.0,
			},
		},
	}

	// Create temporary file
	tempFile, err := os.CreateTemp("", "stats_test_*.json")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	tempFile.Close()
	defer os.Remove(tempFile.Name())

	// Export stats
	err = ExportStatsJSON(stats, tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to export stats: %v", err)
	}

	// Check if file exists and has content
	info, err := os.Stat(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to stat output file: %v", err)
	}
	if info.Size() == 0 {
		t.Errorf("Expected non-empty JSON file")
	}
}
