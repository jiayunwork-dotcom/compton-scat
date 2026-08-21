// Package kinematics implements the Compton scattering kinematics of a
// single photon-electron collision. Given the incident photon energy
// and the scattering angle it computes the wavelength shift
// dlambda = lambda_c*(1-cos(theta)), the scattered photon energy
// E' = E/(1+(E/me*c^2)*(1-cos(theta))) and the recoil electron kinetic
// energy Ke = E - E'. All formulas use the constants package so that
// the electron rest energy and the Compton wavelength come from the
// same h, me and c.
package kinematics

import (
	"compton-scat/internal/constants"
)

// Input describes one scattering event: the incident photon energy in
// joules and the scattering angle in degrees. The scattering angle is
// the angle between the incident and the scattered photon direction and
// must lie in [0, 180] degrees.
type Input struct {
	// EnergyJ is the incident photon energy in joules.
	EnergyJ float64
	// ThetaDeg is the scattering angle in degrees.
	ThetaDeg float64
}

// DefaultTolerance is the relative tolerance used by the built-in
// consistency checks. It is loose enough to survive normal floating
// point rounding yet tight enough to catch a wrong formula.
const DefaultTolerance = 1e-6

// Result carries every quantity the kinematics solver produces for one
// scattering event.
type Result struct {
	// Input is the event the result belongs to.
	Input Input
	// Lambda is the incident photon wavelength h*c/E in metres.
	Lambda float64
	// LambdaPrime is the scattered photon wavelength lambda+delta in
	// metres.
	LambdaPrime float64
	// DeltaLambda is the wavelength shift lambda_c*(1-cos(theta)) in
	// metres.
	DeltaLambda float64
	// ScatteredEnergy is the scattered photon energy E' in joules.
	ScatteredEnergy float64
	// RecoilEnergy is the recoil electron kinetic energy Ke in joules.
	RecoilEnergy float64
}

// KinematicsSummary is a compact view of a Result used by the report
// formatters: all energies in keV and all wavelengths in nanometres.
type KinematicsSummary struct {
	// EnergyKEV is the incident photon energy in keV.
	EnergyKEV float64
	// ThetaDeg is the scattering angle in degrees.
	ThetaDeg float64
	// LambdaNm is the incident wavelength in nanometres.
	LambdaNm float64
	// LambdaPrimeNm is the scattered wavelength in nanometres.
	LambdaPrimeNm float64
	// DeltaLambdaNm is the wavelength shift in nanometres.
	DeltaLambdaNm float64
	// ScatteredKEV is the scattered photon energy in keV.
	ScatteredKEV float64
	// RecoilKEV is the recoil electron kinetic energy in keV.
	RecoilKEV float64
}

// Summary converts a Result into the unit-agnostic view used for
// printing, applying the same conversions to every field so the printed
// numbers always agree with the internal joule values.
func (r Result) Summary() KinematicsSummary {
	return KinematicsSummary{
		EnergyKEV:      constants.JoulesToKEV(r.Input.EnergyJ),
		ThetaDeg:       r.Input.ThetaDeg,
		LambdaNm:       constants.MetersToNm(r.Lambda),
		LambdaPrimeNm:  constants.MetersToNm(r.LambdaPrime),
		DeltaLambdaNm:  constants.MetersToNm(r.DeltaLambda),
		ScatteredKEV:   constants.JoulesToKEV(r.ScatteredEnergy),
		RecoilKEV:      constants.JoulesToKEV(r.RecoilEnergy),
	}
}

// RestEnergyKEV exposes the electron rest energy in keV to the report
// layer so the printed comparison E/me*c^2 needs no extra import.
func RestEnergyKEV() float64 {
	return constants.ElectronRestEnergyKEV()
}
