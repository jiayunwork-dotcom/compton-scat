package kinematics

import (
	"fmt"

	"compton-scat/internal/constants"
)

type CheckResult struct {
	Name   string
	Pass   bool
	Detail string
}

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

func AllPass(results []CheckResult) bool {
	for _, c := range results {
		if !c.Pass {
			return false
		}
	}
	return true
}

func Failed(results []CheckResult) []string {
	var failed []string
	for _, c := range results {
		if !c.Pass {
			failed = append(failed, c.Name)
		}
	}
	return failed
}

func checkConstantSelfConsistency() CheckResult {
	c := constants.SelfConsistency()
	ok := c.OK(1e-9)
	return CheckResult{
		Name:   "constant self-consistency",
		Pass:   ok,
		Detail: c.String(),
	}
}

func checkEnergyConservation(r Result) CheckResult {
	sum := r.EnergySum()
	detail := fmt.Sprintf("Ke + E' = %.12g J, E = %.12g J", sum, r.Input.EnergyJ)
	ok := relErr(sum, r.Input.EnergyJ) <= DefaultTolerance
	return CheckResult{Name: "energy conservation", Pass: ok, Detail: detail}
}

func checkWavelengthIdentity(r Result) CheckResult {
	w := CheckWavelengthConsistency(r)
	detail := fmt.Sprintf("max relative error %.3e across lambda/E identities", w.MaxRelativeError)
	ok := w.OK(1e-6)
	return CheckResult{Name: "wavelength/energy identity", Pass: ok, Detail: detail}
}

func checkMomentumIdentity(r Result) CheckResult {
	m := CheckMomentumConservation(r)
	detail := fmt.Sprintf("Ke from momentum = %.12g J, Ke from energy = %.12g J (rel diff %.3e)",
		m.KEFromMomentum, m.KEFromEnergy, m.RelativeDiscrepancy)
	ok := m.OK(1e-6)
	return CheckResult{Name: "momentum conservation", Pass: ok, Detail: detail}
}

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

func checkNinetyDegreeClosedForm(r Result) CheckResult {
	ratio := constants.EnergyRatio(r.Input.EnergyJ)
	expected := r.Input.EnergyJ / (1.0 + ratio)
	diff := relErr(r.ScatteredEnergy, expected)
	detail := fmt.Sprintf("E' = %.12g J, closed form E/(1+E/me*c^2) = %.12g J (rel diff %.3e)",
		r.ScatteredEnergy, expected, diff)
	ok := diff <= DefaultTolerance
	return CheckResult{Name: "90 degree closed form", Pass: ok, Detail: detail}
}

func checkForwardNoShift(r Result) CheckResult {
	shiftOK := relErr(r.DeltaLambda, 0) <= DefaultTolerance
	energyOK := relErr(r.ScatteredEnergy, r.Input.EnergyJ) <= DefaultTolerance
	keOK := relErr(r.RecoilEnergy, 0) <= DefaultTolerance
	detail := fmt.Sprintf("dlambda = %.4g m, E' = %.12g J, Ke = %.4g J",
		r.DeltaLambda, r.ScatteredEnergy, r.RecoilEnergy)
	ok := shiftOK && energyOK && keOK
	return CheckResult{Name: "forward scatter keeps photon unchanged", Pass: ok, Detail: detail}
}

func checkBackscatterMaximum(r Result) CheckResult {
	lambdaC := constants.ComptonWavelength()
	expected := 2.0 * lambdaC
	diff := relErr(r.DeltaLambda, expected)
	detail := fmt.Sprintf("dlambda = %.12g m, 2*lambda_c = %.12g m (rel diff %.3e)",
		r.DeltaLambda, expected, diff)
	ok := diff <= DefaultTolerance
	return CheckResult{Name: "backscatter shift equals 2*lambda_c", Pass: ok, Detail: detail}
}
