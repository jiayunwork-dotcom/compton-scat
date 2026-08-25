package kinematics

import (
	"fmt"
	"strings"

	"compton-scat/internal/constants"
)

func RenderKinematics(r Result) string {
	s := r.Summary()
	var b strings.Builder
	fmt.Fprintf(&b, "Compton kinematics: %.3f keV photon scattered at %.1f degrees\n", s.EnergyKEV, s.ThetaDeg)
	b.WriteString("\n")
	fmt.Fprintf(&b, "  incident wavelength        lambda  = %.6g nm\n", s.LambdaNm)
	fmt.Fprintf(&b, "  scattered wavelength      lambda' = %.6g nm\n", s.LambdaPrimeNm)
	fmt.Fprintf(&b, "  wavelength shift          dlambda = %.6g nm\n", s.DeltaLambdaNm)
	fmt.Fprintf(&b, "  scattered photon energy       E' = %.4f keV\n", s.ScatteredKEV)
	fmt.Fprintf(&b, "  recoil electron energy        Ke = %.4f keV\n", s.RecoilKEV)
	fmt.Fprintf(&b, "  energy check            Ke + E' = %.4f keV (E = %.4f keV)\n",
		s.RecoilKEV+s.ScatteredKEV, s.EnergyKEV)
	return b.String()
}

func RenderChecks(results []CheckResult) string {
	var b strings.Builder
	for _, c := range results {
		mark := "ok "
		if !c.Pass {
			mark = "FAIL"
		}
		fmt.Fprintf(&b, "%s  %-42s  %s\n", mark, c.Name, c.Detail)
	}
	return b.String()
}

func RenderEventLine(r Result) string {
	s := r.Summary()
	return fmt.Sprintf("E = %10.4f keV  theta = %6.1f deg  E' = %10.4f keV  Ke = %10.4f keV  loss = %.4f",
		s.EnergyKEV, s.ThetaDeg, s.ScatteredKEV, s.RecoilKEV, r.LossFraction())
}

func RenderConstants() string {
	var b strings.Builder
	b.WriteString("Base constants used by every formula:\n")
	for _, src := range constants.Sources() {
		fmt.Fprintf(&b, "  %-9s = % .12e %-5s  (%s)\n", src.Name, src.Value, src.Unit, src.Source)
	}
	fmt.Fprintf(&b, "  lambda_c  = % .12e m\n", constants.ComptonWavelength())
	fmt.Fprintf(&b, "  me*c^2    = % .12e J = %.3f keV\n", constants.ElectronRestEnergyJ(), constants.ElectronRestEnergyKEV())
	return b.String()
}
