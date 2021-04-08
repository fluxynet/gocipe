package mongo

import (
	"go.mongodb.org/mongo-driver/bson"

	"github.com/fluxynet/gocipe/repository"
	"github.com/fluxynet/gocipe/values"
)

// ConditionsToBsonD returns the filter for mongo.
func ConditionsToBsonD(c []repository.Condition) (bson.D, error) {
	var n = len(c)
	if n == 0 {
		return nil, nil
	}

	var filters = make(bson.D, n)

	for i := range c {
		var (
			name = c[i].Attribute
			val  = c[i].Value
		)

		if name == "" {
			return nil, repository.ErrInvalidAttribute
		}

		switch c[i].Operator {
		default:
			return nil, repository.ErrInvalidConditionOperator
		case repository.Equals:
			filters[i] = bson.E{
				Key:   name,
				Value: bson.M{"$eq": val},
			}
		case repository.NotEquals:
			filters[i] = bson.E{
				Key:   name,
				Value: bson.M{"$ne": val},
			}
		case repository.GreaterThan:
			filters[i] = bson.E{
				Key:   name,
				Value: bson.M{"$gt": val},
			}
		case repository.GreaterOrEqualTo:
			filters[i] = bson.E{
				Key:   name,
				Value: bson.M{"$gte": val},
			}
		case repository.LessThan:
			filters[i] = bson.E{
				Key:   name,
				Value: bson.M{"$lt": val},
			}
		case repository.LessOrEqualTo:
			filters[i] = bson.E{
				Key:   name,
				Value: bson.M{"$lte": val},
			}
		case repository.Like:
			filters[i] = bson.E{
				Key:   name,
				Value: bson.M{"$regex": val},
			}
		case repository.In:
			filters[i] = bson.E{
				Key:   name,
				Value: bson.M{"$in": val},
			}
		case repository.NotIn:
			filters[i] = bson.E{
				Key:   name,
				Value: bson.M{"$nin": val},
			}
		}
	}

	return filters, nil
}

// ValuesToBsonM converts values to Bson that can be used for insert
func ValuesToBsonM(vals *values.Values) bson.M {
	if vals == nil {
		return nil
	}

	var (
		data = make(bson.M, vals.Length())
		it   = vals.Iterator()
	)

	for it.Next() {
		var v = it.Value()

		data[v.Name] = v.Value
	}

	return data
}
