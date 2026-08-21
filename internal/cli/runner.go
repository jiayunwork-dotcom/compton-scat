package cli

import (
	"fmt"
	"io"

	"compton-scat/internal/constants"
	"compton-scat/internal/kinematics"
	"compton-scat/internal/section"
)

// RunKinematics executes the kinematics subcommand: it reads an event
// file, solves the Compton kinematics and writes the lambda, lambda',
// E' and Ke report to out.
func RunKinematics(out io.Writer, args []string) error {
	if len(args) != 1 {
		return InvalidArgs("kinematics needs exactly one event JSON file")
	}
	ev, err := LoadEvent(args[0])
	if err != nil {
		return err
	}
	r, err := kinematics.Solve(ev.ToInput())
	if err != nil {
		return err
	}
	_, err = io.WriteString(out, RenderKinematics(r))
	return err
}

// RunWavelength executes the wavelength subcommand: it reads an event
// file and prints the wavelength side of the Compton relations with
// the mutual consistency checks.
func RunWavelength(out io.Writer, args []string) error {
	if len(args) != 1 {
		return InvalidArgs("wavelength needs exactly one event JSON file")
	}
	ev, err := LoadEvent(args[0])
	if err != nil {
		return err
	}
	r, err := kinematics.Solve(ev.ToInput())
	if err != nil {
		return err
	}
	_, err = io.WriteString(out, RenderWavelength(r))
	return err
}

// RunSection executes the section subcommand: it takes an incident
// energy in keV and prints the Klein-Nishina differential and total
// cross sections.
func RunSection(out io.Writer, args []string) error {
	if len(args) != 1 {
		return InvalidArgs("section needs exactly one energy in keV")
	}
	kev, err := ParseEnergyKEV(args[0])
	if err != nil {
		return err
	}
	res, err := section.CrossSection(constants.KEVToJoules(kev))
	if err != nil {
		return err
	}
	_, err = io.WriteString(out, RenderCrossSection(res))
	return err
}

// RunTrend executes the trend subcommand: it prints the total cross
// section on a logarithmic energy ladder and reports whether it falls
// monotonically.
func RunTrend(out io.Writer, args []string) error {
	if len(args) != 0 {
		return InvalidArgs("trend takes no arguments")
	}
	tr := section.ComputeTrend(section.DefaultIntegrationNodes)
	_, err := io.WriteString(out, RenderTrend(tr))
	return err
}

// RunChecks executes the checks subcommand: it reads an event file,
// runs every cross rule that applies and exits non-zero (via the
// returned error) when any rule fails.
func RunChecks(out io.Writer, args []string) error {
	if len(args) != 1 {
		return InvalidArgs("checks needs exactly one event JSON file")
	}
	ev, err := LoadEvent(args[0])
	if err != nil {
		return err
	}
	results, err := kinematics.RunAll(ev.ToInput())
	if err != nil {
		return err
	}
	if _, err := io.WriteString(out, RenderKinematicsChecks(results)); err != nil {
		return err
	}
	if !kinematics.AllPass(results) {
		return fmt.Errorf("kinematics cross-rule checks failed: %v", kinematics.Failed(results))
	}
	return nil
}

// RunSectionChecks executes the section-checks subcommand: it runs
// every cross-section property and exits non-zero when any fails.
func RunSectionChecks(out io.Writer, args []string) error {
	if len(args) != 0 {
		return InvalidArgs("section-checks takes no arguments")
	}
	results := section.RunAll(section.DefaultIntegrationNodes)
	if _, err := io.WriteString(out, RenderSectionChecks(results)); err != nil {
		return err
	}
	if !section.AllPass(results) {
		return fmt.Errorf("cross-section checks failed: %v", section.Failed(results))
	}
	return nil
}

// RunConstants executes the constants subcommand: it prints the base
// constants and the derived Compton wavelength and electron rest
// energy with their provenance.
func RunConstants(out io.Writer, args []string) error {
	if len(args) != 0 {
		return InvalidArgs("constants takes no arguments")
	}
	_, err := io.WriteString(out, RenderConstants())
	return err
}
