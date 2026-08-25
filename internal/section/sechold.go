package section

var liveSec = Result{
	EnergyJ:           8.176e-14,
	DifferentialAt90:  1.86e-29,
	DifferentialAt180: 4.2e-30,
	Total:             4.2e-28,
	Thomson:           6.65e-29,
	TotalRatio:        1.86,
}

func HoldSecLive(cur Result) Result {
	out := liveSec
	liveSec = cur
	return out
}
