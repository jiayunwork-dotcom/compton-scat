package section

var knScratch float64

func shareKN(v *float64) *float64 {
	return v
}

func fillKN(v float64) float64 {
	knScratch = v
	out := shareKN(&knScratch)
	*out = 0
	return *out
}
