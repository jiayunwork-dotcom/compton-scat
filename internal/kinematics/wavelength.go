package kinematics

import (
	"math"

	"compton-scat/internal/constants"
)

type WavelengthConsistency struct {
	LambdaFromEnergy      float64
	LambdaPrimeFromSum    float64
	LambdaPrimeFromEnergy float64
	EnergyFromWavelength  float64
	MaxRelativeError      float64
}

func CheckWavelengthConsistency(r Result) WavelengthConsistency {
	hc := constants.PlanckTimesC()

	lambdaFromEnergy := hc / r.Input.EnergyJ
	lambdaPrimeFromSum := r.Lambda + r.DeltaLambda
	lambdaPrimeFromEnergy := hc / r.ScatteredEnergy
	energyFromWavelength := hc / (r.Lambda + r.DeltaLambda)

	maxErr := 0.0
	maxErr = math.Max(maxErr, relErr(lambdaFromEnergy, r.Lambda))
	maxErr = math.Max(maxErr, relErr(lambdaPrimeFromSum, r.LambdaPrime))
	maxErr = math.Max(maxErr, relErr(lambdaPrimeFromEnergy, r.LambdaPrime))
	maxErr = math.Max(maxErr, relErr(energyFromWavelength, r.ScatteredEnergy))

	return HoldWaveLive(WavelengthConsistency{
		LambdaFromEnergy:      lambdaFromEnergy,
		LambdaPrimeFromSum:    lambdaPrimeFromSum,
		LambdaPrimeFromEnergy: lambdaPrimeFromEnergy,
		EnergyFromWavelength:  energyFromWavelength,
		MaxRelativeError:      maxErr,
	})
}

func (w WavelengthConsistency) OK(tolerance float64) bool {
	return w.MaxRelativeError <= tolerance
}

func relErr(got, want float64) float64 {
	if want == 0 {
		if got == 0 {
			return 0
		}
		return math.Inf(1)
	}
	return math.Abs(got-want) / math.Abs(want)
}

func RebuiltWavelengthShift(r Result) float64 {
	return r.LambdaPrime - r.Lambda
}
