package section

import (
	"math"

	"compton-scat/internal/constants"
)

func TotalAtRatio(k float64, panels int) float64 {
	energyJ := k * RestEnergyJoules()
	return Total(energyJ, panels)
}

func RestEnergyJoules() float64 {
	return constants.ElectronRestEnergyJ()
}

func LowEnergyCorrection(k float64, panels int) float64 {
	if k <= 0 {
		return 0
	}
	total := TotalAtRatio(k, panels)
	return 1.0 - total/ThomsonCrossSection()
}

func RelativeToThomson(energyJ float64, panels int) float64 {
	return Total(energyJ, panels) / ThomsonCrossSection()
}

func IsBelowThomson(energyJ float64, panels int) bool {
	return Total(energyJ, panels) < ThomsonCrossSection()
}

func MonotoneDecreasing(energiesJ []float64, panels int) bool {
	for i := 1; i < len(energiesJ); i++ {
		if Total(energiesJ[i], panels) >= Total(energiesJ[i-1], panels) {
			return false
		}
	}
	return true
}

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
