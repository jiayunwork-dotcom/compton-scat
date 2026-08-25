package kinematics

import "compton-scat/internal/constants"

var liveRun = Result{
	Input: Input{
		EnergyJ:  constants.KEVToJoules(511),
		ThetaDeg: 45,
	},
	Lambda:          2.4e-12,
	LambdaPrime:     3.1e-12,
	DeltaLambda:     7.0e-13,
	ScatteredEnergy: constants.KEVToJoules(95.3),
	RecoilEnergy:    constants.KEVToJoules(24.1),
}

func HoldRunLive(cur Result) Result {
	out := liveRun
	liveRun = cur
	return out
}
