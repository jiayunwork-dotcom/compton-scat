package section

import (
	"testing"

	"compton-scat/internal/constants"
)

// TestDifferential90 checks the differential cross section at 90
// degrees against the Klein-Nishina closed form for a 511 keV photon,
// where k = E/me*c^2 is close to one.
func TestDifferential90(t *testing.T) {
	e := constants.KEVToJoules(511.0)
	re := constants.ElectronClassicalRadius()
	k := constants.EnergyRatio(e)
	p := 1.0 / (1.0 + k)
	want := (re * re / 2.0) * p * p * (1.0/p + p - 1.0)
	got := Differential(e, 90)
	rel := (got - want) / want
	if rel > 1e-9 || rel < -1e-9 {
		t.Errorf("Differential(90) = %g m^2/sr, want %g m^2/sr (rel %.2e)", got, want, rel)
	}
}

// TestForwardConstant checks that the forward differential cross
// section equals r_e^2 at several energies.
func TestForwardConstant(t *testing.T) {
	for _, kev := range []float64{1, 10, 511, 1000, 5000} {
		e := constants.KEVToJoules(kev)
		got, want, rel := VerifyForwardConstant(e)
		if rel > 1e-12 {
			t.Errorf("forward dsigma at %g keV = %g m^2/sr, want r_e^2 = %g m^2/sr", kev, got, want)
		}
	}
}

// TestThomsonLimit checks that the total cross section at low energy
// approaches the Thomson value.
func TestThomsonLimit(t *testing.T) {
	e := constants.KEVToJoules(1.0)
	ratio := Total(e, DefaultIntegrationNodes) / ThomsonCrossSection()
	if !(ratio < 1.0 && ratio > 0.99) {
		t.Errorf("sigma(1 keV)/sigma_Thomson = %g, want in (0.99, 1)", ratio)
	}
}

// TestTotalDecreasesWithEnergy checks the central property: the total
// cross section falls as the incident energy rises.
func TestTotalDecreasesWithEnergy(t *testing.T) {
	low := constants.KEVToJoules(10.0)
	high := constants.KEVToJoules(1000.0)
	sigmaLow := Total(low, DefaultIntegrationNodes)
	sigmaHigh := Total(high, DefaultIntegrationNodes)
	if !(sigmaHigh < sigmaLow) {
		t.Errorf("sigma(1 MeV) = %g m^2, sigma(10 keV) = %g m^2, want sigma to drop", sigmaHigh, sigmaLow)
	}
	if !(sigmaLow < ThomsonCrossSection()) {
		t.Errorf("sigma(10 keV) = %g m^2, want below Thomson %g m^2", sigmaLow, ThomsonCrossSection())
	}
}

// TestIntegrationConvergence checks that doubling the panels changes
// the total cross section by a negligible amount.
func TestIntegrationConvergence(t *testing.T) {
	e := constants.KEVToJoules(511.0)
	coarse, _ := TotalError(e, 2000)
	fine := Total(e, 4000)
	rel := (coarse - fine) / fine
	if rel < 0 {
		rel = -rel
	}
	if rel > 1e-7 {
		t.Errorf("2000 vs 4000 panels differ by rel %.2e, want below 1e-7", rel)
	}
}

// TestNegativeEnergy checks that a non-positive energy is rejected.
func TestNegativeEnergy(t *testing.T) {
	for _, e := range []float64{0, -1} {
		if _, err := CrossSection(e); err == nil {
			t.Errorf("CrossSection(%g) accepted, want error", e)
		}
	}
}

// TestTotalAtRatio checks the dimensionless ratio form agrees with the
// absolute energy form.
func TestTotalAtRatio(t *testing.T) {
	e := constants.KEVToJoules(511.0)
	k := constants.EnergyRatio(e)
	viaRatio := TotalAtRatio(k, DefaultIntegrationNodes)
	direct := Total(e, DefaultIntegrationNodes)
	rel := (viaRatio - direct) / direct
	if rel < 0 {
		rel = -rel
	}
	if rel > 1e-9 {
		t.Errorf("TotalAtRatio = %g m^2, Total = %g m^2, rel %.2e", viaRatio, direct, rel)
	}
}

// TestRunAllChecks verifies that every cross-section property passes
// for the default integration.
func TestRunAllChecks(t *testing.T) {
	results := RunAll(DefaultIntegrationNodes)
	if !AllPass(results) {
		t.Errorf("section checks failed: %v", Failed(results))
	}
}

// TestBackwardDifferential checks the backscatter differential section
// stays positive and finite at high energy.
func TestBackwardDifferential(t *testing.T) {
	for _, kev := range []float64{1, 511, 10000} {
		e := constants.KEVToJoules(kev)
		v := DifferentialBackward(e)
		if !(v > 0) {
			t.Errorf("DifferentialBackward(%g keV) = %g m^2/sr, want positive", kev, v)
		}
		if v > DifferentialForward(e) {
			t.Errorf("backscatter %g m^2/sr exceeds forward %g m^2/sr at %g keV", v, DifferentialForward(e), kev)
		}
	}
}
