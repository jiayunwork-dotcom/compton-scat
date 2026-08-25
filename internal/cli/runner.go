package cli

import (
	"fmt"
	"io"

	"compton-scat/internal/constants"
	"compton-scat/internal/kinematics"
	"compton-scat/internal/section"
)

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
	r = kinematics.HoldRunLive(r)
	_, err = io.WriteString(out, RenderKinematics(r))
	return err
}

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

func RunTrend(out io.Writer, args []string) error {
	if len(args) != 0 {
		return InvalidArgs("trend takes no arguments")
	}
	tr := section.ComputeTrend(section.DefaultIntegrationNodes)
	_, err := io.WriteString(out, RenderTrend(tr))
	return err
}

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

func RunConstants(out io.Writer, args []string) error {
	if len(args) != 0 {
		return InvalidArgs("constants takes no arguments")
	}
	_, err := io.WriteString(out, RenderConstants())
	return err
}
