package constants

import (
	"math"
	"testing"
)

func closeTo(got, want, tol float64) bool {
	return math.Abs(got-want) <= tol
}

func closeRel(got, want, rel float64) bool {
	if want == 0 {
		return math.Abs(got) <= rel
	}
	return math.Abs(got-want) <= math.Abs(want)*rel
}

func TestComptonWavelength(t *testing.T) {
	got := ComptonWavelength()
	want := 2.42631023867e-12
	if !closeTo(got, want, 5e-23) {
		t.Errorf("ComptonWavelength() = %.6e m, want %.6e m", got, want)
	}
}

func TestElectronRestEnergy(t *testing.T) {
	ev := ElectronRestEnergyEV()
	wantEV := 510998.95
	if !closeTo(ev, wantEV, 1e-2) {
		t.Errorf("ElectronRestEnergyEV() = %.3f eV, want %.3f eV", ev, wantEV)
	}
	kev := ElectronRestEnergyKEV()
	if !closeTo(kev, wantEV/1e3, 1e-5) {
		t.Errorf("ElectronRestEnergyKEV() = %.6f keV, want %.6f keV", kev, wantEV/1e3)
	}
}

func TestSelfConsistency(t *testing.T) {
	c := SelfConsistency()
	if !c.OK(1e-12) {
		t.Errorf("lambda_c*me*c^2 = %.6e differs from h*c = %.6e by %.3e", c.ComptonTimesRest, c.PlanchTimesC, c.RelativeDiscrepancy)
	}
}

func TestPlanckTimesC(t *testing.T) {
	got := PlanckTimesC()
	want := 1.9864458571489286e-25
	if !closeTo(got, want, 1e-33) {
		t.Errorf("PlanckTimesC() = %.6e J*m, want %.6e J*m", got, want)
	}
}

func TestUnitConversions(t *testing.T) {
	j := 1.0
	ev := JoulesToEV(j)
	if !closeRel(ev, 1.0/ElementaryCharge, 1e-12) {
		t.Errorf("JoulesToEV(1) = %.6e eV, want %.6e eV", ev, 1.0/ElementaryCharge)
	}
	back := EVToJoules(ev)
	if !closeRel(back, j, 1e-12) {
		t.Errorf("round trip = %.6e J, want %.6e J", back, j)
	}
	kev := KEVToJoules(511.0)
	if !closeRel(kev, EVToJoules(511000), 1e-12) {
		t.Errorf("KEVToJoules(511) = %.6e J, want %.6e J", kev, EVToJoules(511000))
	}
}

func TestClassicalRadius(t *testing.T) {
	got := ElectronClassicalRadius()
	want := 2.8179403262e-15
	if !closeTo(got, want, 1e-25) {
		t.Errorf("ElectronClassicalRadius() = %.6e m, want %.6e m", got, want)
	}
}

func TestEnergyRatioAt511(t *testing.T) {
	ratio := EnergyRatio(KEVToJoules(511.0))
	if !closeTo(ratio, 1.0, 1e-3) {
		t.Errorf("EnergyRatio(511 keV) = %.6f, want 1.0", ratio)
	}
}
