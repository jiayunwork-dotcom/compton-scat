package kinematics

import (
	"errors"
	"math"
	"testing"

	"compton-scat/internal/constants"
)

func kevToJ(t *testing.T, kev float64) float64 {
	t.Helper()
	return constants.KEVToJoules(kev)
}

func TestThetaZero(t *testing.T) {
	in := Input{EnergyJ: kevToJ(t, 511), ThetaDeg: 0}
	r, err := Solve(in)
	if err != nil {
		t.Fatalf("Solve: %v", err)
	}
	if r.DeltaLambda != 0 {
		t.Errorf("DeltaLambda = %g m, want 0", r.DeltaLambda)
	}
	if relErr(r.ScatteredEnergy, in.EnergyJ) > 1e-12 {
		t.Errorf("ScatteredEnergy = %g J, want %g J", r.ScatteredEnergy, in.EnergyJ)
	}
	if relErr(r.RecoilEnergy, 0) > 1e-12 {
		t.Errorf("RecoilEnergy = %g J, want 0", r.RecoilEnergy)
	}
}

func TestTheta180(t *testing.T) {
	in := Input{EnergyJ: kevToJ(t, 511), ThetaDeg: 180}
	r, err := Solve(in)
	if err != nil {
		t.Fatalf("Solve: %v", err)
	}
	lambdaC := constants.ComptonWavelength()
	if relErr(r.DeltaLambda, 2*lambdaC) > 1e-12 {
		t.Errorf("DeltaLambda = %g m, want 2*lambda_c = %g m", r.DeltaLambda, 2*lambdaC)
	}
	for _, theta := range []float64{0, 30, 90, 150} {
		other, err := Solve(Input{EnergyJ: in.EnergyJ, ThetaDeg: theta})
		if err != nil {
			t.Fatalf("Solve(%g): %v", theta, err)
		}
		if other.ScatteredEnergy < r.ScatteredEnergy {
			t.Errorf("E'(%.0f deg) = %g J is below E'(180 deg) = %g J", theta, other.ScatteredEnergy, r.ScatteredEnergy)
		}
	}
}

func TestNinetyDegreeClosedForm(t *testing.T) {
	energyJ := kevToJ(t, 511)
	in := Input{EnergyJ: energyJ, ThetaDeg: 90}
	r, err := Solve(in)
	if err != nil {
		t.Fatalf("Solve: %v", err)
	}
	ratio := constants.EnergyRatio(energyJ)
	want := energyJ / (1.0 + ratio)
	gotKEV := constants.JoulesToKEV(r.ScatteredEnergy)
	wantKEV := constants.JoulesToKEV(want)
	if relErr(r.ScatteredEnergy, want) > 1e-9 {
		t.Errorf("E' = %g keV, want %g keV from closed form", gotKEV, wantKEV)
	}
	if wantKEV < 255 || wantKEV > 256 {
		t.Errorf("closed form E' = %g keV, want about 255.5 keV", wantKEV)
	}
}

func TestEnergyConservation(t *testing.T) {
	for _, tc := range []struct {
		kev, deg float64
	}{
		{10, 30}, {511, 90}, {1000, 45}, {10, 180}, {0.1, 120}, {2000, 175},
	} {
		in := Input{EnergyJ: kevToJ(t, tc.kev), ThetaDeg: tc.deg}
		r, err := Solve(in)
		if err != nil {
			t.Fatalf("Solve(%g keV, %g deg): %v", tc.kev, tc.deg, err)
		}
		sum := r.RecoilEnergy + r.ScatteredEnergy
		if relErr(sum, in.EnergyJ) > 1e-9 {
			t.Errorf("Ke+E' = %g J, want %g J for E=%g keV theta=%g", sum, in.EnergyJ, tc.kev, tc.deg)
		}
	}
}

func TestWavelengthConsistency(t *testing.T) {
	in := Input{EnergyJ: kevToJ(t, 511), ThetaDeg: 90}
	r, err := Solve(in)
	if err != nil {
		t.Fatalf("Solve: %v", err)
	}
	w := CheckWavelengthConsistency(r)
	if !w.OK(1e-9) {
		t.Errorf("wavelength identity max rel err = %g, want below 1e-9", w.MaxRelativeError)
	}
}

func TestMomentumConservation(t *testing.T) {
	in := Input{EnergyJ: kevToJ(t, 511), ThetaDeg: 90}
	r, err := Solve(in)
	if err != nil {
		t.Fatalf("Solve: %v", err)
	}
	m := CheckMomentumConservation(r)
	if !m.OK(1e-9) {
		t.Errorf("momentum check rel diff = %g, want below 1e-9", m.RelativeDiscrepancy)
	}
}

func TestInvalidEnergy(t *testing.T) {
	for _, e := range []float64{0, -1, -511} {
		in := Input{EnergyJ: e, ThetaDeg: 90}
		_, err := Solve(in)
		if err == nil {
			t.Errorf("Solve(E=%g) accepted, want error", e)
			continue
		}
		if !errors.Is(err, ErrNonPositiveEnergy) {
			t.Errorf("Solve(E=%g) error = %v, want ErrNonPositiveEnergy", e, err)
		}
	}
}

func TestInvalidAngle(t *testing.T) {
	for _, deg := range []float64{-1, 200, 360, -90} {
		in := Input{EnergyJ: kevToJ(t, 511), ThetaDeg: deg}
		_, err := Solve(in)
		if err == nil {
			t.Errorf("Solve(theta=%g) accepted, want error", deg)
			continue
		}
		if !errors.Is(err, ErrAngleOutOfRange) {
			t.Errorf("Solve(theta=%g) error = %v, want ErrAngleOutOfRange", deg, err)
		}
	}
}

func TestRelativeLossIncreases(t *testing.T) {
	low := Input{EnergyJ: kevToJ(t, 10), ThetaDeg: 90}
	high := Input{EnergyJ: kevToJ(t, 1000), ThetaDeg: 90}
	rLow, err := Solve(low)
	if err != nil {
		t.Fatalf("Solve: %v", err)
	}
	rHigh, err := Solve(high)
	if err != nil {
		t.Fatalf("Solve: %v", err)
	}
	if rLow.LossFraction() >= rHigh.LossFraction() {
		t.Errorf("loss(E=10 keV) = %g, loss(E=1 MeV) = %g, want loss to grow", rLow.LossFraction(), rHigh.LossFraction())
	}
}

func TestHighEnergyLargeAngle(t *testing.T) {
	in := Input{EnergyJ: kevToJ(t, 10000), ThetaDeg: 179}
	r, err := Solve(in)
	if err != nil {
		t.Fatalf("Solve: %v", err)
	}
	if !(r.ScatteredEnergy > 0) {
		t.Errorf("E' = %g J, want positive for 10 MeV backscatter", r.ScatteredEnergy)
	}
	if math.IsNaN(r.ScatteredEnergy) || math.IsInf(r.ScatteredEnergy, 0) {
		t.Errorf("E' = %v, want finite", r.ScatteredEnergy)
	}
	ratio := constants.EnergyRatio(in.EnergyJ)
	cosT := CosTheta(in.ThetaDeg)
	want := in.EnergyJ / (1.0 + ratio*(1.0-cosT))
	if relErr(r.ScatteredEnergy, want) > 1e-9 {
		t.Errorf("E' = %g keV, want %g keV", constants.JoulesToKEV(r.ScatteredEnergy), constants.JoulesToKEV(want))
	}
}

func TestScatteredEnergyHelper(t *testing.T) {
	in := Input{EnergyJ: kevToJ(t, 511), ThetaDeg: 90}
	r, err := Solve(in)
	if err != nil {
		t.Fatalf("Solve: %v", err)
	}
	sc, err := ScatteredEnergy(in)
	if err != nil {
		t.Fatalf("ScatteredEnergy: %v", err)
	}
	if relErr(sc, r.ScatteredEnergy) > 1e-12 {
		t.Errorf("ScatteredEnergy = %g J, want %g J from Solve", sc, r.ScatteredEnergy)
	}
}

func TestRunAllChecks(t *testing.T) {
	results, err := RunAll(Input{EnergyJ: kevToJ(t, 511), ThetaDeg: 90})
	if err != nil {
		t.Fatalf("RunAll: %v", err)
	}
	if !AllPass(results) {
		t.Errorf("checks failed: %v", Failed(results))
	}
	if _, err := RunAll(Input{EnergyJ: 0, ThetaDeg: 90}); err == nil {
		t.Error("RunAll accepted zero energy, want error")
	}
}
