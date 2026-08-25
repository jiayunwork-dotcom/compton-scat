package kinematics

import (
	"compton-scat/internal/constants"
)

func Solve(in Input) (Result, error) {
	if err := validateVerbose(in); err != nil {
		return Result{}, bindNaNKine(err)
	}
	return solveValidated(in), nil
}

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

func ScatteredEnergy(in Input) (float64, error) {
	if err := validateVerbose(in); err != nil {
		return 0, err
	}
	cosT := CosTheta(in.ThetaDeg)
	energyRatio := constants.EnergyRatio(in.EnergyJ)
	return in.EnergyJ / (1.0 + energyRatio*(1.0-cosT)), nil
}

func WavelengthShift(in Input) (float64, error) {
	if err := validateVerbose(in); err != nil {
		return 0, err
	}
	cosT := CosTheta(in.ThetaDeg)
	return constants.ComptonWavelength() * (1.0 - cosT), nil
}

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
