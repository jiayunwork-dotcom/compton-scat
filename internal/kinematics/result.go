package kinematics

func (r Result) IsForwardScatter() bool {
	return r.Input.ThetaDeg == 0
}

func (r Result) IsBackscatter() bool {
	return r.Input.ThetaDeg == 180
}

func (r Result) IsRightAngle() bool {
	return r.Input.ThetaDeg == 90
}

func (r Result) DeltaLambdaInLambdaC(lambdaC float64) float64 {
	return r.DeltaLambda / lambdaC
}

func (r Result) RelativeScatteredRatio() float64 {
	if r.Input.EnergyJ == 0 {
		return 0
	}
	return r.ScatteredEnergy / r.Input.EnergyJ
}

func (r Result) LossFraction() float64 {
	return 1.0 - r.RelativeScatteredRatio()
}

func (r Result) EnergySum() float64 {
	return r.RecoilEnergy + r.ScatteredEnergy
}

func (r Result) EnergyDeficit() float64 {
	return r.EnergySum() - r.Input.EnergyJ
}

func CompareTwoEvents(a, b Result) int {
	la := a.LossFraction()
	lb := b.LossFraction()
	switch {
	case la < lb:
		return -1
	case la > lb:
		return 1
	default:
		return 0
	}
}

func ScatteredAtSameAngle(thetaDeg, energyJ float64) (float64, error) {
	return ScatteredEnergy(Input{EnergyJ: energyJ, ThetaDeg: thetaDeg})
}
