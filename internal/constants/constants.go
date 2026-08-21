// Package constants holds the fixed physical constants used by the
// Compton scattering kinematics: the Planck constant, the electron rest
// mass, the speed of light and the elementary charge. Every derived
// quantity in this package (Compton wavelength, electron rest energy,
// h*c product) is computed from these constants, never from a separate
// memorised number, so that the electron rest energy and the Compton
// wavelength stay mutually consistent by construction.
package constants

// SI base constants (CODATA 2018 recommended values where the value is
// not exact by definition of the SI).
const (
	// PlanckConstant is the Planck constant h in joule-seconds. Its
	// value is exact by the 2019 redefinition of the SI.
	PlanckConstant = 6.62607015e-34

	// ElectronMass is the electron rest mass me in kilograms.
	ElectronMass = 9.1093837015e-31

	// SpeedOfLight is the speed of light in vacuum c in metres per
	// second. Its value is exact by definition.
	SpeedOfLight = 299792458.0

	// ElementaryCharge is the elementary charge e in coulombs. Its
	// value is exact by the 2019 redefinition of the SI.
	ElementaryCharge = 1.602176634e-19

	// VacuumPermittivity is the vacuum electric permittivity epsilon0
	// in farads per metre. It is used to express the classical
	// electron radius from the same charge and mass constants.
	VacuumPermittivity = 8.8541878128e-12
)

// Named constants for the SI prefixes used in the unit helpers.
const (
	// Kilo is the decimal prefix for 10^3.
	Kilo = 1e3
	// Mega is the decimal prefix for 10^6.
	Mega = 1e6
	// Nano is the decimal prefix for 10^-9.
	Nano = 1e-9
)

// ConstantSource describes one fixed constant and where its numeric
// value comes from, so the report commands can quote the provenance of
// every number they print.
type ConstantSource struct {
	// Name is the symbol used in the formulas, e.g. "h".
	Name string
	// Value is the numeric value in the SI unit of the quantity.
	Value float64
	// Unit is the SI unit of the quantity.
	Unit string
	// Source is a short provenance note, e.g. "exact, SI definition".
	Source string
}

// Sources returns the list of base constants accepted by the package,
// in a stable order suitable for printing.
func Sources() []ConstantSource {
	return []ConstantSource{
		{Name: "h", Value: PlanckConstant, Unit: "J*s", Source: "exact, SI definition (2019)"},
		{Name: "me", Value: ElectronMass, Unit: "kg", Source: "CODATA 2018 recommended value"},
		{Name: "c", Value: SpeedOfLight, Unit: "m/s", Source: "exact, SI definition"},
		{Name: "e", Value: ElementaryCharge, Unit: "C", Source: "exact, SI definition (2019)"},
		{Name: "epsilon0", Value: VacuumPermittivity, Unit: "F/m", Source: "CODATA 2018 recommended value"},
	}
}
