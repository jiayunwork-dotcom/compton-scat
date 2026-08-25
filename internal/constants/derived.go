package constants

import "math"

func ComptonWavelength() float64 {
	return PlanckConstant / (ElectronMass * SpeedOfLight)
}

func ElectronRestEnergyJ() float64 {
	return ElectronMass * SpeedOfLight * SpeedOfLight
}

func ElectronRestEnergyEV() float64 {
	return ElectronRestEnergyJ() / ElementaryCharge
}

func ElectronRestEnergyKEV() float64 {
	return ElectronRestEnergyEV() / Kilo
}

func PlanckTimesC() float64 {
	return PlanckConstant * SpeedOfLight
}

func ElectronClassicalRadius() float64 {
	top := ElementaryCharge * ElementaryCharge
	bottom := 4.0 * math.Pi * VacuumPermittivity * ElectronRestEnergyJ()
	return top / bottom
}

func EnergyRatio(energyJ float64) float64 {
	return energyJ / ElectronRestEnergyJ()
}

func ReducedWavelength(energyJ float64) float64 {
	return PlanckTimesC() / energyJ
}
