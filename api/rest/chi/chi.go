package chi

import (
	"net/http"

	"github.com/fluxynet/gocipe/api"
	"github.com/fluxynet/gocipe/api/rest"
	"github.com/fluxynet/gocipe/repository"

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
func Register(r chi.Router, db repository.Repositorium, res api.Resource) {
	var p = rest.Server{
		IdGetter: GetIdFunc,
		Entity:   res,
		Repo:     db,
	}

	var (
		n       = res.Path()
		actions = res.Actions()
	)

	if actions.Has(api.ActionRead) {
		r.Get(n+"/{id}", p.Get)
	}

	if actions.Has(api.ActionReplace) {
		r.Put(n+"/{id}", p.Replace)
	}

	if actions.Has(api.ActionUpdate) {
		r.Patch(n+"/{id}", p.Update)
	}

	if actions.Has(api.ActionDelete) {
		r.Delete(n+"/{id}", p.Delete)
	}

	if actions.Has(api.ActionCreate) {
		r.Post(n, p.Create)
	}

	if actions.Has(api.ActionList) {
		r.Get(n, p.List)
	}
}
