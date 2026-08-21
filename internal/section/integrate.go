package section

import (
	"math"
)

// Simpson integrates a function over [a, b] with the composite
// Simpson rule using the given number of panels. The panel count must
// be even so the rule is closed; the caller is responsible for passing
// an even count.
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

// SimpsonError estimates the integration error by comparing the
// result on n panels with the result on 2n panels. The difference of
// the two approximations is returned as the absolute error estimate.
func SimpsonError(a, b float64, panels int, f func(float64) float64) float64 {
	coarse := Simpson(a, b, panels, f)
	fine := Simpson(a, b, 2*panels, f)
	return math.Abs(coarse - fine)
}

// solidAngleWeight is the integrand of the total cross section when
// the polar angle is expressed in degrees: the differential cross
// section times the solid angle element 2*pi*sin(theta) and the
// degree-to-radian Jacobian pi/180.
func solidAngleWeight(energyJ, thetaDeg float64) float64 {
	ds := Differential(energyJ, thetaDeg)
	sinT := math.Sin(thetaDeg * math.Pi / 180.0)
	return 2.0 * math.Pi * ds * sinT * (math.Pi / 180.0)
}

// Total computes the total Klein-Nishina cross section in m^2 by
// integrating the differential cross section over the full sphere with
// a composite Simpson rule. The integrand is smooth on [0, pi], so 2000
// panels give a relative error far below one part in a million.
func Total(energyJ float64, panels int) float64 {
	if panels < 2 {
		panels = 2
	}
	// Force an even panel count for the closed Simpson rule.
	if panels%2 == 1 {
		panels++
	}
	integrand := func(deg float64) float64 {
		return solidAngleWeight(energyJ, deg)
	}
	return Simpson(0, 180, panels, integrand)
}

// TotalError returns the total cross section and an absolute error
// estimate based on a finer integration, so a caller can assert the
// integration has converged.
func TotalError(energyJ float64, panels int) (total, err float64) {
	total = Total(energyJ, panels)
	err = SimpsonError(0, 180, panels, func(deg float64) float64 {
		return solidAngleWeight(energyJ, deg)
	})
	return total, err
}
