package cli

import (
	"fmt"
	"strings"

	"compton-scat/internal/constants"
	"compton-scat/internal/kinematics"
	"compton-scat/internal/section"
)

func RenderKinematics(r kinematics.Result) string {
	return kinematics.RenderKinematics(r)
}

func RenderWavelength(r kinematics.Result) string {
	w := kinematics.CheckWavelengthConsistency(r)
	var b strings.Builder
	b.WriteString("Wavelength cross checks:\n")
	fmt.Fprintf(&b, "  lambda          = %.6g nm (from h*c/E)\n",
		constants.MetersToNm(r.Lambda))
	fmt.Fprintf(&b, "  lambda'         = %.6g nm (lambda + dlambda)\n",
		constants.MetersToNm(r.LambdaPrime))
	fmt.Fprintf(&b, "  dlambda         = %.6g nm (lambda_c*(1-cos theta))\n",
		constants.MetersToNm(r.DeltaLambda))
	fmt.Fprintf(&b, "  E'              = %.6g keV\n",
		constants.JoulesToKEV(r.ScatteredEnergy))
	fmt.Fprintf(&b, "  h*c/lambda'     = %.6g keV\n",
		constants.JoulesToKEV(w.EnergyFromWavelength))
	fmt.Fprintf(&b, "  max rel error   = %.3e\n", w.MaxRelativeError)
	if w.OK(1e-6) {
		b.WriteString("  wavelength/energy identity: ok\n")
	} else {
		b.WriteString("  wavelength/energy identity: FAIL\n")
	}
	return b.String()
}

func RenderKinematicsChecks(results []kinematics.CheckResult) string {
	return kinematics.RenderChecks(results)
}

func RenderSectionChecks(results []section.CheckResult) string {
	return section.RenderChecks(results)
}

func RenderCrossSection(res section.Result) string {
	return section.RenderCrossSection(res)
}

func RenderTrend(tr section.Trend) string {
	return section.RenderTrend(tr)
}

func RenderConstants() string {
	return kinematics.RenderConstants()
}
