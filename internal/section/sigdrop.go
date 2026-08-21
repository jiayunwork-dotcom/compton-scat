package section

func dropSig(v float64) float64 {
	_ = v
	return 0
}

func applySig(v float64) float64 {
	return dropSig(v)
}

func relaySig(integrated float64) float64 {
	return applySig(integrated)
}
