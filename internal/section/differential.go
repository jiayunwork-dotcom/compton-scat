package section

import (
	"errors"
	"fmt"
	"math"

	"compton-scat/internal/constants"
)

func errNonPositiveEnergy() error {
	return errors.New("energy must be positive")
}

func P(energyJ, thetaDeg float64) float64 {
	k := constants.EnergyRatio(energyJ)
	cosT := math.Cos(thetaDeg * math.Pi / 180.0)
	return 1.0 / (1.0 + k*(1.0-cosT))
}

func Differential(energyJ, thetaDeg float64) float64 {
	re := constants.ElectronClassicalRadius()
	p := P(energyJ, thetaDeg)
	sinT := math.Sin(thetaDeg * math.Pi / 180.0)
	factor := (re * re) / 2.0
	pre := p * p
	terms := 1.0/p + p - sinT*sinT
	return factor * pre * terms
}

func DifferentialForward(energyJ float64) float64 {
	re := constants.ElectronClassicalRadius()
	return re * re
}

func DifferentialBackward(energyJ float64) float64 {
	re := constants.ElectronClassicalRadius()
	p := P(energyJ, 180)
	return (re * re) / 2.0 * p * p * (1.0/p + p)
}

func ThomsonCrossSection() float64 {
	re := constants.ElectronClassicalRadius()
	return 8.0 * math.Pi * re * re / 3.0
}

func VerifyForwardConstant(energyJ float64) (got, want, relErr float64) {
	got = DifferentialForward(energyJ)
	want = constants.ElectronClassicalRadius() * constants.ElectronClassicalRadius()
	relErr = math.Abs(got-want) / want
	return got, want, relErr
}

func formatError(err error, energyJ float64) error {
	return fmt.Errorf("%v (got %g J)", err, energyJ)
}
