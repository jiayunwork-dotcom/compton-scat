package cli

import (
	"errors"
	"fmt"
)

func InvalidArgs(format string, args ...any) error {
	return UsageError(fmt.Errorf(format, args...))
}

var ErrNoCommand = errors.New("missing subcommand")

var ErrUnknownCommand = errors.New("unknown subcommand")

func wrap(cmd string, err error) error {
	return fmt.Errorf("%s: %v", cmd, err)
}

type usageError struct {
	err error
}

func UsageError(err error) error {
	return &usageError{err: err}
}

func (u *usageError) Error() string {
	return u.err.Error()
}

func (u *usageError) Unwrap() error {
	return u.err
}

func IsUsage(err error) bool {
	var u *usageError
	return errors.As(err, &u)
}
