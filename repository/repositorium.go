package repository

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/fluxynet/gocipe"
	"github.com/fluxynet/gocipe/fields"
	"github.com/fluxynet/gocipe/values"
)

var (
	// ErrNotFound is when the item was not found
	ErrNotFound = errors.New("item not found")

	// ErrUnknownSortAttribute when an unknown sort attribute is passed
	ErrUnknownSortAttribute = errors.New("unknown sort attribute")

	// ErrInvalidConditionOperator when an invalid conditional operator is used
	ErrInvalidConditionOperator = errors.New("invalid conditional operator")

	// ErrInvalidAttribute when an invalid attribute is passed
	ErrInvalidAttribute = errors.New("invalid conditional attribute")
)

// ConditionOperator represents the condition wrt the value
type ConditionOperator uint8

const (
	// Equals denotes Equal
	Equals = ConditionOperator(0)

	// NotEquals denotes Not Equal
	NotEquals = ConditionOperator(1)

	// GreaterThan denotes Greater Than
	GreaterThan = ConditionOperator(2)

	// GreaterOrEqualTo denotes Greater or Equal To
	GreaterOrEqualTo = ConditionOperator(3)

	// LessThan denotes Less than
	LessThan = ConditionOperator(4)

	// LessOrEqualTo denotes Less than or Equal to
	LessOrEqualTo = ConditionOperator(5)

	// Like for strings that match a term
	Like = ConditionOperator(6)

	// In denotes In
	In = ConditionOperator(7)

	// NotIn denotes Not in
	NotIn = ConditionOperator(8)
)

// String representation (mainly for testing / debugging / logging)m
func (c ConditionOperator) String() string {
	switch c {
	case Equals:
		return "="
	case NotEquals:
		return "!="
	case GreaterThan:
		return ">"
	case GreaterOrEqualTo:
		return ">="
	case LessThan:
		return "<"
	case LessOrEqualTo:
		return "<="
	case Like:
		return "~"
	case In:
		return "IN"
	case NotIn:
		return "NOT IN"
	}

	return "?? " + strconv.Itoa(int(c))
}

// ConditionType represents what kind of condition it is with respect to other conditions (AND / OR)
type ConditionType uint8

const (
	// And means the condition must be inclusive wrt other conditions
	And = ConditionType(0)

	// Or means the condition needs not be inclusive wrt other conditions
	Or = ConditionType(1)
)

// Condition represents filter criteria for fetching multiple elements
type Condition struct {
	Attribute string
	Operator  ConditionOperator
	Value     interface{}
	Type      ConditionType
}

// ConditionsFromMap get conditions from a map of key => values
func ConditionsFromMap(m map[string][]string, f fields.Fields) ([]Condition, error) {
	if f.IsEmpty() || len(m) == 0 {
		return nil, nil
	}

	var (
		err   error
		conds []Condition
		it    = f.Iterator()
	)

	for it.Next() {
		var (
			i     = it.Field()
			v, ok = m[i.Name]
		)

		if !ok {
			continue
		}

		var t = len(v)

		if t != 1 {
			continue // todo deal with multiple values for IN and NOT IN
		}

		var c = Condition{Attribute: i.Name}

		var w = v[0]
		if p := strings.Index(w, ":"); p != -1 {
			var o string
			o, w = w[0:p], w[p+1:]
			switch o {
			case "eq":
				c.Operator = Equals
			case "ne":
				c.Operator = NotEquals
			case "gt":
				c.Operator = GreaterThan
			case "gte":
				c.Operator = GreaterOrEqualTo
			case "lt":
				c.Operator = LessThan
			case "lte":
				c.Operator = LessOrEqualTo
			case "li":
				if i.Kind != gocipe.String {
					return nil, ErrInvalidConditionOperator
				}
				c.Operator = Like
			default:
				return nil, ErrInvalidConditionOperator
			}
		}

		switch i.Kind {
		case gocipe.Bool:
			if c.Operator != Equals && c.Operator != NotEquals {
				return nil, ErrInvalidConditionOperator
			}
			c.Value, err = gocipe.BoolFromString(w)
		case gocipe.String:
			c.Value, err = w, nil
		case gocipe.Int64:
			c.Value, err = gocipe.Int64FromString(w)
		case gocipe.Float64:
			c.Value, err = gocipe.Float64FromString(w)
		}

		if err != nil {
			return nil, err
		}

		conds = append(conds, c)
	}

	return conds, err
}

// OrderSort is the sort order of data fetched
type OrderSort uint8

// String representation (mainly for testing / debugging / logging)
func (o OrderSort) String() string {
	switch o {
	case Ascending:
		return "Ascending"
	case Descending:
		return "Descending"
	}

	return "?? " + strconv.Itoa(int(o))
}

const (
	// Ascending alphanumeric sorting first
	Ascending = OrderSort(0)

	// Descending alphanumeric sorting reversed
	Descending = OrderSort(1)
)

type OrderBy struct {
	Attribute string
	Sort      OrderSort
}

// OrderByFromString returns order by from a string (typically from uri query)
func OrderByFromString(s string, f fields.Fields) ([]OrderBy, error) {
	if s == "" {
		return nil, nil
	}

	var p = strings.Split(s, ",")
	var o = make([]OrderBy, len(p))

	for i := range p {
		var (
			attr string
			sort OrderSort
		)

		if strings.HasPrefix(p[i], "-") {
			attr = p[i][1:]
			sort = Descending
		} else {
			attr = p[i]
		}

		if attr == "" {
			continue
		}

		if !f.Contains(attr) {
			return nil, ErrUnknownSortAttribute
		}

		o[i] = OrderBy{
			Attribute: attr,
			Sort:      sort,
		}
	}

	return o, nil
}

// Pagination represents offset and limit when fetching multiple elements
type Pagination struct {
	Offset int
	Limit  int
	Order  []OrderBy
}

// Persistable is something that can be persisted by a repository
type Persistable interface {
	// Identifier returns the id of the item
	Identifier() string

	// AssignValues returns an updated version of the Persistable given a value set
	AssignValues(values values.Values) Persistable

	// Values returns values representation
	Values() values.Values
}

// Repositorium allows persistence of Entity
type Repositorium interface {
	// Get a single Entity by id
	Get(ctx context.Context, entity string, f fields.Fields, id string) (*values.Values, error)

	// List multiple Entity with pagination rules and conditions
	List(ctx context.Context, entity string, f fields.Fields, p Pagination, c ...Condition) ([]values.Values, error)

	// Delete a single Entity by id
	Delete(ctx context.Context, entity, id string) error

	// DeleteWhere delete multiple Entity based on conditions
	DeleteWhere(ctx context.Context, entity string, c ...Condition) error

	// Create a new Entity in persistent storage
	Create(ctx context.Context, entity string, vals *values.Values) (string, error)

	// Update an existing Entity in persistent storage
	Update(ctx context.Context, entity string, id string, vals *values.Values) error

	// UpdateValuesWhere Values in persistent storage
	UpdateWhere(ctx context.Context, entity string, vals *values.Values, c ...Condition) error

	// Close connection to the repo
	Close() error
}
