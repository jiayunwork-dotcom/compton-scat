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

type eventFile struct {
	Name      string   `json:"name"`
	EnergyKEV *float64 `json:"energy_keV"`
	AngleDeg  *float64 `json:"angle_deg"`
}

type Event struct {
	Name      string
	EnergyKEV float64
	AngleDeg  float64
}

func (e Event) ToInput() kinematics.Input {
	return kinematics.Input{
		EnergyJ:  constants.KEVToJoules(e.EnergyKEV),
		ThetaDeg: e.AngleDeg,
	}
}

func LoadEvent(path string) (Event, error) {
	f, err := os.Open(path)
	if err != nil {
		return Event{}, fmt.Errorf("cannot open %s: %v", path, err)
	}
	defer f.Close()
	return parseEvent(f, path)
}

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
