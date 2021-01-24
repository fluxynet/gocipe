package rest

import (
	"encoding/json"
	"github.com/fluxynet/gocipe"
	"github.com/fluxynet/gocipe/fields"
	"github.com/fluxynet/gocipe/repository"
	"github.com/fluxynet/gocipe/values"
	"net/http"
)

// Server represents a REST server
type Server struct {
	// IdGetter to read ids from http.Request
	IdGetter GetIdFunc

	// Fields for the Entity represented
	Fields fields.Fields

	// Entity name
	Entity string

	// Repo for data persistence
	Repo repository.Repositorium

	// NoGet disables Get method on the resource if set to true
	NoGet bool

	// NoList disables List method on the resource if set to true
	NoList bool

	// NoDelete disables Delete method on the resource if set to true
	NoDelete bool

	// NoCreate disables Create method on the resource if set to true
	NoCreate bool

	// NoReplace disables Replace method on the resource if set to true
	NoReplace bool

	// NoUpdate disables Update method on the resource if set to true
	NoUpdate bool
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
			disabled = s.NoList
			handler = s.List
		} else {
			disabled = s.NoGet
			handler = s.Get
		}
	case http.MethodPost:
		disabled = s.NoCreate
		handler = s.Create
	case http.MethodPut:
		disabled = s.NoReplace
		handler = s.Replace
	case http.MethodPatch:
		disabled = s.NoUpdate
		handler = s.Update
	case http.MethodDelete:
		disabled = s.NoDelete
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
		vals   values.Values
		status = http.StatusOK
		ctx    = r.Context()
	)

	id, err = s.IdGetter(r)
	if err == ErrIdNotPresent {
		status = http.StatusBadRequest
	}

	if err == nil {
		vals, err = s.Repo.Get(ctx, s.Entity, s.Fields, id)
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
		b []byte

		p = repository.Pagination{}  // todo
		c = []repository.Condition{} // todo

		status    = http.StatusOK
		ctx       = r.Context()
		vals, err = s.Repo.List(ctx, s.Entity, s.Fields, p, c...)
	)

	if err == repository.ErrNotFound {
		status = http.StatusNotFound
	} else if err != nil {
		//
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
		vals   values.Values
		err    error
		id     string
		b      []byte
		status = http.StatusOK
		ctx    = r.Context()
	)

	vals, err = values.FromJSON(r.Body, s.Fields)
	defer gocipe.Closed(r.Body, &err)

	if err == nil {
		id, err = s.Repo.Create(ctx, s.Entity, vals)
	}

	if err == nil {
		vals.Set("id", id)
		var data = vals.ToMap()
		b, err = json.Marshal(data)
	}

	if err != nil {
		status = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if status == http.StatusOK {
		_, err = w.Write(b)
	}
}

// Replace an entity by another instance of itself; ids cannot be updated.
func (s *Server) Replace(w http.ResponseWriter, r *http.Request) {
	var (
		vals   values.Values
		err    error
		id     string
		b      []byte
		status = http.StatusOK
		ctx    = r.Context()
	)

	id, err = s.IdGetter(r)
	if err == ErrIdNotPresent {
		status = http.StatusBadRequest
	}

	vals, err = s.Repo.Get(ctx, s.Entity, s.Fields, id)
	if err == repository.ErrNotFound {
		status = http.StatusNotFound
	}

	if err == nil {
		vals, err = values.FromJSON(r.Body, s.Fields)
		defer gocipe.Closed(r.Body, &err)
	}

	if err == nil {
		vals.Set("id", id)
		err = s.Repo.Update(ctx, s.Entity, id, vals)
	}

	if err == repository.ErrNotFound {
		status = http.StatusNotFound
	} else if err == nil {
		vals.Set("id", id)
		var data = vals.ToMap()
		b, err = json.Marshal(data)
	}

	if err != nil {
		status = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if status == http.StatusOK {
		_, err = w.Write(b)
	}
}

// Update is partial update of an entity, typically Patch; ids cannot be updated
func (s *Server) Update(w http.ResponseWriter, r *http.Request) {
	var (
		vals   values.Values
		err    error
		id     string
		b      []byte
		status = http.StatusOK
		ctx    = r.Context()
	)

	id, err = s.IdGetter(r)
	if err == ErrIdNotPresent {
		status = http.StatusBadRequest
	}

	vals, err = s.Repo.Get(ctx, s.Entity, s.Fields, id)
	if err == repository.ErrNotFound {
		status = http.StatusNotFound
	}

	if err == nil {
		vals, err = values.FromJSON(r.Body, s.Fields)
		defer gocipe.Closed(r.Body, &err)
	}

	if err == nil {
		vals.Set("id", id)
		err = s.Repo.Update(ctx, s.Entity, id, vals)
	}

	if err == repository.ErrNotFound {
		status = http.StatusNotFound
	} else if err == nil {
		vals.Set("id", id)
		var data = vals.ToMap()
		b, err = json.Marshal(data)
	}

	if err != nil {
		status = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if status == http.StatusOK {
		_, err = w.Write(b)
	}
}
