package openapi

import (
	"strings"

	"github.com/fluxynet/gocipe/api"
	"github.com/fluxynet/gocipe/types/fields"
)

// Resource is an openapi resource
type Resource struct {
	name        string
	description string
	actions     api.ActionSet
	props       Properties
	path        string
}

func (r *Resource) SetName(name string) *Resource {
	r.name = strings.ToLower(name)
	return r
}

func (r Resource) Name() string {
	return r.name
}

func (r *Resource) SetDescription(description string) *Resource {
	r.description = description
	return r
}

func (r Resource) Description() string {
	return r.description
}

func (r *Resource) SetProperties(props ...Property) *Resource {
	r.props = props
	return r
}

func (r *Resource) SetPath(path string) *Resource {
	r.path = path
	return r
}

func (r Resource) Path() string {
	if r.path == "" && r.name != "" {
		return "/" + r.name
	}
	return "/" + r.path
}

func (r *Resource) SetActions(actions api.ActionSet) *Resource {
	r.actions = actions
	return r
}

func (r Resource) Actions() api.ActionSet {
	return r.actions
}

func (r Resource) Properties() Properties {
	return r.props
}

func (r Resource) Fields() fields.Fields {
	return r.props.Fields()
}
