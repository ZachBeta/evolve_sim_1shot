package config

import (
	"encoding/json"
	"math"
	"os"
)

// Version is the current application version
const Version = "0.1.0"

// WorldConfig holds settings for the simulation world
type WorldConfig struct {
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

// OrganismConfig holds settings for the simulated organisms
type OrganismConfig struct {
	Count                        int     `json:"count"`
	Speed                        float64 `json:"speed"`
	SensorDistance               float64 `json:"sensorDistance"`
	TurnSpeed                    float64 `json:"turnSpeed"` // radians per step
	PreferenceDistributionMean   float64 `json:"preferenceDistributionMean"`
	PreferenceDistributionStdDev float64 `json:"preferenceDistributionStdDev"`
}

// EnergyConfig holds settings for the energy system
type EnergyConfig struct {
	InitialEnergy         float64    `json:"initialEnergy"`         // Starting energy for new organisms
	MaximumEnergy         float64    `json:"maximumEnergy"`         // Maximum energy capacity
	BaseMetabolicRate     float64    `json:"baseMetabolicRate"`     // Energy consumed per second just existing
	MovementCostFactor    float64    `json:"movementCostFactor"`    // Energy cost per unit of movement
	SensingCostBase       float64    `json:"sensingCostBase"`       // Energy cost for sensor operations
	OptimalEnergyGainRate float64    `json:"optimalEnergyGainRate"` // Maximum energy gain per second
	EnergyEfficiencyRange [2]float64 `json:"energyEfficiencyRange"` // Min/max for random initialization
}

// ReproductionConfig holds settings for the reproduction system
type ReproductionConfig struct {
	ReproductionThreshold float64 `json:"reproductionThreshold"` // Energy required to reproduce
	EnergyTransferRatio   float64 `json:"energyTransferRatio"`   // Portion of energy given to offspring
	OffspringDistance     float64 `json:"offspringDistance"`     // How far offspring spawns from parent
	MutationRate          float64 `json:"mutationRate"`          // Probability of trait mutation
	MutationMagnitude     float64 `json:"mutationMagnitude"`     // Maximum percent change when mutation occurs
	MaxPopulation         int     `json:"maxPopulation"`         // Optional cap on total population
}

// ChemicalConfig holds settings for chemical sources
type ChemicalConfig struct {
	Count          int     `json:"count"`
	MinStrength    float64 `json:"minStrength"`
	MaxStrength    float64 `json:"maxStrength"`
	MinDecayFactor float64 `json:"minDecayFactor"`
	MaxDecayFactor float64 `json:"maxDecayFactor"`
	// New fields for energy balance
	DepletionRate           float64 `json:"depletionRate"`
	RegenerationProbability float64 `json:"regenerationProbability"`
	TargetSystemEnergy      float64 `json:"targetSystemEnergy"`
}

// RenderConfig holds settings for visualization
type RenderConfig struct {
	WindowWidth  int  `json:"windowWidth"`
	WindowHeight int  `json:"windowHeight"`
	FrameRate    int  `json:"frameRate"`
	ShowGrid     bool `json:"showGrid"`
	ShowSensors  bool `json:"showSensors"`
	ShowLegend   bool `json:"showLegend"`
}

// SimulationConfig holds all configuration for the simulation
type SimulationConfig struct {
	Version         string             `json:"version"`
	World           WorldConfig        `json:"world"`
	Organism        OrganismConfig     `json:"organism"`
	Chemical        ChemicalConfig     `json:"chemical"`
	Render          RenderConfig       `json:"render"`
	Energy          EnergyConfig       `json:"energy"`       // New energy configuration
	Reproduction    ReproductionConfig `json:"reproduction"` // New reproduction configuration
	RandomSeed      int64              `json:"randomSeed"`
	SimulationSpeed float64            `json:"simulationSpeed"`
}

// DefaultConfig returns a default configuration with reasonable values
func DefaultConfig() SimulationConfig {
	return SimulationConfig{
		Version: Version,
		World: WorldConfig{
			Width:  1000.0,
			Height: 1000.0,
		},
		Organism: OrganismConfig{
			Count:                        100,
			Speed:                        2.0,
			SensorDistance:               10.0,
			TurnSpeed:                    math.Pi / 10, // 18 degrees per step
			PreferenceDistributionMean:   50.0,
			PreferenceDistributionStdDev: 10.0,
		},
		Energy: EnergyConfig{
			InitialEnergy:         80.0,                 // Start with 80% of maximum
			MaximumEnergy:         100.0,                // Base energy capacity
			BaseMetabolicRate:     0.1,                  // Energy consumed per second just existing
			MovementCostFactor:    0.02,                 // Energy cost per unit of movement
			SensingCostBase:       0.01,                 // Energy cost for sensing operations
			OptimalEnergyGainRate: 0.5,                  // Maximum energy gain per second
			EnergyEfficiencyRange: [2]float64{0.8, 1.2}, // Range for random efficiency
		},
		Reproduction: ReproductionConfig{
			ReproductionThreshold: 0.75, // 75% of max energy required to reproduce
			EnergyTransferRatio:   0.3,  // 30% of energy given to offspring
			OffspringDistance:     10.0, // Units away from parent
			MutationRate:          0.2,  // 20% chance of mutation per trait
			MutationMagnitude:     0.1,  // 10% maximum change when mutation occurs
			MaxPopulation:         500,  // Maximum allowed population
		},
		Chemical: ChemicalConfig{
			Count:          5,
			MinStrength:    100.0,
			MaxStrength:    500.0,
			MinDecayFactor: 0.001,
			MaxDecayFactor: 0.01,
			// Default values for energy balance
			DepletionRate:           0.2,
			RegenerationProbability: 0.2,
			TargetSystemEnergy:      10000.0,
		},
		Render: RenderConfig{
			WindowWidth:  800,
			WindowHeight: 800,
			FrameRate:    60,
			ShowGrid:     true,
			ShowSensors:  true,
			ShowLegend:   true,
		},
		RandomSeed:      0, // 0 means use current time as seed
		SimulationSpeed: 10.0,
	}
}

// LoadFromFile loads configuration from a JSON file
func LoadFromFile(filename string) (SimulationConfig, error) {
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

// SaveToFile saves configuration to a JSON file
func SaveToFile(config SimulationConfig, filename string) error {
	// Convert to JSON with indentation for readability
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	// Write to file
	return os.WriteFile(filename, data, 0644)
}
