package kinematics

import (
	"compton-scat/internal/constants"
)

// Energy relations derived from the scattered photon energy. They are
// separated from the solver so the cross checks can verify each
// conservation law on its own.

// RecoilFromEnergies computes the recoil electron kinetic energy from
// the incident and scattered photon energy. Conservation of energy in
// the collision reads E + me*c^2 = E' + E_e with E_e = Ke + me*c^2,
// which reduces to Ke = E - E'.
func RecoilFromEnergies(incidentJ, scatteredJ float64) float64 {
	return incidentJ - scatteredJ
}

// SumRecoilAndScattered returns Ke + E'. Energy conservation requires
// this to equal the incident photon energy E.
func SumRecoilAndScattered(recoilJ, scatteredJ float64) float64 {
	return recoilJ + scatteredJ
}

// RelativeEnergyLoss returns the fraction (E - E')/E of the incident
// energy transferred to the recoil electron. For a fixed scattering
// angle this ratio grows as the incident energy rises.
func RelativeEnergyLoss(incidentJ, scatteredJ float64) float64 {
	if incidentJ == 0 {
		return 0
	}
	return (incidentJ - scatteredJ) / incidentJ
}

// EnergyRatioInKEV returns the dimensionless ratio E/me*c^2 using keV
// values. It is the same number as the joule-based EnergyRatio from the
// constants package because both units cancel.
func EnergyRatioInKEV(energyKEV float64) float64 {
	return energyKEV / constants.ElectronRestEnergyKEV()
}

// ScatteredKEVFromRatio converts a scattered energy ratio back into
// keV for reporting: given the ratio E'/E it returns E' in keV.
func ScatteredKEVFromRatio(incidentKEV, ratio float64) float64 {
	return incidentKEV * ratio
}

// CompareScatteredEnergy checks the scattered energy against an
// independently computed expectation and returns the relative
// discrepancy. The expectation is normally E/(1+(E/me*c^2)(1-cos(theta)))
// evaluated from scratch, so this comparison catches a solver that has
// drifted from the closed form.
func CompareScatteredEnergy(r Result, expectedJ float64) float64 {
	return relErr(r.ScatteredEnergy, expectedJ)
}

// ScatteredRatioAtTheta returns E'/E for a scattering angle theta. The
// ratio lies in (0, 1] for every valid theta and is independent of the
// absolute energy scale.
func ScatteredRatioAtTheta(energyRatio, cosTheta float64) float64 {
	return 1.0 / (1.0 + energyRatio*(1.0-cosTheta))
}

// RatioForBackscatter returns E'/E at theta = 180 degrees, the
// smallest scattered ratio for a given incident energy.
func RatioForBackscatter(energyRatio float64) float64 {
	return 1.0 / (1.0 + 2.0*energyRatio)
}
