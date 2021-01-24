package repository

import (
	"context"
	"errors"
	"github.com/fluxynet/gocipe/fields"
	"github.com/fluxynet/gocipe/values"
)

var (
	// ErrNotFound is when the item was not found
	ErrNotFound = errors.New("item not found")
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

	// In denotes In
	In = ConditionOperator(6)

	// NotIn denotes Not in
	NotIn = ConditionOperator(7)
)

// ConditionType represents what kind of condition it is with respect to other conditions (AND / OR)
type ConditionType uint8

const (
	// And means the condition must be inclusive wrt other conditions
	And = ConditionType(0)

	// Or means the condition needs not be inclusive wrt other conditions
	Or = ConditionType(1)
)

// OrderSort is the sort order of data fetched
type OrderSort uint8

const (
	// Ascending alphanumeric sorting first
	Ascending = OrderSort(0)

	// Descending alphanumeric sorting reversed
	Descending = OrderSort(1)
)

// Condition represents filter criteria for fetching multiple elements
type Condition struct {
	Property string
	Operator ConditionOperator
	Value    interface{}
	Type     ConditionType
}

type OrderBy struct {
	Attribute string
	Sort      OrderSort
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
	Get(ctx context.Context, entity string, f fields.Fields, id string) (values.Values, error)

	// List multiple Entity with pagination rules and conditions
	List(ctx context.Context, entity string, f fields.Fields, p Pagination, c ...Condition) ([]values.Values, error)

	// Delete a single Entity by id
	Delete(ctx context.Context, entity, id string) error

	// DeleteWhere delete multiple Entity based on conditions
	DeleteWhere(ctx context.Context, entity string, c ...Condition) error

	// Create a new Entity in persistent storage
	Create(ctx context.Context, entity string, vals values.Values) (string, error)

	// Update an existing Entity in persistent storage
	Update(ctx context.Context, entity string, id string, vals values.Values) error

	// UpdateValuesWhere Values in persistent storage
	UpdateValues(ctx context.Context, entity string, vals values.Values, c ...Condition) error
}
