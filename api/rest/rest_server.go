package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/fluxynet/gocipe/api"
	"github.com/fluxynet/gocipe/repository"
	"github.com/fluxynet/gocipe/types/fields/entity"
	"github.com/fluxynet/gocipe/util"
	"github.com/fluxynet/gocipe/values"
)

// Server represents a REST server endpoint set for a specific entity
type Server struct {
	// IdGetter to read ids from http.Request
	IdGetter GetIdFunc

	// Entity represented
	Entity entity.Entity

	// Repo for data persistence
	Repo repository.Repositorium

	// Actions enabled
	Actions api.ActionSet
}

// ServeHTTP is a simple muxer based on method (and also presence of id in url in case of GET)
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		disabled bool
		handler  http.HandlerFunc
	)

	switch r.Method {
	case http.MethodHead, http.MethodGet:
		if _, err := s.IdGetter(r); err == ErrIdNotPresent {
			disabled = s.Actions.NotHas(api.ActionList)
			handler = s.List
		} else {
			disabled = s.Actions.NotHas(api.ActionRead)
			handler = s.Get
		}
	case http.MethodPost:
		disabled = s.Actions.NotHas(api.ActionCreate)
		handler = s.Create
	case http.MethodPut:
		disabled = s.Actions.NotHas(api.ActionReplace)
		handler = s.Replace
	case http.MethodPatch:
		disabled = s.Actions.NotHas(api.ActionUpdate)
		handler = s.Update
	case http.MethodDelete:
		disabled = s.Actions.NotHas(api.ActionDelete)
		handler = s.Delete
	case http.MethodOptions:
		// todo
	}

	if handler == nil || disabled {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	handler(w, r)
	// todo check accepted types
}

func (s *Server) Get(w http.ResponseWriter, r *http.Request) {
	var (
		b []byte

		id     string
		err    error
		vals   *values.Values
		status = http.StatusOK
		ctx    = r.Context()
	)

	id, err = s.IdGetter(r)
	if err == ErrIdNotPresent {
		status = http.StatusBadRequest
	}

	if err == nil {
		vals, err = s.Repo.Get(ctx, s.Entity, id)
	}

	if err == repository.ErrNotFound {
		status = http.StatusNotFound
	} else if err != nil {
		status = http.StatusInternalServerError
	} else {
		var data = vals.ToMap()
		b, err = json.Marshal(data)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if status == http.StatusOK {
		w.Write(b)
	}
}

func (s *Server) List(w http.ResponseWriter, r *http.Request) {
	var (
		b    []byte
		err  error
		c    []repository.Condition
		vals []values.Values

		q = r.URL.Query()
		p = repository.Pagination{} // todo
		f = s.Entity.Fields()

		status = http.StatusOK
		ctx    = r.Context()
	)

	c, err = repository.ConditionsFromMap(q, f)
	if err != nil {
		status = http.StatusBadRequest
		err = fmt.Errorf("filters could not be parsed. %w", err)
	}

	if err == nil {
		p.Offset, err = util.GetSingleInteger(q, "__offset")
		if err != nil {
			status = http.StatusBadRequest
			err = fmt.Errorf("offset parameter could not be parsed. %w", err)
		}
	}

	if err == nil {
		p.Limit, err = util.GetSingleInteger(q, "__limit")
		if err != nil {
			status = http.StatusBadRequest
			err = fmt.Errorf("limit parameter could not be parsed. %w", err)
		}
	}

	if err != nil {
		// sad
	} else if v, ok := q["__sort"]; !ok {
		// not present
	} else if len(v) != 1 {
		status = http.StatusBadRequest
		err = fmt.Errorf("multiple sort values obtained. %w", api.ErrInvalidRequestParameters)
	} else {
		p.Order, err = repository.OrderByFromString(v[0], f)
		if err != nil {
			status = http.StatusBadRequest
			err = fmt.Errorf("sort parameter could not be parsed. %w", api.ErrInvalidRequestParameters)
		}
	}

	if err == nil {
		vals, err = s.Repo.List(ctx, s.Entity, p, c...)
	}

	if err == repository.ErrNotFound {
		status = http.StatusNotFound
	} else if err != nil {
		status = http.StatusInternalServerError
	} else {
		var data = make([]map[string]interface{}, len(vals))
		for i := range vals {
			data[i] = vals[i].ToMap()
		}
		b, err = json.Marshal(data)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if status == http.StatusOK {
		w.Write(b)
	} else {
		w.Write([]byte(`{"error": "` + strings.ReplaceAll(err.Error(), `"`, `\"`) + `"}`))
	}
}

func (s *Server) Delete(w http.ResponseWriter, r *http.Request) {
	var (
		id     string
		err    error
		status = http.StatusOK
		ctx    = r.Context()
	)

	id, err = s.IdGetter(r)
	if err == ErrIdNotPresent {
		status = http.StatusBadRequest
	}

	if err == nil {
		err = s.Repo.Delete(ctx, s.Entity, id)
	}

	if err == repository.ErrNotFound {
		status = http.StatusNotFound
	} else if err != nil {
		status = http.StatusInternalServerError
	}

	w.WriteHeader(status)
}

func (s *Server) Create(w http.ResponseWriter, r *http.Request) {
	var (
		vals   *values.Values
		err    error
		id     string
		status = http.StatusOK
		ctx    = r.Context()
	)

	vals, err = values.FromJSON(r.Body, s.Entity.Fields(), false)
	defer util.Closed(r.Body, &err)

	if err == nil {
		id, err = s.Repo.Create(ctx, s.Entity, vals)
	}

	if err == nil {
		vals.Set("id", id)
	}

	if err != nil {
		status = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(nil)
}

// Replace an entity by another instance of itself; ids cannot be updated.
func (s *Server) Replace(w http.ResponseWriter, r *http.Request) {
	var (
		vals   *values.Values
		err    error
		id     string
		status = http.StatusOK
		ctx    = r.Context()
	)

	id, err = s.IdGetter(r)
	if err == ErrIdNotPresent {
		status = http.StatusBadRequest
	}

	vals, err = s.Repo.Get(ctx, s.Entity, id)
	if err == repository.ErrNotFound {
		status = http.StatusNotFound
	}

	if err == nil {
		vals, err = values.FromJSON(r.Body, s.Entity.Fields(), false)
		defer util.Closed(r.Body, &err)
	}

	if err == nil {
		vals.Set("id", id)
		err = s.Repo.Update(ctx, s.Entity, id, vals)
	}

	if err == repository.ErrNotFound {
		status = http.StatusNotFound
	}

	if err != nil {
		status = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(nil)
}

// Update is partial update of an entity, typically Patch; ids cannot be updated
func (s *Server) Update(w http.ResponseWriter, r *http.Request) {
	var (
		vals   *values.Values
		err    error
		id     string
		status = http.StatusOK
		ctx    = r.Context()
	)

	id, err = s.IdGetter(r)
	if err == ErrIdNotPresent {
		status = http.StatusBadRequest
	}

	vals, err = s.Repo.Get(ctx, s.Entity, id)
	if err == repository.ErrNotFound {
		status = http.StatusNotFound
	}

	if err == nil {
		vals, err = values.FromJSON(r.Body, s.Entity.Fields(), true)
		defer util.Closed(r.Body, &err)
	}

	if err == nil {
		vals.Set("id", id)
		err = s.Repo.Update(ctx, s.Entity, id, vals)
	}

	if err == repository.ErrNotFound {
		status = http.StatusNotFound
	} else if err != nil {
		status = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(nil)
}
