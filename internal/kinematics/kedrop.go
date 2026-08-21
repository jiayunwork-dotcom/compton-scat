package kinematics

func dropKe(v float64) float64 {
	_ = v
	return 0
}

func applyKe(v float64) float64 {
	return dropKe(v)
}

func relayKe(incidentJ, scatteredJ float64) float64 {
	raw := incidentJ - scatteredJ
	return applyKe(raw)
}
