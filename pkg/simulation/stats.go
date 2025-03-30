package simulation

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/zachbeta/evolve_sim/pkg/types"
)

// OrganismStats holds statistics about organisms in the simulation
type OrganismStats struct {
	Count                   int
	AveragePreference       float64
	PreferenceStdDev        float64
	MinPreference           float64
	MaxPreference           float64
	AverageConcentration    float64
	PreferenceHistogram     map[string]int // Bucketized preferences
	PreferenceExposureRatio float64        // Average ratio of preference to actual concentration
	AverageEnergy           float64        // Average energy level of organisms
	EnergyRatio             float64        // Average energy as percentage of capacity
}

// ChemicalStats holds statistics about chemical concentrations
type ChemicalStats struct {
	SourceCount            int
	AverageConcentration   float64
	MaxConcentration       float64
	MinConcentration       float64
	ConcentrationHistogram map[string]int // Bucketized concentrations
}

// SimulationStats holds all statistics for a simulation
type SimulationStats struct {
	Time            float64
	RealTimeElapsed time.Duration
	Organisms       OrganismStats
	Chemicals       ChemicalStats
}

// Histogram bucket size
const histogramBucketSize = 5.0

// calculateOrganismStats calculates statistics about organisms
func calculateOrganismStats(organisms []types.Organism, world interface{ GetConcentrationAt(types.Point) float64 }) OrganismStats {
	if len(organisms) == 0 {
		return OrganismStats{
			Count:               0,
			PreferenceHistogram: make(map[string]int),
		}
	}

	// Initialize stats
	stats := OrganismStats{
		Count:               len(organisms),
		MinPreference:       math.MaxFloat64,
		MaxPreference:       -math.MaxFloat64,
		PreferenceHistogram: make(map[string]int),
	}

	// Sum for average calculation
	var preferenceSum float64
	var concentrationSum float64
	var preferenceDiffSum float64
	var exposureRatioSum float64
	var energySum float64
	var energyRatioSum float64
	preferences := make([]float64, len(organisms))

	// Collect data
	for i, org := range organisms {
		pref := org.ChemPreference
		preferences[i] = pref
		preferenceSum += pref

		// Update min/max
		if pref < stats.MinPreference {
			stats.MinPreference = pref
		}
		if pref > stats.MaxPreference {
			stats.MaxPreference = pref
		}

		// Build histogram
		bucket := fmt.Sprintf("%.0f", math.Floor(pref/histogramBucketSize)*histogramBucketSize)
		stats.PreferenceHistogram[bucket]++

		// Get actual concentration at organism position
		conc := world.GetConcentrationAt(org.Position)
		concentrationSum += conc

		// Calculate preference exposure ratio (how close organism is to its preferred concentration)
		// Avoid division by zero
		if conc > 0 {
			ratio := pref / conc
			if ratio > 1 {
				ratio = 1 / ratio // Normalize to 0-1 range
			}
			exposureRatioSum += ratio
		}

		// Add energy statistics
		energySum += org.Energy
		energyRatioSum += org.Energy / org.EnergyCapacity
	}

	// Calculate averages
	stats.AveragePreference = preferenceSum / float64(len(organisms))
	stats.AverageConcentration = concentrationSum / float64(len(organisms))
	stats.PreferenceExposureRatio = exposureRatioSum / float64(len(organisms))
	stats.AverageEnergy = energySum / float64(len(organisms))
	stats.EnergyRatio = energyRatioSum / float64(len(organisms))

	// Calculate standard deviation
	for _, pref := range preferences {
		diff := pref - stats.AveragePreference
		preferenceDiffSum += diff * diff
	}
	stats.PreferenceStdDev = math.Sqrt(preferenceDiffSum / float64(len(organisms)))

	return stats
}

// calculateChemicalStats calculates statistics about chemical concentrations
func calculateChemicalStats(sources []types.ChemicalSource, world interface{ GetConcentrationAt(types.Point) float64 }, bounds types.Rect) ChemicalStats {
	stats := ChemicalStats{
		SourceCount:            len(sources),
		MinConcentration:       math.MaxFloat64,
		MaxConcentration:       -math.MaxFloat64,
		ConcentrationHistogram: make(map[string]int),
	}

	// Simple sampling grid for concentration statistics
	const samplesX = 20
	const samplesY = 20
	var concentrationSum float64
	var samples int

	// Sample concentrations
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	for x := 0; x < samplesX; x++ {
		for y := 0; y < samplesY; y++ {
			// Calculate sample position
			point := types.Point{
				X: bounds.Min.X + width*float64(x)/float64(samplesX-1),
				Y: bounds.Min.Y + height*float64(y)/float64(samplesY-1),
			}

			// Get concentration
			conc := world.GetConcentrationAt(point)

			// Update stats
			concentrationSum += conc
			samples++

			// Update min/max
			if conc < stats.MinConcentration {
				stats.MinConcentration = conc
			}
			if conc > stats.MaxConcentration {
				stats.MaxConcentration = conc
			}

			// Add to histogram
			bucket := fmt.Sprintf("%.0f", math.Floor(conc/histogramBucketSize)*histogramBucketSize)
			stats.ConcentrationHistogram[bucket]++
		}
	}

	// Calculate average
	if samples > 0 {
		stats.AverageConcentration = concentrationSum / float64(samples)
	}

	return stats
}

// CollectStats collects statistics for the current simulation state
func (s *Simulator) CollectStats() SimulationStats {
	return SimulationStats{
		Time:            s.Time,
		RealTimeElapsed: time.Duration(0), // Will be set by caller if needed
		Organisms:       calculateOrganismStats(s.World.GetOrganisms(), s.World),
		Chemicals:       calculateChemicalStats(s.World.GetChemicalSources(), s.World, s.World.GetBounds()),
	}
}

// ExportStatsCSV exports a time series of simulation statistics to a CSV file
func ExportStatsCSV(stats []SimulationStats, filename string) error {
	// Create file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{
		"Time",
		"OrganismCount",
		"AveragePreference",
		"PreferenceStdDev",
		"AverageConcentration",
		"PreferenceExposureRatio",
		"MaxConcentration",
	}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Write data rows
	for _, stat := range stats {
		row := []string{
			fmt.Sprintf("%.2f", stat.Time),
			fmt.Sprintf("%d", stat.Organisms.Count),
			fmt.Sprintf("%.2f", stat.Organisms.AveragePreference),
			fmt.Sprintf("%.2f", stat.Organisms.PreferenceStdDev),
			fmt.Sprintf("%.2f", stat.Organisms.AverageConcentration),
			fmt.Sprintf("%.2f", stat.Organisms.PreferenceExposureRatio),
			fmt.Sprintf("%.2f", stat.Chemicals.MaxConcentration),
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

// ExportStatsJSON exports simulation statistics to a JSON file
func ExportStatsJSON(stats []SimulationStats, filename string) error {
	// Marshal data to JSON
	data, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		return err
	}

	// Write to file
	return os.WriteFile(filename, data, 0644)
}
