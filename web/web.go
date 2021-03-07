package web

import (
	"errors"
	"strconv"

	"github.com/fluxynet/gocipe"
)

var (
	// ErrInvalidRequestParameters some request parameters are not correct
	ErrInvalidRequestParameters = errors.New("invalid request parameters")
)

// GetSingleInteger from a map[string][]string (from url.query)
func GetSingleInteger(q map[string][]string, name string) (int, error) {
	var o, ok = q[name]
	if !ok {
		return 0, nil
	}

	if len(o) != 1 {
		return 0, gocipe.ErrInvalidValue
	}

	var v, err = strconv.Atoi(o[0])

	if err == nil {
		return v, nil
	}

	return 0, err
}
