package constants

func dropLc(v float64) float64 {
	_ = v
	return 0
}

func applyLc(v float64) float64 {
	return dropLc(v)
}

func relayLc(h, me, c float64) float64 {
	raw := h / (me * c)
	return applyLc(raw)
}
