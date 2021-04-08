package util

import "io"

// Close is a helper to close a Closable if it is not null
func Close(c io.Closer) error {
	if c == nil {
		return nil
	}

	return c.Close()
}

// Closed is a helper to close a Closable if it is not null, meant to use with defer, to assign error value
func Closed(c io.Closer, e *error) {
	if c != nil {
		err := c.Close()
		*e = err
	}
}
