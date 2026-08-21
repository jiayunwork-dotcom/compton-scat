package kinematics

import (
	"math"

	"compton-scat/internal/constants"
)

// WavelengthConsistency cross-checks the wavelength side of the
// Compton relations against the energy side. For any valid event the
// three identities
//
//	lambda   = h*c / E
//	lambda'  = lambda + dlambda
//	E'       = h*c / lambda'
//
// must hold simultaneously. The check reports the largest relative
// discrepancy found among them.
type WavelengthConsistency struct {
	// LambdaFromEnergy is h*c/E for the incident photon.
	LambdaFromEnergy float64
	// LambdaPrimeFromSum is lambda + dlambda.
	LambdaPrimeFromSum float64
	// LambdaPrimeFromEnergy is h*c/E'.
	LambdaPrimeFromEnergy float64
	// EnergyFromWavelength is h*c/lambda' rebuilt from the scattered
	// wavelength.
	EnergyFromWavelength float64
	// MaxRelativeError is the largest relative discrepancy among the
	// identities above.
	MaxRelativeError float64
}

// CheckWavelengthConsistency runs the three identities above for a
// solved result. A correct solver returns a MaxRelativeError below
// about 1e-12; a result whose formulas no longer agree (for example a
// denominator that has been changed) shows a large error or a negative
// energy that the check reports.
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

	return WavelengthConsistency{
		LambdaFromEnergy:      lambdaFromEnergy,
		LambdaPrimeFromSum:    lambdaPrimeFromSum,
		LambdaPrimeFromEnergy: lambdaPrimeFromEnergy,
		EnergyFromWavelength:  energyFromWavelength,
		MaxRelativeError:      maxErr,
	}
}

// OK reports whether every identity holds within the given relative
// tolerance.
func (w WavelengthConsistency) OK(tolerance float64) bool {
	return w.MaxRelativeError <= tolerance
}

// relErr returns the relative difference between got and want, defined
// as |got-want|/|want|. A zero want value yields a large error so that
// a vanished energy is always flagged.
func relErr(got, want float64) float64 {
	if want == 0 {
		if got == 0 {
			return 0
		}
		return math.Inf(1)
	}
	return math.Abs(got-want) / math.Abs(want)
}

// RebuiltWavelengthShift returns dlambda as the difference between the
// scattered and the incident wavelength, the direct measurement a
// spectroscopist would make. It must agree with lambda_c*(1-cos(theta)).
func RebuiltWavelengthShift(r Result) float64 {
	return r.LambdaPrime - r.Lambda
}
