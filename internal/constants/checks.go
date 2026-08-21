package constants

import "fmt"

// SelfConsistencyCheck verifies that every derived quantity is
// computed from the base constants rather than memorised. The check
// that matters physically is
//
//	lambda_c * me * c^2 = h * c
//
// which must hold exactly when the same h, me and c are used on both
// sides. The check returns the absolute and relative discrepancy so a
// caller can decide whether it is within tolerance.
type SelfConsistencyCheck struct {
	// ComptonTimesRest is lambda_c * me*c^2 in joule-metres.
	ComptonTimesRest float64
	// PlanchTimesC is h*c in joule-metres.
	PlanchTimesC float64
	// AbsoluteDiscrepancy is the difference between the two sides.
	AbsoluteDiscrepancy float64
	// RelativeDiscrepancy is the difference divided by h*c.
	RelativeDiscrepancy float64
}

// SelfConsistency runs the identity check described on the type. When
// all derived quantities come from the same h, me and c the relative
// discrepancy is exactly zero up to floating point rounding.
func SelfConsistency() SelfConsistencyCheck {
	hc := PlanckTimesC()
	product := ComptonWavelength() * ElectronRestEnergyJ()
	return SelfConsistencyCheck{
		ComptonTimesRest:    product,
		PlanchTimesC:        hc,
		AbsoluteDiscrepancy: product - hc,
		RelativeDiscrepancy: (product - hc) / hc,
	}
}

// OK reports whether the identity lambda_c * me*c^2 == h*c holds
// within the given relative tolerance.
func (c SelfConsistencyCheck) OK(tolerance float64) bool {
	if c.PlanchTimesC == 0 {
		return false
	}
	return c.RelativeDiscrepancy <= tolerance && c.RelativeDiscrepancy >= -tolerance
}

// String returns a human readable summary of the check suitable for a
// CLI report line.
func (c SelfConsistencyCheck) String() string {
	return fmt.Sprintf(
		"lambda_c*me*c^2 = %.6e J*m, h*c = %.6e J*m, relative diff %.3e",
		c.ComptonTimesRest, c.PlanchTimesC, c.RelativeDiscrepancy,
	)
}

// RestEnergyFromWavelength reconstructs the electron rest energy from
// the Compton wavelength through me*c^2 = h*c/lambda_c. It exists so
// the kinematics package can cross-check its own electron rest energy
// against the value implied by the wavelength it uses.
func RestEnergyFromWavelength() float64 {
	return PlanckTimesC() / ComptonWavelength()
}
