package parser

import (
	"errors"
	"fmt"
	"phi/process"
)

// ParseError is the type of error when parsing a process.
type ParseError struct {
	Pos TokenPos
	Err string // Error string returned from parser.
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("Parse failed at %s: %s", e.Pos, e.Err)
}

// ImmutableNameError is the type of error when trying
// to modify a name without setter methods.
type ImmutableNameError struct {
	Name process.Name
}

func (e ImmutableNameError) Error() string {
	return fmt.Sprintf("cannot set Name: %s is an immutable Name implementation", e.Name.String())
	// return fmt.Sprintf("cannot set Name: ?? is an immutable Name implementation")
}

var ErrInvalid = errors.New("invalid argument")

// UnknownProcessError is the type of error
// when a type switch encounters an unknown
// Process implementation.
//
// The Process implementation may be valid but
// the caller does not anticipate or handle it.
type UnknownProcessError struct {
	Proc process.Process
}

func (e UnknownProcessError) Error() string {
	return fmt.Sprintf("unknown process: %s (type: %T)", e.Proc.String(), e.Proc.String())
}
