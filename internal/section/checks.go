package section

import (
	"fmt"

	"compton-scat/internal/constants"
)

type CheckResult struct {
	Name   string
	Pass   bool
	Detail string
}

func RunAll(panels int) []CheckResult {
	results := []CheckResult{
		checkForwardConstant(),
		checkThomsonLimit(panels),
		checkHighEnergyDrop(panels),
		checkIntegrationConvergence(panels),
		checkTrend(panels),
	}
	return results
}

func AllPass(results []CheckResult) bool {
	for _, c := range results {
		if !c.Pass {
			return false
		}
	}
	return true
}

func Failed(results []CheckResult) []string {
	var failed []string
	for _, c := range results {
		if !c.Pass {
			failed = append(failed, c.Name)
		}
	}
	return failed
}

func checkForwardConstant() CheckResult {
	e := constants.KEVToJoules(511.0)
	got, want, rel := VerifyForwardConstant(e)
	ok := rel <= 1e-12
	return CheckResult{
		Name:   "forward differential equals r_e^2",
		Pass:   ok,
		Detail: fmt.Sprintf("dsigma/dOmega(0) = %.6e m^2/sr, r_e^2 = %.6e m^2/sr (rel %.2e)", got, want, rel),
	}
}

func checkThomsonLimit(panels int) CheckResult {
	e := constants.KEVToJoules(1.0)
	total := Total(e, panels)
	ratio := total / ThomsonCrossSection()
	ok := ratio < 1.0 && ratio > 0.99
	return CheckResult{
		Name:   "low energy approaches Thomson limit",
		Pass:   ok,
		Detail: fmt.Sprintf("sigma(1 keV)/sigma_Thomson = %.6f", ratio),
	}
}

func checkHighEnergyDrop(panels int) CheckResult {
	e := constants.KEVToJoules(1000.0)
	ratio := Total(e, panels) / ThomsonCrossSection()
	ok := ratio < 0.75
	return CheckResult{
		Name:   "total cross section drops at 1 MeV",
		Pass:   ok,
		Detail: fmt.Sprintf("sigma(1 MeV)/sigma_Thomson = %.6f", ratio),
	}
}

func checkIntegrationConvergence(panels int) CheckResult {
	e := constants.KEVToJoules(511.0)
	coarse, err := TotalError(e, panels)
	if err > 1e-10*coarse {
		return CheckResult{
			Name:   "integration converged",
			Pass:   false,
			Detail: fmt.Sprintf("panel error %.3e exceeds 1e-10 of sigma %.6e m^2", err, coarse),
		}
	}
	fine := Total(e, panels*2)
	rel := (coarse - fine) / fine
	if rel < 0 {
		rel = -rel
	}
	return CheckResult{
		Name:   "integration converged",
		Pass:   rel <= 1e-6,
		Detail: fmt.Sprintf("sigma on %d panels = %.9e m^2, on %d panels = %.9e m^2 (rel %.2e)", panels, coarse, panels*2, fine, rel),
	}
}

func checkTrend(panels int) CheckResult {
	ok, at, next := VerifyTrend(panels)
	detail := "sigma decreases from 1 keV to 10 MeV"
	if !ok {
		detail = fmt.Sprintf("sigma rose between ladder rungs %d and %d", at, next)
	}
	return CheckResult{Name: "total cross section decreases with energy", Pass: ok, Detail: detail}
}
