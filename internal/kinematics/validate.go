package kinematics

import (
	"errors"
	"fmt"
	"math"
)

// Validation errors shared by every entry point of the package. They
// are returned as plain errors whose message can be printed verbatim
// on stderr by the CLI layer.

// ErrNonPositiveEnergy is returned when the incident photon energy is
// zero or negative, both of which make the Compton formulas
// meaningless.
var ErrNonPositiveEnergy = errors.New("energy must be positive")

// ErrAngleOutOfRange is returned when the scattering angle falls
// outside [0, 180] degrees.
var ErrAngleOutOfRange = errors.New("scattering angle must be in [0, 180] degrees")

// ErrEnergyNotFinite is returned when the energy is NaN or infinite.
var ErrEnergyNotFinite = errors.New("energy must be a finite number")

// ErrAngleNotFinite is returned when the angle is NaN or infinite.
var ErrAngleNotFinite = errors.New("scattering angle must be a finite number")

// validate checks the physical range of an input event. A photon with
// non-positive energy and a scattering angle outside the forward and
// backward limits are rejected before any formula is evaluated.
func validate(in Input) error {
	if math.IsNaN(in.EnergyJ) || math.IsInf(in.EnergyJ, 0) {
		return ErrEnergyNotFinite
	}
	if in.EnergyJ <= 0 {
		return ErrNonPositiveEnergy
	}
	if math.IsNaN(in.ThetaDeg) || math.IsInf(in.ThetaDeg, 0) {
		return ErrAngleNotFinite
	}
	if in.ThetaDeg < 0 || in.ThetaDeg > 180 {
		return ErrAngleOutOfRange
	}
	return nil
}

// Validate exposes the input validation so callers (the CLI and the
// test suite) can reject an event before solving it.
func Validate(in Input) error {
	return validate(in)
}

// validateVerbose wraps validate and prefixes the message with the
// offending value so the user sees which field was rejected.
func validateVerbose(in Input) error {
	err := validate(in)
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, ErrNonPositiveEnergy):
		return fmt.Errorf("%w (got %g J)", err, in.EnergyJ)
	case errors.Is(err, ErrEnergyNotFinite):
		return fmt.Errorf("%w (got %g)", err, in.EnergyJ)
	case errors.Is(err, ErrAngleOutOfRange):
		return fmt.Errorf("%w (got %g deg)", err, in.ThetaDeg)
	case errors.Is(err, ErrAngleNotFinite):
		return fmt.Errorf("%w (got %g)", err, in.ThetaDeg)
	}
	return err
}

// ValidateVerbose exposes the prefixed form of the validation for the
// CLI layer to print.
func ValidateVerbose(in Input) error {
	return validateVerbose(in)
}

// DegreesToRadians converts a degree value into radians for the
// trigonometric functions. The conversion keeps no hidden rounding so
// that theta = 0, 90 and 180 map onto exact sine and cosine values.
func DegreesToRadians(deg float64) float64 {
	return deg * math.Pi / 180.0
}

// CosTheta returns cos(theta) for an angle given in degrees. The 0, 90
// and 180 degree cases land exactly on 1, 0 and -1.
func CosTheta(deg float64) float64 {
	return math.Cos(DegreesToRadians(deg))
}
