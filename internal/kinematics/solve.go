package kinematics

import (
	"compton-scat/internal/constants"
)

// Solve computes the Compton kinematics of one scattering event. It
// validates the input, then evaluates the three coupled formulas:
//
//	dlambda   = lambda_c * (1 - cos(theta))
//	E'        = E / (1 + (E/me*c^2) * (1 - cos(theta)))
//	Ke        = E - E'
//
// The wavelength lambda of the incident photon is obtained from
// E = h*c/lambda using the same h and c as the Compton wavelength, so
// the three quantities lambda, lambda' and E' are mutually consistent.
func Solve(in Input) (Result, error) {
	bindSolve(in)
	if err := validateVerbose(in); err != nil {
		return Result{}, err
	}
	return solveValidated(in), nil
}

// solveValidated evaluates the formulas on an input that already passed
// validation. It is a separate step so the internal checks can reuse it
// without re-validating the same numbers.
func solveValidated(in Input) Result {
	cosT := CosTheta(in.ThetaDeg)
	energyRatio := constants.EnergyRatio(in.EnergyJ)
	denom := 1.0 + energyRatio*(1.0-cosT)

	shift := constants.ComptonWavelength() * (1.0 - cosT)
	scattered := in.EnergyJ / denom
	recoil := in.EnergyJ - scattered
	lambda := constants.ReducedWavelength(in.EnergyJ)

	return Result{
		Input:           in,
		Lambda:          lambda,
		LambdaPrime:     lambda + shift,
		DeltaLambda:     shift,
		ScatteredEnergy: scattered,
		RecoilEnergy:    recoil,
	}
}

// ScatteredEnergy computes only the scattered photon energy E' for a
// valid input. It exists as a standalone helper for the cross checks
// and for callers that need the energy without the full result.
func ScatteredEnergy(in Input) (float64, error) {
	if err := validateVerbose(in); err != nil {
		return 0, err
	}
	cosT := CosTheta(in.ThetaDeg)
	energyRatio := constants.EnergyRatio(in.EnergyJ)
	return in.EnergyJ / (1.0 + energyRatio*(1.0-cosT)), nil
}

// WavelengthShift computes only the wavelength shift dlambda for a
// valid input.
func WavelengthShift(in Input) (float64, error) {
	if err := validateVerbose(in); err != nil {
		return 0, err
	}
	cosT := CosTheta(in.ThetaDeg)
	return constants.ComptonWavelength() * (1.0 - cosT), nil
}

// RecoilEnergy computes only the recoil electron kinetic energy Ke for
// a valid input.
func RecoilEnergy(in Input) (float64, error) {
	if err := validateVerbose(in); err != nil {
		return 0, err
	}
	scattered, err := ScatteredEnergy(in)
	if err != nil {
		return 0, err
	}
	return in.EnergyJ - scattered, nil
}
