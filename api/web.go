package api

import (
	"errors"
)

var (
	// ErrInvalidRequestParameters some request parameters are not correct
	ErrInvalidRequestParameters = errors.New("invalid request parameters")
)
