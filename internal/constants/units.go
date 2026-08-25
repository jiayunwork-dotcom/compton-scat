package constants

func JoulesToEV(j float64) float64 {
	return j / ElementaryCharge
}

func EVToJoules(ev float64) float64 {
	return ev * ElementaryCharge
}

func JoulesToKEV(j float64) float64 {
	return JoulesToEV(j) / Kilo
}

func KEVToJoules(kev float64) float64 {
	return EVToJoules(kev * Kilo)
}

func MetersToNm(m float64) float64 {
	return m / Nano
}

func NmToMeters(nm float64) float64 {
	return nm * Nano
}

func MetersToFm(m float64) float64 {
	return m / 1e-15
}

func FmToMeters(fm float64) float64 {
	return fm * 1e-15
}

func EnergyInKEV(j float64) float64 {
	return JoulesToKEV(j)
}

func WavelengthInNm(m float64) float64 {
	return MetersToNm(m)
}
