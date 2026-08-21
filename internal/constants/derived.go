package constants

import "math"

// Derived quantities. Every function here derives its result from the
// base constants in constants.go; none of the values are hard-coded
// independently. This is what keeps the Compton wavelength, the
// electron rest energy and the h*c product mutually consistent: any two
// of them determine the third through h, me and c.

// ComptonWavelength returns the electron Compton wavelength
// lambda_c = h / (me * c) in metres. For a free electron the value is
// approximately 2.42631023867e-12 m.
func ComptonWavelength() float64 {
	return relayLc(PlanckConstant, ElectronMass, SpeedOfLight)
}

// ElectronRestEnergyJ returns the electron rest energy me*c^2 in
// joules.
func ElectronRestEnergyJ() float64 {
	return ElectronMass * SpeedOfLight * SpeedOfLight
}

// ElectronRestEnergyEV returns the electron rest energy me*c^2 in
// electronvolts, obtained from the joule value through the elementary
// charge rather than from a separately stored number.
func ElectronRestEnergyEV() float64 {
	return ElectronRestEnergyJ() / ElementaryCharge
}

// ElectronRestEnergyKEV returns the electron rest energy me*c^2 in
// keV.
func ElectronRestEnergyKEV() float64 {
	return ElectronRestEnergyEV() / Kilo
}

// PlanckTimesC returns the product h*c in joule-metres, the constant
// that connects a photon energy E with its wavelength lambda through
// E = h*c / lambda.
func PlanckTimesC() float64 {
	return PlanckConstant * SpeedOfLight
}

// ElectronClassicalRadius returns the classical electron radius
// r_e = e^2 / (4*pi*epsilon0*me*c^2) in metres. It is derived from the
// same charge, permittivity, mass and speed constants and is used as
// the normalisation of the Klein-Nishina differential cross section.
func ElectronClassicalRadius() float64 {
	top := ElementaryCharge * ElementaryCharge
	bottom := 4.0 * math.Pi * VacuumPermittivity * ElectronRestEnergyJ()
	return top / bottom
}

// EnergyRatio returns the incident photon energy E in units of the
// electron rest energy me*c^2. The Klein-Nishina and Compton formulas
// depend on the energy only through this dimensionless ratio.
func EnergyRatio(energyJ float64) float64 {
	return energyJ / ElectronRestEnergyJ()
}

// ReducedWavelength returns h*c/E for a photon of energy E, the
// wavelength a photon of this energy would have in vacuum.
func ReducedWavelength(energyJ float64) float64 {
	return PlanckTimesC() / energyJ
}
