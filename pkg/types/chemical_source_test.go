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
}

func TestChemicalSourceEdgeCases(t *testing.T) {
	// Test with zero strength
	csZeroStrength := NewChemicalSource(NewPoint(0, 0), 0.0, 0.1)
	if csZeroStrength.GetConcentrationAt(NewPoint(5, 5)) != 0 {
		t.Error("Concentration with zero strength should be zero")
	}

	// Test with zero decay factor
	csZeroDecay := NewChemicalSource(NewPoint(0, 0), 100.0, 0.0)
	if csZeroDecay.GetConcentrationAt(NewPoint(100, 100)) != 100.0 {
		t.Error("Concentration with zero decay factor should equal strength at any distance")
	}
}
