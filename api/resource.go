package api

import (
	"github.com/fluxynet/gocipe/types/fields/entity"
)

// ActionSet enabled on a resource
type ActionSet uint8

const (
	// ActionRead Endpoint to GET is enabled
	ActionRead = ActionSet(1)

	// ActionList Endpoint to LIST is enabled
	ActionList = ActionSet(2)

	// ActionDelete Endpoint to DELETE is enabled
	ActionDelete = ActionSet(4)

	// ActionCreate Endpoint to CREATE is enabled
	ActionCreate = ActionSet(8)

	// ActionReplace Endpoint to REPLACE is enabled
	ActionReplace = ActionSet(16)

	// ActionUpdate Endpoint to UPDATE is enabled
	ActionUpdate = ActionSet(32)

	// ActionAll all endpoints enabled
	ActionAll = ActionSet(33)
)

// Has checks if an action set is contained with a value
func (s ActionSet) Has(v ActionSet) bool {
	if s == ActionAll {
		return true
	}

	return (s & v) != 0
}

// NotHas checks if an action set is NOT contained with a value
func (s ActionSet) NotHas(v ActionSet) bool {
	if s == ActionAll {
		return false
	}

	return (s & v) == 0
}

type ResourceOpts struct {
	entity  entity.Entity
	Path    string
	Actions ActionSet
}

func New(opts ResourceOpts) Resource {
	return resource{
		name:    opts.entity.Name(),
		path:    opts.Path,
		actions: opts.Actions,
	}

}

// Resource is an entity served via REST
type Resource interface {
	entity.Entity

	// Path to the resource
	Path() string

	// Actions allowed for this resource
	Actions() ActionSet
}

// resource is an entity served via REST
type resource struct {
	entity.Entity

	// name of the resource
	name string

	// Path to the resource, omit prefix and trailing slash
	path string

	// Actions enabled
	actions ActionSet
}

func (r resource) Path() string {
	if r.path == "" {
		return "/" + r.name
	}

	return "/" + r.path
}

func (r resource) Actions() ActionSet {
	return r.actions
}
