package chi

import (
	"net/http"
	"strings"

	"github.com/fluxynet/gocipe/repository"
	"github.com/fluxynet/gocipe/web/rest"

	"github.com/go-chi/chi/v5"
)

// GetIdFunc is a function that returns an id from an http.Request
func GetIdFunc(r *http.Request) (string, error) {
	var id = chi.URLParam(r, "id")

	if id == "" {
		return "", rest.ErrIdNotPresent
	}

	return id, nil
}

// Register a series of resource definitions on a chi router
func Register(r chi.Router, db repository.Repositorium, def []rest.ResourceDef) {
	for d := range def {
		var p = rest.Resource{
			IdGetter: GetIdFunc,
			Fields:   def[d].Fields,
			Entity:   def[d].Entity,
			Repo:     db,
		}

		var n string
		if def[d].Path == "" {
			n = "/" + strings.ToLower(def[d].Entity)
		} else {
			n = "/" + def[d].Path
		}

		if def[d].Methods.Has(rest.ENDPOINTS_GET) {
			r.Get(n+"/{id}", p.Get)
		}

		if def[d].Methods.Has(rest.ENDPOINTS_REPLACE) {
			r.Put(n+"/{id}", p.Replace)
		}

		if def[d].Methods.Has(rest.ENDPOINTS_UPDATE) {
			r.Patch(n+"/{id}", p.Update)
		}

		if def[d].Methods.Has(rest.ENDPOINTS_DELETE) {
			r.Delete(n+"/{id}", p.Delete)
		}

		if def[d].Methods.Has(rest.ENDPOINTS_CREATE) {
			r.Post(n, p.Create)
		}

		if def[d].Methods.Has(rest.ENDPOINTS_LIST) {
			r.Get(n, p.List)
		}
	}
}
