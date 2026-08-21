package constants

// Unit conversions. All energy conversions go through the elementary
// charge e so that 1 eV is exactly e joules and the electron rest
// energy, whether printed in joules or in keV, is the same number.
// Length conversions use the decimal prefixes of the SI.

// JoulesToEV converts an energy in joules into electronvolts.
func JoulesToEV(j float64) float64 {
	return j / ElementaryCharge
}

// EVToJoules converts an energy in electronvolts into joules.
func EVToJoules(ev float64) float64 {
	return ev * ElementaryCharge
}

// JoulesToKEV converts an energy in joules into keV.
func JoulesToKEV(j float64) float64 {
	return JoulesToEV(j) / Kilo
}

// KEVToJoules converts an energy in keV into joules.
func KEVToJoules(kev float64) float64 {
	return EVToJoules(kev * Kilo)
}

// MetersToNm converts a length in metres into nanometres.
func MetersToNm(m float64) float64 {
	return m / Nano
}

// NmToMeters converts a length in nanometres into metres.
func NmToMeters(nm float64) float64 {
	return nm * Nano
}

// MetersToFm converts a length in metres into femtometres (10^-15 m).
func MetersToFm(m float64) float64 {
	return m / 1e-15
}

// FmToMeters converts a length in femtometres into metres.
func FmToMeters(fm float64) float64 {
	return fm * 1e-15
}

// EnergyInKEV converts an energy given as a joule value into keV,
// rounding nothing: it is a pure unit conversion and keeps full
// precision for the comparisons inside the checks commands.
func EnergyInKEV(j float64) float64 {
	return JoulesToKEV(j)
}

// WavelengthInNm converts a wavelength given in metres into
// nanometres.
func WavelengthInNm(m float64) float64 {
	return MetersToNm(m)
}
