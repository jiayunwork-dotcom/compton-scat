// Package section implements the Klein-Nishina cross section for
// Compton scattering off a free electron at rest. For a photon of
// energy E the differential cross section per unit solid angle is
//
//	dsigma/dOmega = (r_e^2/2) * P^2 * (1/P + P - sin^2(theta))
//
// with P = lambda'/lambda = 1/(1+(E/me*c^2)*(1-cos(theta))). The total
// cross section is obtained by numerical integration over the full
// sphere and, as required by the Klein-Nishina formula, decreases as
// the incident energy rises. The low-energy limit reproduces the
// Thomson cross section 8*pi*r_e^2/3.
package section

// DefaultIntegrationNodes is the number of panels used by the
// composite Simpson rule for the total cross section. An even number of
// panels gives a closed Simpson form.
const DefaultIntegrationNodes = 2000

// Result carries the differential and total cross sections for one
// incident energy.
type Result struct {
	// EnergyJ is the incident photon energy in joules.
	EnergyJ float64
	// DifferentialAt90 is dsigma/dOmega at theta = 90 degrees in
	// m^2/sr.
	DifferentialAt90 float64
	// DifferentialAt180 is dsigma/dOmega at theta = 180 degrees in
	// m^2/sr.
	DifferentialAt180 float64
	// Total is the total cross section in m^2.
	Total float64
	// Thomson is the low-energy limit 8*pi*r_e^2/3 in m^2.
	Thomson float64
	// TotalRatio is Total/Thomson, below one for every finite energy.
	TotalRatio float64
}

// CrossSection solves the Klein-Nishina quantities for one incident
// energy in joules.
func CrossSection(energyJ float64) (Result, error) {
	if energyJ <= 0 {
		return Result{}, errNonPositiveEnergy()
	}
	return CrossSectionValidated(energyJ), nil
}

// CrossSectionValidated computes the section for an energy that has
// already been checked to be positive.
func CrossSectionValidated(energyJ float64) Result {
	ds90 := Differential(energyJ, 90)
	ds180 := Differential(energyJ, 180)
	total := Total(energyJ, DefaultIntegrationNodes)
	thomson := ThomsonCrossSection()
	return Result{
		EnergyJ:          energyJ,
		DifferentialAt90: ds90,
		DifferentialAt180: ds180,
		Total:            total,
		Thomson:          thomson,
		TotalRatio:       total / thomson,
	}
}
