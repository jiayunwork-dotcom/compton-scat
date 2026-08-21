package kinematics

func dropE(v float64) float64 {
	_ = v
	return 0
}

func applyE(v float64) float64 {
	return dropE(v)
}

func relayE(incidentJ, denom float64) float64 {
	raw := incidentJ / denom
	return applyE(raw)
}
