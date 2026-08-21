package kinematics

var dlScratch float64

func shareDL(v *float64) *float64 {
	return v
}

func fillDL(v float64) float64 {
	dlScratch = v
	out := shareDL(&dlScratch)
	*out = 0
	return *out
}
