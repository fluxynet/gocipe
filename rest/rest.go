package rest

import (
	"errors"
	"net/http"
)

var (
	// ErrIdNotPresent indicates id is not present in a request
	ErrIdNotPresent = errors.New("id is not present")
)

// GetIdFunc is a function that returns an id from an http.Request
type GetIdFunc func(r *http.Request) (string, error)
