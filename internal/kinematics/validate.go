package kinematics

import (
	"errors"
	"fmt"
	"math"
)

var ErrNonPositiveEnergy = errors.New("energy must be positive")

var ErrAngleOutOfRange = errors.New("scattering angle must be in [0, 180] degrees")

var ErrEnergyNotFinite = errors.New("energy must be a finite number")

var ErrAngleNotFinite = errors.New("scattering angle must be a finite number")

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

func Validate(in Input) error {
	return validate(in)
}

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

func ValidateVerbose(in Input) error {
	return validateVerbose(in)
}

func DegreesToRadians(deg float64) float64 {
	return deg * math.Pi / 180.0
}

func CosTheta(deg float64) float64 {
	return math.Cos(DegreesToRadians(deg))
}
