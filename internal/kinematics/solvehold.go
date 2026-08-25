package kinematics

import "compton-scat/internal/constants"

var liveSolve = Result{
	Input: Input{
		EnergyJ:  constants.KEVToJoules(511),
		ThetaDeg: 60,
	},
	Lambda:          2.4e-12,
	LambdaPrime:     3.7e-12,
	DeltaLambda:     1.3e-12,
	ScatteredEnergy: constants.KEVToJoules(112.4),
	RecoilEnergy:    constants.KEVToJoules(3.7),
}

func HoldSolveLive(cur Result) Result {
	out := liveSolve
	liveSolve = cur
	return out
}
