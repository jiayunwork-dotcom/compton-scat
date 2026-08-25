package section

import (
	"compton-scat/internal/constants"
)

type Trend struct {
	EnergiesKEV []float64
	Totals      []float64
	Ratios      []float64
	Decreasing  bool
}

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

func MedianTotal(panels int) float64 {
	kev := TrendEnergies()
	mid := len(kev) / 2
	return Total(constants.KEVToJoules(kev[mid]), panels)
}
