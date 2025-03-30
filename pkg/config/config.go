package config

import (
	"encoding/json"
	"math"
	"os"
)

// WorldConfig holds settings for the simulation world
type WorldConfig struct {
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

// OrganismConfig holds settings for the simulated organisms
type OrganismConfig struct {
	Count          int     `json:"count"`
	Speed          float64 `json:"speed"`
	SensorDistance float64 `json:"sensorDistance"`
	TurnSpeed      float64 `json:"turnSpeed"` // radians per step
	PrefMean       float64 `json:"preferenceDistributionMean"`
	PrefStdDev     float64 `json:"preferenceDistributionStdDev"`
}

// ChemicalConfig holds settings for chemical sources
type ChemicalConfig struct {
	Count          int     `json:"count"`
	MinStrength    float64 `json:"minStrength"`
	MaxStrength    float64 `json:"maxStrength"`
	MinDecayFactor float64 `json:"minDecayFactor"`
	MaxDecayFactor float64 `json:"maxDecayFactor"`
}

// RenderConfig holds settings for visualization
type RenderConfig struct {
	WindowWidth       int  `json:"windowWidth"`
	WindowHeight      int  `json:"windowHeight"`
	FrameRate         int  `json:"frameRate"`
	ShowGrid          bool `json:"showGrid"`
	ShowConcentration bool `json:"showConcentration"`
	ShowSensors       bool `json:"showSensors"`
}

// SimulationConfig holds all configuration for the simulation
type SimulationConfig struct {
	World           WorldConfig    `json:"world"`
	Organism        OrganismConfig `json:"organism"`
	Chemical        ChemicalConfig `json:"chemical"`
	Render          RenderConfig   `json:"render"`
	RandomSeed      int64          `json:"randomSeed"`
	SimulationSpeed float64        `json:"simulationSpeed"`
}

// DefaultConfig returns a default configuration with reasonable values
func DefaultConfig() SimulationConfig {
	return SimulationConfig{
		World: WorldConfig{
			Width:  1000.0,
			Height: 1000.0,
		},
		Organism: OrganismConfig{
			Count:          100,
			Speed:          2.0,
			SensorDistance: 10.0,
			TurnSpeed:      math.Pi / 10, // 18 degrees per step
			PrefMean:       50.0,
			PrefStdDev:     10.0,
		},
		Chemical: ChemicalConfig{
			Count:          5,
			MinStrength:    100.0,
			MaxStrength:    500.0,
			MinDecayFactor: 0.001,
			MaxDecayFactor: 0.01,
		},
		Render: RenderConfig{
			WindowWidth:       800,
			WindowHeight:      800,
			FrameRate:         60,
			ShowGrid:          true,
			ShowConcentration: true,
			ShowSensors:       true,
		},
		RandomSeed:      0, // 0 means use current time as seed
		SimulationSpeed: 1.0,
	}
}

// LoadConfig loads configuration from a JSON file
func LoadConfig(filename string) (SimulationConfig, error) {
	// Start with default config
	config := DefaultConfig()

	// Read the file
	data, err := os.ReadFile(filename)
	if err != nil {
		return config, err
	}

	// Parse JSON
	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

// SaveConfig saves configuration to a JSON file
func SaveConfig(config SimulationConfig, filename string) error {
	// Convert to JSON with indentation for readability
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	// Write to file
	return os.WriteFile(filename, data, 0644)
}
