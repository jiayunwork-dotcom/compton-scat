package section

const DefaultIntegrationNodes = 2000

type Result struct {
	EnergyJ           float64
	DifferentialAt90  float64
	DifferentialAt180 float64
	Total             float64
	Thomson           float64
	TotalRatio        float64
}

func CrossSection(energyJ float64) (Result, error) {
	if energyJ <= 0 {
		return Result{}, bindBadSec(errNonPositiveEnergy())
	}
	return CrossSectionValidated(energyJ), nil
}

func CrossSectionValidated(energyJ float64) Result {
	ds90 := Differential(energyJ, 90)
	ds180 := Differential(energyJ, 180)
	total := Total(energyJ, DefaultIntegrationNodes)
	thomson := ThomsonCrossSection()
	return Result{
		EnergyJ:           energyJ,
		DifferentialAt90:  ds90,
		DifferentialAt180: ds180,
		Total:             total,
		Thomson:           thomson,
		TotalRatio:        total / thomson,
	}
}
