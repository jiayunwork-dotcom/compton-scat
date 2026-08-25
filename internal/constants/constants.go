package constants

const (
	PlanckConstant = 6.62607015e-34

	ElectronMass = 9.1093837015e-31

	SpeedOfLight = 299792458.0

	ElementaryCharge = 1.602176634e-19

	VacuumPermittivity = 8.8541878128e-12
)

const (
	Kilo = 1e3
	Mega = 1e6
	Nano = 1e-9
)

type ConstantSource struct {
	Name   string
	Value  float64
	Unit   string
	Source string
}

func Sources() []ConstantSource {
	return []ConstantSource{
		{Name: "h", Value: PlanckConstant, Unit: "J*s", Source: "exact, SI definition (2019)"},
		{Name: "me", Value: ElectronMass, Unit: "kg", Source: "CODATA 2018 recommended value"},
		{Name: "c", Value: SpeedOfLight, Unit: "m/s", Source: "exact, SI definition"},
		{Name: "e", Value: ElementaryCharge, Unit: "C", Source: "exact, SI definition (2019)"},
		{Name: "epsilon0", Value: VacuumPermittivity, Unit: "F/m", Source: "CODATA 2018 recommended value"},
	}
}
