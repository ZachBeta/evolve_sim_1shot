package types

import (
	"math"
	"testing"
)

func TestNewChemicalSource(t *testing.T) {
	position := NewPoint(10, 20)
	cs := NewChemicalSource(position, 100.0, 0.1)

	if cs.Position.X != 10 || cs.Position.Y != 20 {
		t.Errorf("ChemicalSource position = %v; want {X:10, Y:20}", cs.Position)
	}

	if cs.Strength != 100.0 {
		t.Errorf("ChemicalSource strength = %v; want 100.0", cs.Strength)
	}

	if cs.DecayFactor != 0.1 {
		t.Errorf("ChemicalSource decayFactor = %v; want 0.1", cs.DecayFactor)
	}

	// Test new energy-related fields
	expectedMaxEnergy := 100.0 * 1000
	if cs.MaxEnergy != expectedMaxEnergy {
		t.Errorf("ChemicalSource maxEnergy = %v; want %v", cs.MaxEnergy, expectedMaxEnergy)
	}

	if cs.Energy != expectedMaxEnergy {
		t.Errorf("ChemicalSource energy = %v; want %v", cs.Energy, expectedMaxEnergy)
	}

	if !cs.IsActive {
		t.Errorf("ChemicalSource isActive = %v; want true", cs.IsActive)
	}

	if cs.DepletionRate != 0.2 {
		t.Errorf("ChemicalSource depletionRate = %v; want 0.2", cs.DepletionRate)
	}
}

func TestChemicalSourceGetConcentrationAt(t *testing.T) {
	cs := NewChemicalSource(NewPoint(0, 0), 100.0, 0.1)

	// Test cases: different points and expected concentrations
	testCases := []struct {
		point    Point
		expected float64
	}{
		{NewPoint(0, 0), 100.0},                       // At source
		{NewPoint(1, 0), 100.0 / (1 + 0.1)},           // Distance = 1
		{NewPoint(3, 4), 100.0 / (1 + 25*0.1)},        // Distance = 5
		{NewPoint(10, 0), 100.0 / (1 + 100*0.1)},      // Distance = 10
		{NewPoint(100, 100), 100.0 / (1 + 20000*0.1)}, // Far away
	}

	for i, tc := range testCases {
		concentration := cs.GetConcentrationAt(tc.point)
		if math.Abs(concentration-tc.expected) > 1e-9 {
			t.Errorf("Case %d: Concentration at %v = %v; want %v",
				i, tc.point, concentration, tc.expected)
		}
	}

	// Test with half energy
	cs.Energy = cs.MaxEnergy / 2
	for i, tc := range testCases {
		halfConcentration := cs.GetConcentrationAt(tc.point)
		expectedHalfConcentration := tc.expected * 0.5
		if math.Abs(halfConcentration-expectedHalfConcentration) > 1e-9 {
			t.Errorf("Case %d (half energy): Concentration at %v = %v; want %v",
				i, tc.point, halfConcentration, expectedHalfConcentration)
		}
	}

	// Test inactive source
	cs.IsActive = false
	for i, tc := range testCases {
		inactiveConcentration := cs.GetConcentrationAt(tc.point)
		if inactiveConcentration != 0 {
			t.Errorf("Case %d (inactive): Concentration at %v = %v; want 0",
				i, tc.point, inactiveConcentration)
		}
	}
}

func TestChemicalSourceEdgeCases(t *testing.T) {
	// Test with zero strength
	csZeroStrength := NewChemicalSource(NewPoint(0, 0), 0.0, 0.1)
	concentration := csZeroStrength.GetConcentrationAt(NewPoint(5, 5))
	if concentration != 0 {
		t.Errorf("Concentration with zero strength = %v; want 0", concentration)
	}

	// Test with zero decay factor
	csZeroDecay := NewChemicalSource(NewPoint(0, 0), 100.0, 0.0)
	expectedConcentration := 100.0 // Full strength at any distance
	actualConcentration := csZeroDecay.GetConcentrationAt(NewPoint(100, 100))
	if math.Abs(actualConcentration-expectedConcentration) > 1e-9 {
		t.Errorf("Concentration with zero decay factor = %v; want %v",
			actualConcentration, expectedConcentration)
	}
}

func TestChemicalSourceDepletion(t *testing.T) {
	cs := NewChemicalSource(NewPoint(0, 0), 100.0, 0.1)
	initialEnergy := cs.Energy
	var worldEnergy float64 = initialEnergy

	// Update for 10 seconds
	deltaTime := 10.0
	cs.Update(deltaTime, &worldEnergy)

	// Calculate expected depletion
	expectedDepletion := deltaTime * cs.DepletionRate
	expectedEnergy := initialEnergy - expectedDepletion

	if math.Abs(cs.Energy-expectedEnergy) > 1e-9 {
		t.Errorf("ChemicalSource energy after depletion = %v; want %v", cs.Energy, expectedEnergy)
	}

	if math.Abs(worldEnergy-expectedEnergy) > 1e-9 {
		t.Errorf("World energy after depletion = %v; want %v", worldEnergy, expectedEnergy)
	}

	// Source should still be active
	if !cs.IsActive {
		t.Error("ChemicalSource should still be active after partial depletion")
	}
}

func TestChemicalSourceDepleteToInactive(t *testing.T) {
	cs := NewChemicalSource(NewPoint(0, 0), 100.0, 0.1)
	cs.Energy = 1.0 // Just a tiny bit of energy left
	var worldEnergy float64 = 1.0

	// Update for enough time to fully deplete
	cs.Update(10.0, &worldEnergy)

	// Energy should be 0 and source should be inactive
	if cs.Energy != 0 {
		t.Errorf("ChemicalSource energy after full depletion = %v; want 0", cs.Energy)
	}

	if cs.IsActive {
		t.Error("ChemicalSource should be inactive after full depletion")
	}

	if worldEnergy != 0 {
		t.Errorf("World energy after full depletion = %v; want 0", worldEnergy)
	}

	// Verify that further updates don't affect energy
	cs.Update(10.0, &worldEnergy)
	if cs.Energy != 0 || worldEnergy != 0 {
		t.Errorf("Update on depleted source changed energy values: source=%v, world=%v",
			cs.Energy, worldEnergy)
	}
}
