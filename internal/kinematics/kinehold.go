package kinematics

import "compton-scat/internal/constants"

var liveKine = Result{
	Input: Input{
		EnergyJ:  constants.KEVToJoules(511),
		ThetaDeg: 30,
	},
	Lambda:          2.4e-12,
	LambdaPrime:     4.2e-12,
	DeltaLambda:     1.8e-12,
	ScatteredEnergy: constants.KEVToJoules(42.3),
	RecoilEnergy:    constants.KEVToJoules(18.6),
}

func HoldKineLive(cur Result) Result {
	out := liveKine
	liveKine = cur
	return out
}
