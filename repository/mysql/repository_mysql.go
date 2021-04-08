package mysql

import (
	"context"
	"database/sql"

	"github.com/google/uuid"

	"github.com/fluxynet/gocipe/repository"
	"github.com/fluxynet/gocipe/types/fields/entity"
	"github.com/fluxynet/gocipe/util"
	"github.com/fluxynet/gocipe/values"
)

func init() {
	var _ repository.Repositorium = &Repo{}
}

// EntityRepo is an implementation of EntityRepository to allow persistence of Name
type Repo struct {
	db *sql.DB
}

func New(db *sql.DB) Repo {
	return Repo{db}
}

// Get a single Name by id
func (r *Repo) Get(ctx context.Context, entity entity.Entity, id string) (*values.Values, error) {
	var (
		vals values.Values
		q    = Get(entity, id)
		f    = entity.Fields()
		it   = f.Iterator()
		dst  = GetScanDest(f)
	)

	var rs, err = r.db.QueryContext(ctx, q.SQL, q.Args...)
	defer util.Closed(rs, &err)

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

// List multiple Name with pagination rules and conditions
func (r *Repo) List(ctx context.Context, entity entity.Entity, p repository.Pagination, c ...repository.Condition) ([]values.Values, error) {
	var (
		l []values.Values
		q = List(entity, p, c...)
		f = entity.Fields()
	)

	var rs, err = r.db.QueryContext(ctx, q.SQL, q.Args...)
	defer util.Closed(rs, &err)

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

// Delete a single Name by id
func (r *Repo) Delete(ctx context.Context, named repository.Named, id string) error {
	var q = Delete(named, id)
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

// DeleteWhere delete multiple Name based on conditions
func (r *Repo) DeleteWhere(ctx context.Context, named repository.Named, c ...repository.Condition) error {
	var q = DeleteWhere(named, c...)
	var _, err = r.db.ExecContext(ctx, q.SQL, q.Args...)
	return err
}

// Create a new Name in persistent storage
func (r *Repo) Create(ctx context.Context, named repository.Named, vals *values.Values) (string, error) {
	var (
		q   Query
		err error
		id  string
	)

	id = uuid.NewString()
	vals.Set("id", id)

	q = Create(named, vals)

	_, err = r.db.ExecContext(ctx, q.SQL, q.Args...)

	return id, err
}

// Update an existing Name in persistent storage
func (r *Repo) Update(ctx context.Context, named repository.Named, id string, vals *values.Values) error {
	vals.Unset("id")

	var q = Update(
		named,
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
func (r *Repo) UpdateWhere(ctx context.Context, named repository.Named, vals *values.Values, c ...repository.Condition) error {
	var q = UpdateWhere(
		named,
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
