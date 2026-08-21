package kinematics

import (
	"fmt"

	"compton-scat/internal/constants"
)

// CheckResult is one named assertion about the kinematics of an event.
// Every item has a Pass flag and a human readable detail line that the
// checks command prints regardless of the outcome, so a failing rule is
// visible on stderr instead of silently passing.
type CheckResult struct {
	// Name identifies the rule being checked.
	Name string
	// Pass reports whether the rule held within tolerance.
	Pass bool
	// Detail describes the numeric values compared.
	Detail string
}

// RunAll executes every applicable cross rule for an input event and
// returns the results in a fixed order. Boundary rules (forward and
// backscatter) only apply when the angle is exactly the boundary value;
// the other rules apply to every valid event.
func RunAll(in Input) ([]CheckResult, error) {
	r, err := Solve(in)
	if err != nil {
		return nil, err
	}
	results := []CheckResult{
		checkConstantSelfConsistency(),
		checkEnergyConservation(r),
		checkWavelengthIdentity(r),
		checkMomentumIdentity(r),
		checkRelativeLossTrend(in),
	}
	if r.IsRightAngle() {
		results = append(results, checkNinetyDegreeClosedForm(r))
	}
	if r.IsForwardScatter() {
		results = append(results, checkForwardNoShift(r))
	}
	if r.IsBackscatter() {
		results = append(results, checkBackscatterMaximum(r))
	}
	return results, nil
}

// AllPass reports whether every check result passed.
func AllPass(results []CheckResult) bool {
	for _, c := range results {
		if !c.Pass {
			return false
		}
	}
	return true
}

// Failed returns the names of the failed checks, or nil when all pass.
func Failed(results []CheckResult) []string {
	var failed []string
	for _, c := range results {
		if !c.Pass {
			failed = append(failed, c.Name)
		}
	}
	return failed
}

// checkConstantSelfConsistency verifies lambda_c*me*c^2 == h*c using
// the constants package, the base identity every other formula relies
// on.
func checkConstantSelfConsistency() CheckResult {
	c := constants.SelfConsistency()
	ok := c.OK(1e-9)
	return CheckResult{
		Name:   "constant self-consistency",
		Pass:   ok,
		Detail: c.String(),
	}
}

// checkEnergyConservation verifies Ke + E' == E for the solved event.
func checkEnergyConservation(r Result) CheckResult {
	sum := r.EnergySum()
	detail := fmt.Sprintf("Ke + E' = %.12g J, E = %.12g J", sum, r.Input.EnergyJ)
	ok := relErr(sum, r.Input.EnergyJ) <= DefaultTolerance
	return CheckResult{Name: "energy conservation", Pass: ok, Detail: detail}
}

// checkWavelengthIdentity verifies lambda' = lambda + dlambda and
// E' = h*c/lambda' for the solved event.
func checkWavelengthIdentity(r Result) CheckResult {
	w := CheckWavelengthConsistency(r)
	detail := fmt.Sprintf("max relative error %.3e across lambda/E identities", w.MaxRelativeError)
	ok := w.OK(1e-6)
	return CheckResult{Name: "wavelength/energy identity", Pass: ok, Detail: detail}
}

// checkMomentumIdentity verifies that the recoil kinetic energy from
// momentum conservation matches Ke from the energy formula.
func checkMomentumIdentity(r Result) CheckResult {
	m := CheckMomentumConservation(r)
	detail := fmt.Sprintf("Ke from momentum = %.12g J, Ke from energy = %.12g J (rel diff %.3e)",
		m.KEFromMomentum, m.KEFromEnergy, m.RelativeDiscrepancy)
	ok := m.OK(1e-6)
	return CheckResult{Name: "momentum conservation", Pass: ok, Detail: detail}
}

// checkRelativeLossTrend verifies the cross rule that raising the
// incident energy at a fixed scattering angle increases the relative
// energy loss.
func checkRelativeLossTrend(in Input) CheckResult {
	low, err := Solve(Input{EnergyJ: in.EnergyJ, ThetaDeg: in.ThetaDeg})
	if err != nil {
		return CheckResult{Name: "relative loss grows with energy", Pass: false, Detail: err.Error()}
	}
	high, err := Solve(Input{EnergyJ: 10.0 * in.EnergyJ, ThetaDeg: in.ThetaDeg})
	if err != nil {
		return CheckResult{Name: "relative loss grows with energy", Pass: false, Detail: err.Error()}
	}
	cmp := CompareTwoEvents(low, high)
	detail := fmt.Sprintf("loss at E = %.4f, loss at 10*E = %.4f",
		low.LossFraction(), high.LossFraction())
	ok := cmp < 0
	return CheckResult{Name: "relative loss grows with energy", Pass: ok, Detail: detail}
}

// checkNinetyDegreeClosedForm verifies the closed form for a right
// angle scattering: E' = E/(1+E/me*c^2).
func checkNinetyDegreeClosedForm(r Result) CheckResult {
	ratio := constants.EnergyRatio(r.Input.EnergyJ)
	expected := r.Input.EnergyJ / (1.0 + ratio)
	diff := relErr(r.ScatteredEnergy, expected)
	detail := fmt.Sprintf("E' = %.12g J, closed form E/(1+E/me*c^2) = %.12g J (rel diff %.3e)",
		r.ScatteredEnergy, expected, diff)
	ok := diff <= DefaultTolerance
	return CheckResult{Name: "90 degree closed form", Pass: ok, Detail: detail}
}

// checkForwardNoShift verifies that theta = 0 leaves the photon
// untouched: dlambda = 0, E' = E and Ke = 0.
func checkForwardNoShift(r Result) CheckResult {
	shiftOK := relErr(r.DeltaLambda, 0) <= DefaultTolerance
	energyOK := relErr(r.ScatteredEnergy, r.Input.EnergyJ) <= DefaultTolerance
	keOK := relErr(r.RecoilEnergy, 0) <= DefaultTolerance
	detail := fmt.Sprintf("dlambda = %.4g m, E' = %.12g J, Ke = %.4g J",
		r.DeltaLambda, r.ScatteredEnergy, r.RecoilEnergy)
	ok := shiftOK && energyOK && keOK
	return CheckResult{Name: "forward scatter keeps photon unchanged", Pass: ok, Detail: detail}
}

// checkBackscatterMaximum verifies that theta = 180 degrees gives the
// maximum wavelength shift dlambda = 2*lambda_c and the minimum
// scattered energy.
func checkBackscatterMaximum(r Result) CheckResult {
	lambdaC := constants.ComptonWavelength()
	expected := 2.0 * lambdaC
	diff := relErr(r.DeltaLambda, expected)
	detail := fmt.Sprintf("dlambda = %.12g m, 2*lambda_c = %.12g m (rel diff %.3e)",
		r.DeltaLambda, expected, diff)
	ok := diff <= DefaultTolerance
	return CheckResult{Name: "backscatter shift equals 2*lambda_c", Pass: ok, Detail: detail}
}
