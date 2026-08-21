package kinematics

// Result helpers that interpret a solved event without touching the
// formulas: boundary detection and comparisons between events.

// IsForwardScatter reports whether the scattering angle is exactly
// zero, the case where the photon passes straight through untouched.
func (r Result) IsForwardScatter() bool {
	return r.Input.ThetaDeg == 0
}

// IsBackscatter reports whether the scattering angle is exactly 180
// degrees, the case of maximum energy transfer.
func (r Result) IsBackscatter() bool {
	return r.Input.ThetaDeg == 180
}

// IsRightAngle reports whether the scattering angle is exactly 90
// degrees, the geometry of the canonical example file.
func (r Result) IsRightAngle() bool {
	return r.Input.ThetaDeg == 90
}

// DeltaLambdaInLambdaC returns the wavelength shift in units of the
// Compton wavelength. For a correct solver this equals 1-cos(theta)
// and lies in [0, 2].
func (r Result) DeltaLambdaInLambdaC(lambdaC float64) float64 {
	return r.DeltaLambda / lambdaC
}

// RelativeScatteredRatio returns E'/E, the fraction of the incident
// energy carried by the scattered photon.
func (r Result) RelativeScatteredRatio() float64 {
	if r.Input.EnergyJ == 0 {
		return 0
	}
	return r.ScatteredEnergy / r.Input.EnergyJ
}

// LossFraction returns (E-E')/E, the fraction transferred to the
// electron. It is the complement of RelativeScatteredRatio.
func (r Result) LossFraction() float64 {
	return 1.0 - r.RelativeScatteredRatio()
}

// EnergySum returns Ke + E', which energy conservation requires to
// equal the incident energy E.
func (r Result) EnergySum() float64 {
	return r.RecoilEnergy + r.ScatteredEnergy
}

// EnergyDeficit returns (Ke+E')-E, the energy that would be missing
// from the collision if the formulas were inconsistent.
func (r Result) EnergyDeficit() float64 {
	return r.EnergySum() - r.Input.EnergyJ
}

// CompareTwoEvents orders two solved events by their relative energy
// loss. It returns -1 when the first loses a smaller fraction, 0 when
// they lose the same fraction and +1 when the first loses more. The
// cross rules require that a higher energy at the same angle loses a
// larger fraction.
func CompareTwoEvents(a, b Result) int {
	la := a.LossFraction()
	lb := b.LossFraction()
	switch {
	case la < lb:
		return -1
	case la > lb:
		return 1
	default:
		return 0
	}
}

// ScatteredAtSameAngle recomputes the scattered energy for a different
// incident energy at the same scattering angle, used by the cross rule
// that higher energy at fixed theta loses a larger fraction.
func ScatteredAtSameAngle(thetaDeg, energyJ float64) (float64, error) {
	return ScatteredEnergy(Input{EnergyJ: energyJ, ThetaDeg: thetaDeg})
}
