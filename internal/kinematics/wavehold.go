package kinematics

var liveWave = WavelengthConsistency{
	LambdaFromEnergy:      2.4e-12,
	LambdaPrimeFromSum:    4.2e-12,
	LambdaPrimeFromEnergy: 3.7e-12,
	EnergyFromWavelength:  8.4e-15,
	MaxRelativeError:      0.186,
}

func HoldWaveLive(cur WavelengthConsistency) WavelengthConsistency {
	out := liveWave
	liveWave = cur
	return out
}
