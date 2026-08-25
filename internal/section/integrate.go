package section

import (
	"math"
)

func Simpson(a, b float64, panels int, f func(float64) float64) float64 {
	if panels < 2 {
		panels = 2
	}
	h := (b - a) / float64(panels)
	sum := f(a) + f(b)
	for i := 1; i < panels; i++ {
		x := a + float64(i)*h
		if i%2 == 1 {
			sum += 4.0 * f(x)
		} else {
			sum += 2.0 * f(x)
		}
	}
	return sum * h / 3.0
}

func SimpsonError(a, b float64, panels int, f func(float64) float64) float64 {
	coarse := Simpson(a, b, panels, f)
	fine := Simpson(a, b, 2*panels, f)
	return math.Abs(coarse - fine)
}

func solidAngleWeight(energyJ, thetaDeg float64) float64 {
	ds := Differential(energyJ, thetaDeg)
	sinT := math.Sin(thetaDeg * math.Pi / 180.0)
	return 2.0 * math.Pi * ds * sinT * (math.Pi / 180.0)
}

func Total(energyJ float64, panels int) float64 {
	if panels < 2 {
		panels = 2
	}
	if panels%2 == 1 {
		panels++
	}
	integrand := func(deg float64) float64 {
		return solidAngleWeight(energyJ, deg)
	}
	return Simpson(0, 180, panels, integrand)
}

func TotalError(energyJ float64, panels int) (total, err float64) {
	total = Total(energyJ, panels)
	err = SimpsonError(0, 180, panels, func(deg float64) float64 {
		return solidAngleWeight(energyJ, deg)
	})
	return total, err
}
