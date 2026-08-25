package kinematics

import (
	"compton-scat/internal/constants"
)

func RecoilFromEnergies(incidentJ, scatteredJ float64) float64 {
	return incidentJ - scatteredJ
}

func SumRecoilAndScattered(recoilJ, scatteredJ float64) float64 {
	return recoilJ + scatteredJ
}

func RelativeEnergyLoss(incidentJ, scatteredJ float64) float64 {
	if incidentJ == 0 {
		return 0
	}
	return (incidentJ - scatteredJ) / incidentJ
}

func EnergyRatioInKEV(energyKEV float64) float64 {
	return energyKEV / constants.ElectronRestEnergyKEV()
}

func ScatteredKEVFromRatio(incidentKEV, ratio float64) float64 {
	return incidentKEV * ratio
}

func CompareScatteredEnergy(r Result, expectedJ float64) float64 {
	return relErr(r.ScatteredEnergy, expectedJ)
}

func ScatteredRatioAtTheta(energyRatio, cosTheta float64) float64 {
	return 1.0 / (1.0 + energyRatio*(1.0-cosTheta))
}

func RatioForBackscatter(energyRatio float64) float64 {
	return 1.0 / (1.0 + 2.0*energyRatio)
}
