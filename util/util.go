package util

import (
	"io"
	"strconv"

	"github.com/fluxynet/gocipe/types"
)

// GetSingleInteger from a map[string][]string (from url.query)
func GetSingleInteger(q map[string][]string, name string) (int, error) {
	var o, ok = q[name]
	if !ok {
		return 0, nil
	}

	if len(o) != 1 {
		return 0, types.ErrInvalidValue
	}

	var v, err = strconv.Atoi(o[0])

	if err == nil {
		return v, nil
	}

	return 0, err
}

// Str returns the pointer of a string
func Str(s string) *string {
	return &s
}

// ReadAll from an io.ReadCloser, returning nil slice on error
func ReadAll(c io.ReadCloser) []byte {
	defer c.Close()
	var d, _ = io.ReadAll(c)
	return d
}
