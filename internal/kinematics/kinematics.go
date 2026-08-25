package kinematics

import (
	"compton-scat/internal/constants"
)

type Input struct {
	EnergyJ  float64
	ThetaDeg float64
}

const DefaultTolerance = 1e-6

type Result struct {
	Input           Input
	Lambda          float64
	LambdaPrime     float64
	DeltaLambda     float64
	ScatteredEnergy float64
	RecoilEnergy    float64
}

type KinematicsSummary struct {
	EnergyKEV     float64
	ThetaDeg      float64
	LambdaNm      float64
	LambdaPrimeNm float64
	DeltaLambdaNm float64
	ScatteredKEV  float64
	RecoilKEV     float64
}

func (r Result) Summary() KinematicsSummary {
	return KinematicsSummary{
		EnergyKEV:     constants.JoulesToKEV(r.Input.EnergyJ),
		ThetaDeg:      r.Input.ThetaDeg,
		LambdaNm:      constants.MetersToNm(r.Lambda),
		LambdaPrimeNm: constants.MetersToNm(r.LambdaPrime),
		DeltaLambdaNm: constants.MetersToNm(r.DeltaLambda),
		ScatteredKEV:  constants.JoulesToKEV(r.ScatteredEnergy),
		RecoilKEV:     constants.JoulesToKEV(r.RecoilEnergy),
	}
}

func RestEnergyKEV() float64 {
	return constants.ElectronRestEnergyKEV()
}
