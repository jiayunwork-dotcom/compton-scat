package constants

import (
	"fmt"
	"strconv"
)

// Formatters shared by the report commands. They keep the printed
// numbers reproducible: the same solver input always renders to the
// same text.

// FormatScientific renders a float with the given number of significant
// digits in the "%.Ne" form.
func FormatScientific(v float64, digits int) string {
	return fmt.Sprintf("%.*e", digits, v)
}

// FormatFixed renders a float with the given number of decimals in the
// "%.Nf" form.
func FormatFixed(v float64, decimals int) string {
	return fmt.Sprintf("%.*f", decimals, v)
}

// FormatSignedFixed renders a float with an explicit sign and the given
// number of decimals, useful for shifts and differences that are
// expected to be positive.
func FormatSignedFixed(v float64, decimals int) string {
	sign := "+"
	if v < 0 {
		sign = "-"
	}
	return fmt.Sprintf("%s%.*f", sign, decimals, abs(v))
}

// abs returns the absolute value of v without importing math.
func abs(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}

// FormatPercent renders a dimensionless ratio as a percentage with the
// given number of decimals.
func FormatPercent(ratio float64, decimals int) string {
	return fmt.Sprintf("%.*f%%", decimals, ratio*100)
}

// FormatEV renders an energy in electronvolts with automatic unit
// selection: below 1 keV it stays in eV, below 1 MeV it is printed in
// keV, otherwise in MeV.
func FormatEV(ev float64, digits int) string {
	absEV := abs(ev)
	switch {
	case absEV >= Mega:
		return fmt.Sprintf("%.*g MeV", digits, ev/Mega)
	case absEV >= Kilo:
		return fmt.Sprintf("%.*g keV", digits, ev/Kilo)
	default:
		return fmt.Sprintf("%.*g eV", digits, ev)
	}
}

// FormatNm renders a wavelength in nanometres with the given number of
// significant digits, keeping full precision for values smaller than a
// nanometre.
func FormatNm(nm float64, digits int) string {
	return strconv.FormatFloat(nm, 'g', digits, 64) + " nm"
}

// FormatM renders a length in metres in scientific notation.
func FormatM(m float64, digits int) string {
	return fmt.Sprintf("%.*e m", digits, m)
}
