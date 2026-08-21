package section

import (
	"math"

	"compton-scat/internal/constants"
)

// LowEnergy behaviour of the Klein-Nishina total cross section. As the
// incident energy tends to zero the total cross section approaches the
// Thomson value 8*pi*r_e^2/3 from below; the relative gap is of order
// k = E/me*c^2.

// TotalAtRatio returns the total cross section for an energy expressed
// in units of me*c^2, i.e. for k = E/me*c^2.
func TotalAtRatio(k float64, panels int) float64 {
	energyJ := k * RestEnergyJoules()
	return Total(energyJ, panels)
}

// RestEnergyJoules returns the electron rest energy in joules, taken
// from the constants package so the section formulas share the same
// me, c and h as the kinematics formulas.
func RestEnergyJoules() float64 {
	return constants.ElectronRestEnergyJ()
}

// LowEnergyCorrection returns 1 - sigma/sigma_Thomson for the given
// dimensionless ratio k. For small k this grows linearly in k; for
// large k it approaches 1.
func LowEnergyCorrection(k float64, panels int) float64 {
	if k <= 0 {
		return 0
	}
	total := TotalAtRatio(k, panels)
	return 1.0 - total/ThomsonCrossSection()
}

// RelativeToThomson returns sigma/sigma_Thomson, which is 1 at zero
// energy and falls towards zero at high energy.
func RelativeToThomson(energyJ float64, panels int) float64 {
	return Total(energyJ, panels) / ThomsonCrossSection()
}

// IsBelowThomson reports whether the total cross section for a
// positive energy is strictly below the Thomson limit.
func IsBelowThomson(energyJ float64, panels int) bool {
	return Total(energyJ, panels) < ThomsonCrossSection()
}

// MonotoneDecreasing reports whether the total cross section decreases
// across the given sorted energies, using the stricter adjacent check
// sigma(E_i+1) < sigma(E_i).
func MonotoneDecreasing(energiesJ []float64, panels int) bool {
	for i := 1; i < len(energiesJ); i++ {
		if Total(energiesJ[i], panels) >= Total(energiesJ[i-1], panels) {
			return false
		}
	}
	return true
}

// logSpace returns n energies spread logarithmically from loJ to hiJ
// inclusive, the set used by the trend command to probe the fall of the
// total cross section.
func logSpace(loJ, hiJ float64, n int) []float64 {
	values := make([]float64, n)
	lo := math.Log(loJ)
	hi := math.Log(hiJ)
	for i := 0; i < n; i++ {
		t := float64(i) / float64(n-1)
		values[i] = math.Exp(lo + t*(hi-lo))
	}
	return values
}

// TrendEnergies returns the default logarithmic energy ladder in keV
// used by the trend command.
func TrendEnergies() []float64 {
	lo := 1.0
	hi := 10000.0
	const count = 9
	kev := make([]float64, count)
	loLog := math.Log(lo)
	hiLog := math.Log(hi)
	for i := 0; i < count; i++ {
		t := float64(i) / float64(count-1)
		kev[i] = math.Exp(loLog + t*(hiLog-lo))
	}
	return kev
}
