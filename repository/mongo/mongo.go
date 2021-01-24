package mongo

import (
	"github.com/fluxynet/gocipe/repository"
	"github.com/fluxynet/gocipe/values"
	"go.mongodb.org/mongo-driver/bson"
)

// ConditionsToBsonD returns the filter for mongo.
func ConditionsToBsonD(c []repository.Condition) bson.D {
	var (
		n       = len(c)
		filters = make(bson.D, n)
	)

	for i := range c {
		var (
			name = c[i].Property
			val  = c[i].Value
		)

		switch c[i].Operator {
		case repository.Equals:
			filters[i] = bson.E{
				Key:   name,
				Value: val,
			}
		case repository.NotEquals:
			filters[i] = bson.E{
				Key: name,
				Value: bson.E{
					Key:   "$ne",
					Value: val,
				},
			}
		case repository.GreaterThan:
			filters[i] = bson.E{
				Key: name,
				Value: bson.E{
					Key:   "$gt",
					Value: val,
				},
			}
		case repository.GreaterOrEqualTo:
			filters[i] = bson.E{
				Key: name,
				Value: bson.E{
					Key:   "$gte",
					Value: val,
				},
			}
		case repository.LessThan:
			filters[i] = bson.E{
				Key: name,
				Value: bson.E{
					Key:   "$lt",
					Value: val,
				},
			}
		case repository.LessOrEqualTo:
			filters[i] = bson.E{
				Key: name,
				Value: bson.E{
					Key:   "$lte",
					Value: val,
				},
			}
		case repository.In:
			filters[i] = bson.E{
				Key: name,
				Value: bson.E{
					Key:   "$in",
					Value: val,
				},
			}
		case repository.NotIn:
			filters[i] = bson.E{
				Key: name,
				Value: bson.E{
					Key:   "$nin",
					Value: val,
				},
			}
		}
	}

	return filters
}

// ValuesToBsonM converts values to Bson that can be used for insert
func ValuesToBsonM(vals values.Values) bson.M {
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
