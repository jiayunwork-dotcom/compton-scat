package section

import (
	"fmt"
	"strings"

	"compton-scat/internal/constants"
)

// RenderCrossSection produces the text printed by the section
// subcommand for one incident energy.
func RenderCrossSection(r Result) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Klein-Nishina cross section at %.3f keV\n",
		constants.JoulesToKEV(r.EnergyJ))
	fmt.Fprintf(&b, "\n")
	fmt.Fprintf(&b, "  differential, theta = 90 deg   %.6e m^2/sr\n", r.DifferentialAt90)
	fmt.Fprintf(&b, "  differential, theta = 180 deg  %.6e m^2/sr\n", r.DifferentialAt180)
	fmt.Fprintf(&b, "  total cross section            %.6e m^2\n", r.Total)
	fmt.Fprintf(&b, "  Thomson limit                  %.6e m^2\n", r.Thomson)
	fmt.Fprintf(&b, "  ratio sigma/sigma_Thomson      %.6f\n", r.TotalRatio)
	return b.String()
}

// RenderTrend produces the energy ladder table printed by the trend
// subcommand, with one row per rung.
func RenderTrend(t Trend) string {
	var b strings.Builder
	b.WriteString("Total Klein-Nishina cross section vs incident energy:\n")
	b.WriteString(fmt.Sprintf("  %12s  %16s  %14s\n", "E (keV)", "sigma (m^2)", "sigma/Thomson"))
	for i := range t.EnergiesKEV {
		fmt.Fprintf(&b, "  %12.3f  %16.6e  %14.6f\n",
			t.EnergiesKEV[i], t.Totals[i], t.Ratios[i])
	}
	if t.Decreasing {
		b.WriteString("  monotone decreasing: yes\n")
	} else {
		b.WriteString("  monotone decreasing: NO\n")
	}
	return b.String()
}

// RenderChecks produces the text printed by the section-checks
// subcommand: one line per cross-section property.
func RenderChecks(results []CheckResult) string {
	var b strings.Builder
	for _, c := range results {
		mark := "ok "
		if !c.Pass {
			mark = "FAIL"
		}
		fmt.Fprintf(&b, "%s  %-46s  %s\n", mark, c.Name, c.Detail)
	}
	return b.String()
}
