package section

import (
	"compton-scat/internal/constants"
)

// Trend describes the behaviour of the total cross section across an
// energy ladder: the value at each rung and whether the sequence is
// strictly decreasing.
type Trend struct {
	// EnergiesKEV is the ladder of incident energies in keV.
	EnergiesKEV []float64
	// Totals is the total cross section at each rung in m^2.
	Totals []float64
	// Ratios is the total cross section relative to the Thomson value
	// at each rung.
	Ratios []float64
	// Decreasing reports whether the totals strictly decrease along
	// the ladder.
	Decreasing bool
}

// ComputeTrend evaluates the total cross section on the default
// logarithmic ladder from 1 keV to 10 MeV.
func ComputeTrend(panels int) Trend {
	kev := TrendEnergies()
	energies := make([]float64, len(kev))
	totals := make([]float64, len(kev))
	ratios := make([]float64, len(kev))
	for i, e := range kev {
		energyJ := constants.KEVToJoules(e)
		energies[i] = e
		totals[i] = Total(energyJ, panels)
		ratios[i] = totals[i] / ThomsonCrossSection()
	}
	dec := true
	for i := 1; i < len(totals); i++ {
		if totals[i] >= totals[i-1] {
			dec = false
			break
		}
	}
	return Trend{
		EnergiesKEV: energies,
		Totals:      totals,
		Ratios:      ratios,
		Decreasing:  dec,
	}
}

// VerifyTrend runs the monotonicity property over a wider, denser
// ladder than the printed one and reports whether it held, plus the
// first offending pair when it did not.
func VerifyTrend(panels int) (ok bool, firstFailAt, firstFailNext int) {
	kev := TrendEnergies()
	for i := 0; i+1 < len(kev); i++ {
		lo := constants.KEVToJoules(kev[i])
		hi := constants.KEVToJoules(kev[i+1])
		if Total(hi, panels) >= Total(lo, panels) {
			return false, i, i + 1
		}
	}
	return true, -1, -1
}

// MedianTotal returns the total cross section at the middle rung of the
// ladder, a stable reference point for the reports.
func MedianTotal(panels int) float64 {
	kev := TrendEnergies()
	mid := len(kev) / 2
	return Total(constants.KEVToJoules(kev[mid]), panels)
}
