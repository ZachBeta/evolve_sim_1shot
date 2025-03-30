package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	// Check a few key values to ensure defaults are set
	if config.World.Width != 1000.0 {
		t.Errorf("Default world width = %v; want 1000.0", config.World.Width)
	}

	if config.Organism.Count != 100 {
		t.Errorf("Default organism count = %v; want 100", config.Organism.Count)
	}

	if config.Chemical.Count != 5 {
		t.Errorf("Default chemical source count = %v; want 5", config.Chemical.Count)
	}

	if config.Render.FrameRate != 60 {
		t.Errorf("Default frame rate = %v; want 60", config.Render.FrameRate)
	}
}

func TestSaveAndLoadConfig(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "config_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a custom config
	config := DefaultConfig()
	config.World.Width = 2000.0
	config.Organism.Count = 200
	config.Chemical.MaxStrength = 1000.0
	config.Render.WindowWidth = 1024

	// Save to a temporary file
	tempFile := filepath.Join(tempDir, "test_config.json")
	err = SaveToFile(config, tempFile)
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Load the config
	loadedConfig, err := LoadFromFile(tempFile)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Check if loaded values match
	if loadedConfig.World.Width != 2000.0 {
		t.Errorf("Loaded world width = %v; want 2000.0", loadedConfig.World.Width)
	}

	if loadedConfig.Organism.Count != 200 {
		t.Errorf("Loaded organism count = %v; want 200", loadedConfig.Organism.Count)
	}

	if loadedConfig.Chemical.MaxStrength != 1000.0 {
		t.Errorf("Loaded max chemical strength = %v; want 1000.0", loadedConfig.Chemical.MaxStrength)
	}

	if loadedConfig.Render.WindowWidth != 1024 {
		t.Errorf("Loaded window width = %v; want 1024", loadedConfig.Render.WindowWidth)
	}
}

func TestLoadMissingConfig(t *testing.T) {
	// Try to load a non-existent file
	config, err := LoadFromFile("non_existent_file.json")

	// Should return an error
	if err == nil {
		t.Error("Expected error when loading non-existent file, got nil")
	}

	// Should still return default config
	if config.World.Width != 1000.0 {
		t.Errorf("Failed load should return default world width, got %v", config.World.Width)
	}
}

func TestLoadInvalidConfig(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "config_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create invalid JSON file
	tempFile := filepath.Join(tempDir, "invalid_config.json")
	err = os.WriteFile(tempFile, []byte("{invalid json}"), 0644)
	if err != nil {
		t.Fatalf("Failed to write invalid config: %v", err)
	}

	// Try to load the invalid file
	config, err := LoadFromFile(tempFile)

	// Should return an error
	if err == nil {
		t.Error("Expected error when loading invalid JSON, got nil")
	}

	// Should still return default config
	if config.World.Width != 1000.0 {
		t.Errorf("Failed load should return default world width, got %v", config.World.Width)
	}
}

func TestPartialConfig(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "config_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a partial config JSON (only modifies some values)
	partialConfig := `{
		"world": {
			"width": 1500.0
		},
		"render": {
			"frameRate": 30
		}
	}`

	tempFile := filepath.Join(tempDir, "partial_config.json")
	err = os.WriteFile(tempFile, []byte(partialConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to write partial config: %v", err)
	}

	// Load the partial config
	config, err := LoadFromFile(tempFile)
	if err != nil {
		t.Fatalf("Failed to load partial config: %v", err)
	}

	// Check that specified values were loaded
	if config.World.Width != 1500.0 {
		t.Errorf("Loaded world width = %v; want 1500.0", config.World.Width)
	}

	if config.Render.FrameRate != 30 {
		t.Errorf("Loaded frame rate = %v; want 30", config.Render.FrameRate)
	}

	// Check that unspecified values remained at defaults
	if config.World.Height != 1000.0 {
		t.Errorf("World height should remain at default 1000.0, got %v", config.World.Height)
	}

	if config.Organism.Count != 100 {
		t.Errorf("Organism count should remain at default 100, got %v", config.Organism.Count)
	}
}
