package mongo

import (
	"context"
	"time"

	"github.com/fluxynet/gocipe/fields"
	"github.com/fluxynet/gocipe/repository"
	"github.com/fluxynet/gocipe/values"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	var _ repository.Repositorium = &Repo{}
}

type Repo struct {
	db *mongo.Database
}

// New repo with mongo
func New(dbname, uri string) (repository.Repositorium, error) {
	var client, err = mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	var ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	var repo = &Repo{db: client.Database(dbname)}

	return repo, err
}

// Get a single Entity by id
func (r Repo) Get(ctx context.Context, entity string, f fields.Fields, id string) (values.Values, error) {
	var (
		err   error
		vals  values.Values
		datum = f.GetEmptyValues()
	)

	err = r.db.Collection(entity).FindOne(ctx, bson.D{{
		Key:   "_id",
		Value: id,
	}}).Decode(&datum)

	if err == mongo.ErrNoDocuments {
		err = repository.ErrNotFound
	}

	if err != nil {
		return vals, err
	}

	vals.FromMap(datum)
	if v := vals.Get("_id"); v != nil {
		vals.Set("id", v)
		vals.Unset("_id")
	}

	return vals, nil
}

// List multiple Entity with pagination rules and conditions
func (r *Repo) List(ctx context.Context, entity string, f fields.Fields, p repository.Pagination, c ...repository.Condition) ([]values.Values, error) {
	var (
		err     error
		data    []values.Values
		cursor  *mongo.Cursor
		filters = ConditionsToBsonD(c)
	)

	// todo(fx) pagination

	cursor, err = r.db.Collection(entity).Find(ctx, filters)
	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		var datum = f.GetEmptyValues()
		err = cursor.Decode(datum)

		if err != nil {
			return nil, err
		}

		var vals = values.FromMap(datum)
		if v := vals.Get("_id"); v != nil {
			vals.Set("id", v)
			vals.Unset("_id")
		}
		data = append(data, vals)
	}

	return nil, err
}

// Delete a single Entity by id
func (r *Repo) Delete(ctx context.Context, entity, id string) error {
	result, err := r.db.Collection(entity).DeleteOne(ctx, bson.D{{
		Key:   "_id",
		Value: id,
	}})

	if err == nil && result.DeletedCount == 0 {
		err = repository.ErrNotFound
	}

	return err
}

// DeleteWhere delete multiple Entity based on conditions
func (r *Repo) DeleteWhere(ctx context.Context, entity string, c ...repository.Condition) error {
	var (
		err     error
		filters = ConditionsToBsonD(c)
	)

	_, err = r.db.Collection(entity).DeleteMany(ctx, filters)

	return err
}

// Create a new Entity in persistent storage
func (r *Repo) Create(ctx context.Context, entity string, vals values.Values) (string, error) {
	var (
		data bson.M
		rs   *mongo.InsertOneResult
		err  error
		id   string
	)

	data = ValuesToBsonM(vals)
	rs, err = r.db.Collection(entity).InsertOne(ctx, data)

	if err != nil {
		return id, err
	}

	if v, ok := rs.InsertedID.(primitive.ObjectID); ok {
		id = v.String()
	}

	return id, err
}

// Update an existing Entity in persistent storage
func (r *Repo) Update(ctx context.Context, entity string, id string, vals values.Values) error {
	var (
		err    error
		result *mongo.UpdateResult
		data   = ValuesToBsonM(vals)
	)

	vals.Unset("id") // mongo has its own representation of id as _id

	result, err = r.db.Collection(entity).UpdateOne(ctx, bson.D{{
		Key:   "_id",
		Value: id,
	}}, data)

	if err == nil && result.UpsertedCount == 0 {
		err = repository.ErrNotFound
	}

	return err
}

// UpdateValuesWhere Values in persistent storage
func (r *Repo) UpdateValues(ctx context.Context, entity string, vals values.Values, c ...repository.Condition) error {
	var (
		err error

		filters = ConditionsToBsonD(c)
		data    = ValuesToBsonM(vals)
	)

	_, err = r.db.Collection(entity).UpdateMany(ctx, filters, data)

	return err
}
