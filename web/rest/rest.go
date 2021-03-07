package rest

import (
	"errors"
	"net/http"

	"github.com/fluxynet/gocipe/fields"
)

// MethodSet enabled on a resource
type MethodSet uint8

const (
	// ENDPOINTS_GET Endpoint to GET is enabled
	ENDPOINTS_GET = MethodSet(1)

	// ENDPOINTS_LIST Endpoint to LIST is enabled
	ENDPOINTS_LIST = MethodSet(2)

	// ENDPOINTS_DELETE Endpoint to DELETE is enabled
	ENDPOINTS_DELETE = MethodSet(4)

	// ENDPOINTS_CREATE Endpoint to CREATE is enabled
	ENDPOINTS_CREATE = MethodSet(6)

	// ENDPOINTS_REPLACE Endpoint to REPLACE is enabled
	ENDPOINTS_REPLACE = MethodSet(8)

	// ENDPOINTS_UPDATE Endpoint to UPDATE is enabled
	ENDPOINTS_UPDATE = MethodSet(10)

	// ENDPOINTS_ALL all endpoints enabled
	ENDPOINTS_ALL = MethodSet(0)
)

// Is checks if a method set is contained with a value
func (m MethodSet) Has(v MethodSet) bool {
	if m == ENDPOINTS_ALL {
		return true
	}

	return (m & v) != 0
}

// IsNot checks if a method set is NOT contained with a value
func (m MethodSet) NotHas(v MethodSet) bool {
	if m == ENDPOINTS_ALL {
		return false
	}

	return (m & v) == 0
}

var (
	// ErrIdNotPresent indicates id is not present in a request
	ErrIdNotPresent = errors.New("id is not present")
)

// GetIdFunc is a function that returns an id from an http.Request
type GetIdFunc func(r *http.Request) (string, error)

// ResourceDef specifies resources for a rest resource endpoint
type ResourceDef struct {
	// Path defaults to entity
	Path string

	// Entity the entity name
	Entity string

	// Fields the fields available
	Fields fields.Fields

	// Methods enabled
	Methods MethodSet
}
