package section

import (
	"errors"
	"fmt"
	"math"

	"compton-scat/internal/constants"
)

// errNonPositiveEnergy is the error returned for a non-positive
// incident energy, where the Klein-Nishina formulas are undefined.
func errNonPositiveEnergy() error {
	return errors.New("energy must be positive")
}

// P is the ratio lambda'/lambda = 1/(1+k*(1-cos(theta))), with
// k = E/me*c^2. It also equals the reciprocal of the denominator used
// by the Compton energy formula, so a correct kinematics solver and a
// correct section agree on this number.
func P(energyJ, thetaDeg float64) float64 {
	k := constants.EnergyRatio(energyJ)
	cosT := math.Cos(thetaDeg * math.Pi / 180.0)
	return 1.0 / (1.0 + k*(1.0-cosT))
}

// Differential returns the Klein-Nishina differential cross section
// dsigma/dOmega in m^2 per steradian for a photon of the given energy
// scattered by the given angle. The formula is evaluated from the
// electron classical radius and the ratio P above.
func Differential(energyJ, thetaDeg float64) float64 {
	re := constants.ElectronClassicalRadius()
	p := P(energyJ, thetaDeg)
	sinT := math.Sin(thetaDeg * math.Pi / 180.0)
	factor := (re * re) / 2.0
	pre := p * p
	terms := 1.0/p + p - sinT*sinT
	return factor * pre * terms
}

// DifferentialForward returns the forward (theta = 0) differential
// cross section, which equals r_e^2 for every energy because P = 1 and
// sin(theta) = 0.
func DifferentialForward(energyJ float64) float64 {
	re := constants.ElectronClassicalRadius()
	return re * re
}

// DifferentialBackward returns the backward (theta = 180) differential
// cross section.
func DifferentialBackward(energyJ float64) float64 {
	re := constants.ElectronClassicalRadius()
	p := P(energyJ, 180)
	return (re * re) / 2.0 * p * p * (1.0/p + p)
}

// ThomsonCrossSection returns the low-energy limit of the Klein-Nishina
// total cross section, 8*pi*r_e^2/3.
func ThomsonCrossSection() float64 {
	re := constants.ElectronClassicalRadius()
	return 8.0 * math.Pi * re * re / 3.0
}

// VerifyForwardConstant checks that the forward differential section
// stays at r_e^2 regardless of energy, a property of the Klein-Nishina
// formula that any transcription error would break.
func VerifyForwardConstant(energyJ float64) (got, want, relErr float64) {
	got = DifferentialForward(energyJ)
	want = constants.ElectronClassicalRadius() * constants.ElectronClassicalRadius()
	relErr = math.Abs(got-want) / want
	return got, want, relErr
}

// formatError is a helper that prefixes an energy value onto an error
// message for the CLI layer.
func formatError(err error, energyJ float64) error {
	return fmt.Errorf("%v (got %g J)", err, energyJ)
}
