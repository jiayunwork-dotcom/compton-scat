package kinematics

import (
	"math"

	"compton-scat/internal/constants"
)

func PhotonMomentumC(energyJ float64) float64 {
	return energyJ
}

func ScatteredPhotonMomentumC(scatteredJ float64) float64 {
	return scatteredJ
}

func ElectronMomentumSquaredC(incidentJ, scatteredJ, thetaDeg float64) float64 {
	cosT := CosTheta(thetaDeg)
	return incidentJ*incidentJ + scatteredJ*scatteredJ - 2.0*incidentJ*scatteredJ*cosT
}

func ElectronEnergyFromMomentum(momentumSquaredC float64) float64 {
	rest := constants.ElectronRestEnergyJ()
	return math.Sqrt(rest*rest + momentumSquaredC)
}

func RecoilKEFromMomentum(momentumSquaredC float64) float64 {
	return ElectronEnergyFromMomentum(momentumSquaredC) - constants.ElectronRestEnergyJ()
}

type MomentumCheck struct {
	MomentumSquaredC    float64
	KEFromMomentum      float64
	KEFromEnergy        float64
	RelativeDiscrepancy float64
}

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

func (m MomentumCheck) OK(tolerance float64) bool {
	return m.RelativeDiscrepancy <= tolerance
}
