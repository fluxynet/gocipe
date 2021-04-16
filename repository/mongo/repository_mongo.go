package mongo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/fluxynet/gocipe/repository"
	"github.com/fluxynet/gocipe/types/fields/entity"
	"github.com/fluxynet/gocipe/values"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	// ErrInvalidID is when an invalid ID is passed
	ErrInvalidID = errors.New("invalid id")
)

func init() {
	var _ repository.Repositorium = &Repo{}
}

type Repo struct {
	db  *mongo.Database
	cli *mongo.Client
}

// New repo with mongo
func New(dbname, uri string) (repository.Repositorium, error) {
	var client, err = mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)

	if err != nil {
		return nil, err
	}

	var repo = &Repo{db: client.Database(dbname), cli: client}

	return repo, err
}

// Get a single Name by id
func (r Repo) Get(ctx context.Context, entity entity.Entity, id string) (*values.Values, error) {
	var (
		vals  values.Values
		datum = entity.Fields().GetEmptyValues()
	)

	var oid, err = primitive.ObjectIDFromHex(id)
	if err == nil {
		err = r.db.Collection(entity.Name()).FindOne(ctx, bson.M{"_id": oid}).Decode(&datum)
	} else {
		return nil, ErrInvalidID
	}

	if err == mongo.ErrNoDocuments {
		err = repository.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	vals.FromMap(datum)
	if v := vals.Get("_id"); v != nil {
		vals.Set("id", v.Value)
		vals.Unset("_id")
	}

	return &vals, nil
}

// List multiple Name with pagination rules and conditions
func (r *Repo) List(ctx context.Context, entity entity.Entity, p repository.Pagination, c ...repository.Condition) ([]values.Values, error) {
	var (
		err     error
		data    []values.Values
		cursor  *mongo.Cursor
		filters bson.D
	)

	filters, err = ConditionsToBsonD(c)
	if err != nil {
		return nil, err
	}

	// todo(fx) pagination

	var f = entity.Fields()

	cursor, err = r.db.Collection(entity.Name()).Find(ctx, filters)
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
			vals.Set("id", v.Value)
			vals.Unset("_id")
		}

		data = append(data, *vals)
	}

	return data, err
}

// Delete a single Name by id
func (r *Repo) Delete(ctx context.Context, named repository.Named, id string) error {
	var oid, err = primitive.ObjectIDFromHex(id)

	var res *mongo.DeleteResult
	if err == nil {
		res, err = r.db.Collection(named.Name()).DeleteOne(ctx, bson.M{"_id": oid})
	}

	if err == nil && res.DeletedCount == 0 {
		err = repository.ErrNotFound
	}

	return err
}

// DeleteWhere delete multiple Name based on conditions
func (r *Repo) DeleteWhere(ctx context.Context, named repository.Named, c ...repository.Condition) error {
	var (
		err     error
		filters bson.D
	)

	filters, err = ConditionsToBsonD(c)
	if err != nil {
		return err
	}

	_, err = r.db.Collection(named.Name()).DeleteMany(ctx, filters)

	return err
}

// Create a new Name in persistent storage
func (r *Repo) Create(ctx context.Context, named repository.Named, vals *values.Values) (string, error) {
	var (
		data bson.M
		rs   *mongo.InsertOneResult
		err  error
		id   string
	)

	if v := vals.Get("id"); v != nil && v.IsString() {
		fmt.Println("id === ", v.String())
		var i, e = primitive.ObjectIDFromHex(v.String())
		if e != nil {
			return "", e
		}

		vals.Set("_id", i)
		vals.Unset("id")
	}

	data = ValuesToBsonM(vals)
	rs, err = r.db.Collection(named.Name()).InsertOne(ctx, data)

	if err != nil {
		return id, err
	}

	if v, ok := rs.InsertedID.(primitive.ObjectID); ok {
		id = v.Hex()
	}

	return id, err
}

// Update an existing Name in persistent storage
func (r *Repo) Update(ctx context.Context, named repository.Named, id string, vals *values.Values) error {
	var (
		result *mongo.UpdateResult
		data   = ValuesToBsonM(vals)
	)

	vals.Unset("id") // mongo has its own representation of id as _id

	var oid, err = primitive.ObjectIDFromHex(id)
	if err == nil {
		result, err = r.db.Collection(named.Name()).UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": data})
	}

	if err == nil && result.MatchedCount == 0 {
		err = repository.ErrNotFound
	}

	return err
}

// UpdateValuesWhere Values in persistent storage
func (r *Repo) UpdateWhere(ctx context.Context, named repository.Named, vals *values.Values, c ...repository.Condition) error {
	var (
		err     error
		filters bson.D
		data    = ValuesToBsonM(vals)
	)

	filters, err = ConditionsToBsonD(c)
	if err != nil {
		return err
	}

	_, err = r.db.Collection(named.Name()).UpdateMany(ctx, filters, data)

	return err
}

// Close db connection
func (r *Repo) Close() error {
	if r.db == nil {
		return nil
	}

	return r.cli.Disconnect(nil)
}
