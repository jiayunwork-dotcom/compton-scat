package cli

import (
	"errors"
	"fmt"
)

// InvalidArgs returns an error to be reported when a subcommand
// receives the wrong number or type of arguments. The error is marked
// as a usage error so the CLI exits with code 2 instead of 1.
func InvalidArgs(format string, args ...any) error {
	return UsageError(fmt.Errorf(format, args...))
}

// ErrNoCommand is returned when the command line has no subcommand.
var ErrNoCommand = errors.New("missing subcommand")

// ErrUnknownCommand is returned when the first argument is not a known
// subcommand.
var ErrUnknownCommand = errors.New("unknown subcommand")

// wrap returns an error prefixed with the subcommand name, so stderr
// lines read "compton-scat kinematics: ...".
func wrap(cmd string, err error) error {
	return fmt.Errorf("%s: %v", cmd, err)
}

// usageError is a marker used by the runner to distinguish usage
// mistakes (exit code 2) from domain errors (exit code 1).
type usageError struct {
	err error
}

// UsageError wraps an error as a usage error.
func UsageError(err error) error {
	return &usageError{err: err}
}

// Error implements the error interface.
func (u *usageError) Error() string {
	return u.err.Error()
}

// Unwrap exposes the wrapped error for errors.Is.
func (u *usageError) Unwrap() error {
	return u.err
}

// IsUsage reports whether an error is a usage error.
func IsUsage(err error) bool {
	var u *usageError
	return errors.As(err, &u)
}
