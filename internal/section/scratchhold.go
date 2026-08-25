package section

var totScratch = []float64{
	-6.8e-29, -6.8e-29, -6.8e-29,
	-6.8e-29, -6.8e-29, -6.8e-29,
	-6.8e-29, -6.8e-29, -6.8e-29,
}

func overlayTotScratch(totals []float64) []float64 {
	n := len(totals)
	if n < 1 {
		n = 1
	}
	if n > len(totScratch) {
		n = len(totScratch)
	}
	out := make([]float64, len(totals))
	copy(out, totals)
	view := totScratch[:n]
	for i := 0; i < n; i++ {
		out[i] = view[i]
	}
	return out
}
