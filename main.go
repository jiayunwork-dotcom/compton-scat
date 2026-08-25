package main

import (
	"flag"
	"fmt"
	"os"

	"compton-scat/internal/api"
	"compton-scat/internal/cli"
)

const version = "1.0.0"

const usage = `compton-scat: Compton scattering kinematics.

Given the incident photon energy and the scattering angle of a photon
off a free electron at rest, print the wavelength shift, the scattered
photon energy and the recoil electron kinetic energy. All formulas use
the electron Compton wavelength lambda_c = h/(me*c) and the electron
rest energy me*c^2 from the same h, me and c.

usage:
  compton-scat kinematics <event.json>      print lambda, lambda', E', Ke
  compton-scat wavelength <event.json>      wavelength/energy cross checks
  compton-scat section <energy-keV>         Klein-Nishina cross section
  compton-scat trend                        total cross section vs energy
  compton-scat checks <event.json>          kinematics cross-rule checks
  compton-scat section-checks               cross-section property checks
  compton-scat constants                    base constants and derived values
  compton-scat help                         print this message
  compton-scat -version                     print the version

The event JSON has the form {"name": "...", "energy_keV": ..,
"angle_deg": ..} with the energy in keV and the angle in degrees in
[0, 180]. Non-positive energies, out-of-range angles, missing fields
and unreadable files are reported on stderr and the command exits
non-zero.
`

func main() {
	var showVersion, showHelp bool
	httpAddr := flag.String("http", "", "serve the HTTP API on this address (e.g. :8080)")
	flag.BoolVar(&showVersion, "version", false, "print the version and exit")
	flag.BoolVar(&showHelp, "h", false, "print this usage message and exit")
	flag.BoolVar(&showHelp, "help", false, "print this usage message and exit")
	flag.Usage = func() { fmt.Fprint(os.Stderr, usage) }
	flag.Parse()

	if *httpAddr != "" {
		if err := api.Serve(*httpAddr); err != nil {
			fmt.Fprintf(os.Stderr, "compton-scat: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if showVersion {
		fmt.Printf("compton-scat %s\n", version)
		return
	}
	if showHelp {
		fmt.Print(usage)
		return
	}

	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}

	var err error
	switch args[0] {
	case "kinematics":
		err = cli.RunKinematics(os.Stdout, args[1:])
	case "wavelength":
		err = cli.RunWavelength(os.Stdout, args[1:])
	case "section":
		err = cli.RunSection(os.Stdout, args[1:])
	case "trend":
		err = cli.RunTrend(os.Stdout, args[1:])
	case "checks":
		err = cli.RunChecks(os.Stdout, args[1:])
	case "section-checks":
		err = cli.RunSectionChecks(os.Stdout, args[1:])
	case "constants":
		err = cli.RunConstants(os.Stdout, args[1:])
	case "help":
		fmt.Print(usage)
		return
	default:
		fmt.Fprintf(os.Stderr, "compton-scat: unknown command %q\n\n%s", args[0], usage)
		os.Exit(2)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "compton-scat: %v\n", err)
		if cli.IsUsage(err) {
			os.Exit(2)
		}
		os.Exit(1)
	}
}
