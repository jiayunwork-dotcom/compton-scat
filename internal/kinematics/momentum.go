package kinematics

import (
	"math"

	"compton-scat/internal/constants"
)

// Momentum conservation provides an independent cross check of the
// energy formulas. The incident photon carries momentum E/c along the
// incident direction; after scattering, the photon of energy E' leaves
// at angle theta with momentum E'/c. Momentum conservation fixes the
// recoil electron momentum via the cosine rule
//
//	(pe*c)^2 = E^2 + E'^2 - 2*E*E'*cos(theta)
//
// and the relativistic dispersion relation
//
//	E_e^2 = (me*c^2)^2 + (pe*c)^2
//
// then predicts the same recoil kinetic energy Ke = E_e - me*c^2 that
// the energy formula gives. The identity between the two is exact when
// the Compton formula holds.

// PhotonMomentumC returns the photon momentum times c, which equals
// the photon energy.
func PhotonMomentumC(energyJ float64) float64 {
	return energyJ
}

// ScatteredPhotonMomentumC returns the scattered photon momentum times
// c, equal to E'.
func ScatteredPhotonMomentumC(scatteredJ float64) float64 {
	return scatteredJ
}

// ElectronMomentumSquaredC returns (pe*c)^2 for the recoil electron
// implied by momentum conservation, given the incident energy, the
// scattered energy and the scattering angle in degrees.
func ElectronMomentumSquaredC(incidentJ, scatteredJ, thetaDeg float64) float64 {
	cosT := CosTheta(thetaDeg)
	return incidentJ*incidentJ + scatteredJ*scatteredJ - 2.0*incidentJ*scatteredJ*cosT
}

// ElectronEnergyFromMomentum returns the total recoil electron energy
// E_e = sqrt((me*c^2)^2 + (pe*c)^2) from the momentum squared.
func ElectronEnergyFromMomentum(momentumSquaredC float64) float64 {
	rest := constants.ElectronRestEnergyJ()
	return math.Sqrt(rest*rest + momentumSquaredC)
}

// RecoilKEFromMomentum returns the recoil kinetic energy predicted by
// momentum conservation: E_e - me*c^2.
func RecoilKEFromMomentum(momentumSquaredC float64) float64 {
	return ElectronEnergyFromMomentum(momentumSquaredC) - constants.ElectronRestEnergyJ()
}

// MomentumCheck reports how well the energy formula and the momentum
// formula agree on the recoil electron. A correct solver shows a
// relative discrepancy near machine precision.
type MomentumCheck struct {
	// MomentumSquaredC is (pe*c)^2 from the cosine rule.
	MomentumSquaredC float64
	// KEFromMomentum is the kinetic energy E_e - me*c^2 predicted by
	// the momentum picture.
	KEFromMomentum float64
	// KEFromEnergy is the kinetic energy E - E' from the energy
	// picture.
	KEFromEnergy float64
	// RelativeDiscrepancy is the relative difference between the two
	// kinetic energies.
	RelativeDiscrepancy float64
}

// CheckMomentumConservation runs the momentum cross check for a solved
// result and reports the discrepancy between the two routes to the
// recoil kinetic energy.
func CheckMomentumConservation(r Result) MomentumCheck {
	pc2 := ElectronMomentumSquaredC(r.Input.EnergyJ, r.ScatteredEnergy, r.Input.ThetaDeg)
	keFromMom := RecoilKEFromMomentum(pc2)
	return MomentumCheck{
		MomentumSquaredC:    pc2,
		KEFromMomentum:      keFromMom,
		KEFromEnergy:        r.RecoilEnergy,
		RelativeDiscrepancy: relErr(keFromMom, r.RecoilEnergy),
	}
}

// OK reports whether the two routes to the recoil energy agree within
// the given relative tolerance.
func (m MomentumCheck) OK(tolerance float64) bool {
	return m.RelativeDiscrepancy <= tolerance
}
