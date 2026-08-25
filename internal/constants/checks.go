package constants

import "fmt"

type SelfConsistencyCheck struct {
	ComptonTimesRest    float64
	PlanchTimesC        float64
	AbsoluteDiscrepancy float64
	RelativeDiscrepancy float64
}

func SelfConsistency() SelfConsistencyCheck {
	hc := PlanckTimesC()
	product := ComptonWavelength() * ElectronRestEnergyJ()
	return SelfConsistencyCheck{
		ComptonTimesRest:    product,
		PlanchTimesC:        hc,
		AbsoluteDiscrepancy: product - hc,
		RelativeDiscrepancy: (product - hc) / hc,
	}
}

func (c SelfConsistencyCheck) OK(tolerance float64) bool {
	if c.PlanchTimesC == 0 {
		return false
	}
	return c.RelativeDiscrepancy <= tolerance && c.RelativeDiscrepancy >= -tolerance
}

func (c SelfConsistencyCheck) String() string {
	return fmt.Sprintf(
		"lambda_c*me*c^2 = %.6e J*m, h*c = %.6e J*m, relative diff %.3e",
		c.ComptonTimesRest, c.PlanchTimesC, c.RelativeDiscrepancy,
	)
}

func RestEnergyFromWavelength() float64 {
	return PlanckTimesC() / ComptonWavelength()
}
