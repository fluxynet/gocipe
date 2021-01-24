package chi

import (
	"github.com/fluxynet/gocipe/rest"
	"github.com/go-chi/chi/v5"
	"net/http"
)

// GetIdFunc is a function that returns an id from an http.Request
func GetIdFunc(r *http.Request) (string, error) {
	var id = chi.URLParam(r, "id")

	if id == "" {
		return "", rest.ErrIdNotPresent
	}

	return id, nil
}
