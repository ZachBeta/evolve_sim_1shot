package types

import (
	"math"
	"testing"
)

func TestNewWorld(t *testing.T) {
	world := NewWorld(100, 200)

	if world.Width != 100 {
		t.Errorf("World width = %v; want 100", world.Width)
	}

	if world.Height != 200 {
		t.Errorf("World height = %v; want 200", world.Height)
	}

	if len(world.Organisms) != 0 {
		t.Errorf("Initial organism count = %v; want 0", len(world.Organisms))
	}

	if len(world.ChemicalSources) != 0 {
		t.Errorf("Initial chemical source count = %v; want 0", len(world.ChemicalSources))
	}

	bounds := world.GetWorldBounds()
	if bounds.X != 0 || bounds.Y != 0 || bounds.Width != 100 || bounds.Height != 200 {
		t.Errorf("World bounds = %v; want {X:0, Y:0, Width:100, Height:200}", bounds)
	}
}

func TestAddOrganism(t *testing.T) {
	world := NewWorld(100, 100)

	// Test adding organism within bounds
	org1 := NewOrganism(NewPoint(50, 50), 0, 5.0, 1.0, DefaultSensorAngles())
	if !world.AddOrganism(org1) {
		t.Error("Failed to add organism within bounds")
	}

	if world.OrganismCount() != 1 {
		t.Errorf("Organism count after adding = %v; want 1", world.OrganismCount())
	}

	// Test adding organism outside bounds
	org2 := NewOrganism(NewPoint(150, 50), 0, 5.0, 1.0, DefaultSensorAngles())
	if world.AddOrganism(org2) {
		t.Error("Should not add organism outside bounds")
	}

	if world.OrganismCount() != 1 {
		t.Errorf("Organism count after failed add = %v; want 1", world.OrganismCount())
	}
}

func TestAddChemicalSource(t *testing.T) {
	world := NewWorld(100, 100)

	// Test adding source within bounds
	source1 := NewChemicalSource(NewPoint(50, 50), 100.0, 0.1)
	if !world.AddChemicalSource(source1) {
		t.Error("Failed to add chemical source within bounds")
	}

	if world.ChemicalSourceCount() != 1 {
		t.Errorf("Source count after adding = %v; want 1", world.ChemicalSourceCount())
	}

	// Test adding source outside bounds
	source2 := NewChemicalSource(NewPoint(150, 50), 100.0, 0.1)
	if world.AddChemicalSource(source2) {
		t.Error("Should not add chemical source outside bounds")
	}

	if world.ChemicalSourceCount() != 1 {
		t.Errorf("Source count after failed add = %v; want 1", world.ChemicalSourceCount())
	}
}

func TestGetConcentrationAt(t *testing.T) {
	world := NewWorld(100, 100)

	// Add two chemical sources
	source1 := NewChemicalSource(NewPoint(25, 25), 100.0, 0.1)
	source2 := NewChemicalSource(NewPoint(75, 75), 50.0, 0.2)
	world.AddChemicalSource(source1)
	world.AddChemicalSource(source2)

	// Calculate expected values directly
	point1 := NewPoint(25, 25)
	point2 := NewPoint(75, 75)
	point3 := NewPoint(50, 50)
	point4 := NewPoint(0, 0)

	// Calculate the actual concentrations
	conc1 := world.GetConcentrationAt(point1)
	conc2 := world.GetConcentrationAt(point2)
	conc3 := world.GetConcentrationAt(point3)
	conc4 := world.GetConcentrationAt(point4)

	// Now verify that concentrations are as expected
	// At source 1
	expected1 := source1.GetConcentrationAt(point1) + source2.GetConcentrationAt(point1)
	if !approximatelyEqual(conc1, expected1, 1e-9) {
		t.Errorf("Concentration at source 1 = %v; want %v", conc1, expected1)
	}

	// At source 2
	expected2 := source1.GetConcentrationAt(point2) + source2.GetConcentrationAt(point2)
	if !approximatelyEqual(conc2, expected2, 1e-9) {
		t.Errorf("Concentration at source 2 = %v; want %v", conc2, expected2)
	}

	// Between sources
	expected3 := source1.GetConcentrationAt(point3) + source2.GetConcentrationAt(point3)
	if !approximatelyEqual(conc3, expected3, 1e-9) {
		t.Errorf("Concentration between sources = %v; want %v", conc3, expected3)
	}

	// Far from sources
	expected4 := source1.GetConcentrationAt(point4) + source2.GetConcentrationAt(point4)
	if !approximatelyEqual(conc4, expected4, 1e-9) {
		t.Errorf("Concentration far from sources = %v; want %v", conc4, expected4)
	}

	// Test with no sources
	emptyWorld := NewWorld(100, 100)
	if emptyWorld.GetConcentrationAt(NewPoint(50, 50)) != 0 {
		t.Error("Concentration in world with no sources should be 0")
	}
}

// Helper function to check if two float64 values are approximately equal
func approximatelyEqual(a, b, epsilon float64) bool {
	diff := math.Abs(a - b)
	return diff < epsilon
}
