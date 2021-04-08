package types

import (
	"errors"
	"strconv"
)

var (
	// ErrInvalidValue is when a value sent is invalid
	ErrInvalidValue = errors.New("invalid value provided")
)

// Type represents a variable type
type Type string

const (
	// Undefined type duh
	Undefined = Type("")

	// Bool indicates native bool
	Bool = Type("bool")

	// String indicates native string
	String = Type("string")

	// Int64 indicates native int64
	Int64 = Type("int64")

	// Float64 indicates native float64
	Float64 = Type("float64")
)

// BoolFromString parses Bool
func BoolFromString(s string) (bool, error) {
	switch s {
	case "true":
		return true, nil
	case "false":
		return false, nil
	}

	return false, ErrInvalidValue
}

// Int64FromString parses Int64
func Int64FromString(s string) (int64, error) {
	return strconv.ParseInt(string(s), 10, 64)
}

// Float64FromString parses Float64
func Float64FromString(s string) (float64, error) {
	return strconv.ParseFloat(string(s), 64)
}

// Default for types
func Default(t Type) interface{} {
	switch t {
	case Bool:
		return true
	case String:
		return ""
	case Int64:
		return 0
	case Float64:
		return float64(0)
	}

	return nil
}

// New returns a pointer for types
func New(t Type) interface{} {
	switch t {
	case Bool:
		return new(bool)
	case String:
		return new(string)
	case Int64:
		return new(int64)
	case Float64:
		return new(float64)
	}

	return nil
}
