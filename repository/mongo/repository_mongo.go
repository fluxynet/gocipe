package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/fluxynet/gocipe/fields"
	"github.com/fluxynet/gocipe/repository"
	"github.com/fluxynet/gocipe/values"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

// Get a single Entity by id
func (r Repo) Get(ctx context.Context, entity string, f fields.Fields, id string) (*values.Values, error) {
	var (
		vals  values.Values
		datum = f.GetEmptyValues()
	)

	var oid, err = primitive.ObjectIDFromHex(id)
	if err == nil {
		err = r.db.Collection(entity).FindOne(ctx, bson.M{"_id": oid}).Decode(&datum)
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

// List multiple Entity with pagination rules and conditions
func (r *Repo) List(ctx context.Context, entity string, f fields.Fields, p repository.Pagination, c ...repository.Condition) ([]values.Values, error) {
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
			vals.Set("id", v.Value)
			vals.Unset("_id")
		}

		data = append(data, *vals)
	}

	return data, err
}

// Delete a single Entity by id
func (r *Repo) Delete(ctx context.Context, entity, id string) error {
	var oid, err = primitive.ObjectIDFromHex(id)

	var res *mongo.DeleteResult
	if err == nil {
		res, err = r.db.Collection(entity).DeleteOne(ctx, bson.M{"_id": oid})
	}

	if err == nil && res.DeletedCount == 0 {
		err = repository.ErrNotFound
	}

	return err
}

// DeleteWhere delete multiple Entity based on conditions
func (r *Repo) DeleteWhere(ctx context.Context, entity string, c ...repository.Condition) error {
	var (
		err     error
		filters bson.D
	)

	filters, err = ConditionsToBsonD(c)
	if err != nil {
		return err
	}

	_, err = r.db.Collection(entity).DeleteMany(ctx, filters)

	return err
}

// Create a new Entity in persistent storage
func (r *Repo) Create(ctx context.Context, entity string, vals *values.Values) (string, error) {
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
		id = v.Hex()
	}

	return id, err
}

// Update an existing Entity in persistent storage
func (r *Repo) Update(ctx context.Context, entity string, id string, vals *values.Values) error {
	var (
		result *mongo.UpdateResult
		data   = ValuesToBsonM(vals)
	)

	vals.Unset("id") // mongo has its own representation of id as _id

	var oid, err = primitive.ObjectIDFromHex(id)
	if err == nil {
		result, err = r.db.Collection(entity).UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": data})
	}

	if err == nil && result.MatchedCount == 0 {
		err = repository.ErrNotFound
	}

	return err
}

// UpdateValuesWhere Values in persistent storage
func (r *Repo) UpdateWhere(ctx context.Context, entity string, vals *values.Values, c ...repository.Condition) error {
	var (
		err     error
		filters bson.D
		data    = ValuesToBsonM(vals)
	)

	filters, err = ConditionsToBsonD(c)
	if err != nil {
		return err
	}

	_, err = r.db.Collection(entity).UpdateMany(ctx, filters, data)

	return err
}

// Close db connection
func (r *Repo) Close() error {
	if r.db == nil {
		return nil
	}

	return r.cli.Disconnect(nil)
}
