package cli

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

// TestParseEvent checks that a well-formed event file decodes into the
// solver input.
func TestParseEvent(t *testing.T) {
	ev, err := parseEvent(strings.NewReader(
		`{"name":"511 keV at 90 deg","energy_keV":511.0,"angle_deg":90.0}`), "e.json")
	if err != nil {
		t.Fatalf("parseEvent: %v", err)
	}
	if ev.Name != "511 keV at 90 deg" {
		t.Errorf("name = %q, want 511 keV at 90 deg", ev.Name)
	}
	if ev.EnergyKEV != 511.0 || ev.AngleDeg != 90.0 {
		t.Errorf("event = (%g keV, %g deg), want (511, 90)", ev.EnergyKEV, ev.AngleDeg)
	}
}

// TestParseMissingField checks that a file missing either numeric field
// is rejected with a message naming the field.
func TestParseMissingField(t *testing.T) {
	empty := `{}`
	if _, err := parseEvent(strings.NewReader(empty), "empty.json"); err == nil {
		t.Error("empty object accepted, want error")
	}
	noAngle := `{"energy_keV":511}`
	_, err := parseEvent(strings.NewReader(noAngle), "noangle.json")
	if err == nil {
		t.Error("file without angle_deg accepted, want error")
	}
	if !strings.Contains(err.Error(), "angle_deg") {
		t.Errorf("error %q does not name angle_deg", err)
	}
	noEnergy := `{"angle_deg":90}`
	_, err = parseEvent(strings.NewReader(noEnergy), "noenergy.json")
	if err == nil {
		t.Error("file without energy_keV accepted, want error")
	}
	if !strings.Contains(err.Error(), "energy_keV") {
		t.Errorf("error %q does not name energy_keV", err)
	}
}

// TestParseUnknownField checks that a mistyped field name is rejected
// instead of silently ignored.
func TestParseUnknownField(t *testing.T) {
	json := `{"energy_keV":511,"angle_deg":90,"engergy":500}`
	if _, err := parseEvent(strings.NewReader(json), "unknown.json"); err == nil {
		t.Error("unknown field accepted, want error")
	}
}

// TestParseInvalidEvent checks that a zero energy and an out-of-range
// angle are rejected at parse time with a non-empty message.
func TestParseInvalidEvent(t *testing.T) {
	zero := `{"energy_keV":0,"angle_deg":90}`
	if _, err := parseEvent(strings.NewReader(zero), "zero.json"); err == nil {
		t.Error("zero energy accepted, want error")
	}
	badAngle := `{"energy_keV":511,"angle_deg":200}`
	if _, err := parseEvent(strings.NewReader(badAngle), "theta200.json"); err == nil {
		t.Error("angle 200 accepted, want error")
	}
}

// TestParseEnergyKEV checks the command line energy parser accepts a
// positive number and rejects zero, negatives and non-numbers.
func TestParseEnergyKEV(t *testing.T) {
	if v, err := ParseEnergyKEV("511"); err != nil || v != 511 {
		t.Errorf("ParseEnergyKEV(511) = %g, %v, want 511, nil", v, err)
	}
	if v, err := ParseEnergyKEV("0"); err == nil {
		t.Errorf("ParseEnergyKEV(0) = %g, want error", v)
	}
	if v, err := ParseEnergyKEV("-5"); err == nil {
		t.Errorf("ParseEnergyKEV(-5) = %g, want error", v)
	}
	if v, err := ParseEnergyKEV("abc"); err == nil {
		t.Errorf("ParseEnergyKEV(abc) = %g, want error", v)
	}
}

// TestRunKinematics checks the runner prints lambda, lambda', E' and Ke
// for a valid event file.
func TestRunKinematics(t *testing.T) {
	path := tempEvent(t, `{"name":"example","energy_keV":511.0,"angle_deg":90.0}`)
	var buf bytes.Buffer
	if err := RunKinematics(&buf, []string{path}); err != nil {
		t.Fatalf("RunKinematics: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"lambda", "lambda'", "E'", "Ke", "255.499"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q:\n%s", want, out)
		}
	}
}

// TestRunKinematicsInvalidInput checks that a zero-energy event file
// produces an error that the CLI reports on stderr.
func TestRunKinematicsInvalidInput(t *testing.T) {
	path := tempEvent(t, `{"energy_keV":0,"angle_deg":90}`)
	var buf bytes.Buffer
	err := RunKinematics(&buf, []string{path})
	if err == nil {
		t.Fatal("RunKinematics accepted zero energy, want error")
	}
	if !strings.Contains(err.Error(), "energy") {
		t.Errorf("error %q does not mention energy", err)
	}
}

// TestRunKinematicsBadArgs checks the argument count is enforced.
func TestRunKinematicsBadArgs(t *testing.T) {
	var buf bytes.Buffer
	err := RunKinematics(&buf, nil)
	if err == nil {
		t.Fatal("RunKinematics with no file accepted, want error")
	}
	if !IsUsage(err) {
		t.Errorf("error %q not marked as usage error", err)
	}
}

// TestRunChecksValid checks the checks subcommand passes on the
// canonical example event.
func TestRunChecksValid(t *testing.T) {
	path := tempEvent(t, `{"energy_keV":511.0,"angle_deg":90.0}`)
	var buf bytes.Buffer
	if err := RunChecks(&buf, []string{path}); err != nil {
		t.Fatalf("RunChecks: %v", err)
	}
	if !strings.Contains(buf.String(), "energy conservation") {
		t.Errorf("checks output missing energy conservation line:\n%s", buf.String())
	}
}

// TestRunSectionValid checks the section subcommand prints the total
// cross section for a positive energy.
func TestRunSectionValid(t *testing.T) {
	var buf bytes.Buffer
	if err := RunSection(&buf, []string{"511"}); err != nil {
		t.Fatalf("RunSection: %v", err)
	}
	if !strings.Contains(buf.String(), "total cross section") {
		t.Errorf("section output missing total cross section:\n%s", buf.String())
	}
}

// TestRunSectionInvalid checks a non-positive energy is rejected by the
// section subcommand.
func TestRunSectionInvalid(t *testing.T) {
	var buf bytes.Buffer
	if err := RunSection(&buf, []string{"0"}); err == nil {
		t.Error("RunSection(0) accepted, want error")
	}
}

// TestRunTrendValid checks the trend subcommand reports a monotone
// decreasing total cross section.
func TestRunTrendValid(t *testing.T) {
	var buf bytes.Buffer
	if err := RunTrend(&buf, nil); err != nil {
		t.Fatalf("RunTrend: %v", err)
	}
	if !strings.Contains(buf.String(), "monotone decreasing: yes") {
		t.Errorf("trend output missing monotone decreasing:\n%s", buf.String())
	}
}

func tempEvent(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := dir + "/event.json"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}
