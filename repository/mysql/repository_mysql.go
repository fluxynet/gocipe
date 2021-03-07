package mysql

import (
	"context"
	"database/sql"

	"github.com/google/uuid"

	"github.com/fluxynet/gocipe"
	"github.com/fluxynet/gocipe/fields"
	"github.com/fluxynet/gocipe/repository"
	"github.com/fluxynet/gocipe/values"
)

func init() {
	var _ repository.Repositorium = &Repo{}
}

// EntityRepo is an implementation of EntityRepository to allow persistence of Entity
type Repo struct {
	db *sql.DB
}

func New(db *sql.DB) Repo {
	return Repo{db}
}

// Get a single Entity by id
func (r *Repo) Get(ctx context.Context, entity string, f fields.Fields, id string) (*values.Values, error) {
	var (
		vals values.Values
		q    = Get(entity, f, id)
		it   = f.Iterator()
		dst  = GetScanDest(f)
	)

	var rs, err = r.db.QueryContext(ctx, q.SQL, q.Args...)
	defer gocipe.Closed(rs, &err)

	if err == sql.ErrNoRows {
		err = repository.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	for rs.Next() {
		err = rs.Scan(dst...)
		if err != nil {
			return nil, err
		}

		for i := 0; it.Next(); i++ {
			vals.Set(it.Field().Name, dst[i])
		}
	}

	return &vals, err
}

// List multiple Entity with pagination rules and conditions
func (r *Repo) List(ctx context.Context, entity string, f fields.Fields, p repository.Pagination, c ...repository.Condition) ([]values.Values, error) {
	var (
		l []values.Values
		q = List(entity, f, p, c...)
	)

	var rs, err = r.db.QueryContext(ctx, q.SQL, q.Args...)
	defer gocipe.Closed(rs, &err)

	if err != nil {
		return nil, err
	}

	for it := f.Iterator(); rs.Next(); it.Reset() {
		var (
			vals values.Values
			dst  = GetScanDest(f)
		)

		err = rs.Scan(dst...)
		if err != nil {
			return nil, err
		}

		for i := 0; it.Next(); i++ {
			vals.Set(it.Field().Name, dst[i])
		}

		l = append(l, vals)
	}

	return l, err
}

// Delete a single Entity by id
func (r *Repo) Delete(ctx context.Context, entity, id string) error {
	var q = Delete(entity, id)
	var res, err = r.db.ExecContext(ctx, q.SQL, q.Args...)
	var n int64

	if err == nil {
		n, err = res.RowsAffected()
	}

	if err == nil && n == 0 {
		err = repository.ErrNotFound
	}

	return err
}

// DeleteWhere delete multiple Entity based on conditions
func (r *Repo) DeleteWhere(ctx context.Context, entity string, c ...repository.Condition) error {
	var q = DeleteWhere(entity, c...)
	var _, err = r.db.ExecContext(ctx, q.SQL, q.Args...)
	return err
}

// Create a new Entity in persistent storage
func (r *Repo) Create(ctx context.Context, entity string, vals *values.Values) (string, error) {
	var (
		q   Query
		err error
		id  string
	)

	id = uuid.NewString()
	vals.Set("id", id)

	q = Create(entity, vals)

	_, err = r.db.ExecContext(ctx, q.SQL, q.Args...)

	return id, err
}

// Update an existing Entity in persistent storage
func (r *Repo) Update(ctx context.Context, entity string, id string, vals *values.Values) error {
	vals.Unset("id")

	var q = Update(
		entity,
		id,
		vals,
	)

	var res, err = r.db.ExecContext(ctx, q.SQL, q.Args...)
	var n int64

	if err == nil {
		n, err = res.RowsAffected()
	}

	if err == nil && n == 0 {
		err = repository.ErrNotFound
	}

	return err
}

// UpdateValuesWhere Values in persistent storage
func (r *Repo) UpdateWhere(ctx context.Context, entity string, vals *values.Values, c ...repository.Condition) error {
	var q = UpdateWhere(
		entity,
		vals,
		c...,
	)

	var _, err = r.db.ExecContext(ctx, q.SQL, q.Args...)

	return err
}

// Close db connection
func (r *Repo) Close() error {
	if r.db == nil {
		return nil
	}

	return r.db.Close()
}
