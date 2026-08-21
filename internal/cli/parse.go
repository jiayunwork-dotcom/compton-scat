package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"

	"compton-scat/internal/constants"
	"compton-scat/internal/kinematics"
)

// eventFile is the JSON layout accepted by the kinematics, wavelength
// and checks commands. Both numeric fields are pointers so a missing
// field can be told apart from a zero value, and both are required.
type eventFile struct {
	// Name is an optional human readable label for the event.
	Name string `json:"name"`
	// EnergyKEV is the incident photon energy in keV.
	EnergyKEV *float64 `json:"energy_keV"`
	// AngleDeg is the scattering angle in degrees.
	AngleDeg *float64 `json:"angle_deg"`
}

// Event is a validated scattering event ready for the solver.
type Event struct {
	// Name is the optional label from the file.
	Name string
	// EnergyKEV is the incident photon energy in keV.
	EnergyKEV float64
	// AngleDeg is the scattering angle in degrees.
	AngleDeg float64
}

// ToInput converts the event into the solver input in joules.
func (e Event) ToInput() kinematics.Input {
	return kinematics.Input{
		EnergyJ:  constants.KEVToJoules(e.EnergyKEV),
		ThetaDeg: e.AngleDeg,
	}
}

// LoadEvent reads an event file from disk and validates it.
func LoadEvent(path string) (Event, error) {
	f, err := os.Open(path)
	if err != nil {
		return Event{}, fmt.Errorf("cannot open %s: %v", path, err)
	}
	defer f.Close()
	return parseEvent(f, path)
}

// parseEvent decodes an event file from a reader. It rejects unknown
// fields and missing numeric fields with a message that names the
// offending field.
func parseEvent(r io.Reader, path string) (Event, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	var raw eventFile
	if err := dec.Decode(&raw); err != nil {
		return Event{}, fmt.Errorf("%s: invalid JSON: %v", path, err)
	}
	if raw.EnergyKEV == nil {
		return Event{}, fmt.Errorf("%s: missing field \"energy_keV\"", path)
	}
	if raw.AngleDeg == nil {
		return Event{}, fmt.Errorf("%s: missing field \"angle_deg\"", path)
	}
	if math.IsNaN(*raw.EnergyKEV) || math.IsInf(*raw.EnergyKEV, 0) {
		return Event{}, fmt.Errorf("%s: energy_keV must be finite", path)
	}
	if math.IsNaN(*raw.AngleDeg) || math.IsInf(*raw.AngleDeg, 0) {
		return Event{}, fmt.Errorf("%s: angle_deg must be finite", path)
	}
	ev := Event{Name: raw.Name, EnergyKEV: *raw.EnergyKEV, AngleDeg: *raw.AngleDeg}
	if err := kinematics.Validate(ev.ToInput()); err != nil {
		return Event{}, fmt.Errorf("%s: %v", path, err)
	}
	return ev, nil
}

// ParseEnergyKEV parses a command line argument as an energy in keV,
// rejecting values that are not numbers or are not positive.
func ParseEnergyKEV(s string) (float64, error) {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, InvalidArgs(fmt.Sprintf("energy must be a number in keV, got %q", s))
	}
	if math.IsNaN(v) || math.IsInf(v, 0) || v <= 0 {
		return 0, fmt.Errorf("energy must be a positive finite number of keV (got %g)", v)
	}
	return v, nil
}
