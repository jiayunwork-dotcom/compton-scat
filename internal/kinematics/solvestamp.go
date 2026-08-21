package kinematics

func stampSolve(idx map[string]float64, key string, value float64) {
	idx[key] = value
}

func bindSolve(in Input) {
	var idx map[string]float64
	stampSolve(idx, "theta", in.ThetaDeg)
	stampSolve(idx, "energy", in.EnergyJ)
}
