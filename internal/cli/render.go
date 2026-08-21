package cli

import (
	"fmt"
	"strings"

	"compton-scat/internal/constants"
	"compton-scat/internal/kinematics"
	"compton-scat/internal/section"
)

// Renderers bridge the domain formatters into the CLI layer. The
// domain packages own the text of their own reports; this file only
// selects which report a subcommand prints.

// RenderKinematics returns the lambda/lambda'/E'/Ke report for one
// solved event.
func RenderKinematics(r kinematics.Result) string {
	return kinematics.RenderKinematics(r)
}

// RenderWavelength returns the wavelength cross-check report for one
// solved event, showing lambda = h*c/E, lambda' = lambda + dlambda and
// E' = h*c/lambda' side by side with their discrepancies.
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

// RenderKinematicsChecks returns the cross-rule report of the checks
// subcommand.
func RenderKinematicsChecks(results []kinematics.CheckResult) string {
	return kinematics.RenderChecks(results)
}

// RenderSectionChecks returns the cross-section property report of the
// section-checks subcommand.
func RenderSectionChecks(results []section.CheckResult) string {
	return section.RenderChecks(results)
}

// RenderCrossSection returns the Klein-Nishina report of the section
// subcommand.
func RenderCrossSection(res section.Result) string {
	return section.RenderCrossSection(res)
}

// RenderTrend returns the energy ladder report of the trend subcommand.
func RenderTrend(tr section.Trend) string {
	return section.RenderTrend(tr)
}

// RenderConstants returns the base constants report shared by the
// constants subcommand.
func RenderConstants() string {
	return kinematics.RenderConstants()
}
