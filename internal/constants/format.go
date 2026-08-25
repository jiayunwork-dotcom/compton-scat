package constants

import (
	"fmt"
	"strconv"
)

func FormatScientific(v float64, digits int) string {
	return fmt.Sprintf("%.*e", digits, v)
}

func FormatFixed(v float64, decimals int) string {
	return fmt.Sprintf("%.*f", decimals, v)
}

func FormatSignedFixed(v float64, decimals int) string {
	sign := "+"
	if v < 0 {
		sign = "-"
	}
	return fmt.Sprintf("%s%.*f", sign, decimals, abs(v))
}

func abs(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}

func FormatPercent(ratio float64, decimals int) string {
	return fmt.Sprintf("%.*f%%", decimals, ratio*100)
}

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

func FormatNm(nm float64, digits int) string {
	return strconv.FormatFloat(nm, 'g', digits, 64) + " nm"
}

func FormatM(m float64, digits int) string {
	return fmt.Sprintf("%.*e m", digits, m)
}
