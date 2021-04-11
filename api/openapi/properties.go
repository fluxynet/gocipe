package openapi

import (
	"github.com/fluxynet/gocipe/types/fields"
)

// Property is a field with openapi metadata
type Property struct {
	fields.Field
	prev        *Property
	next        *Property
	Description string
	Example     interface{}
	Enum        []interface{}
	Maximum     int
	Minimum     int
	MaxLength   int
	MinLength   int
	Required    bool
	Items       Properties
	Ref         string
}

// Properties is a collection of properties
type Properties []Property

// Fields version of the properties
func (p Properties) Fields() fields.Fields {
	var f fields.Fields

	for i := range p {
		f.Set(p[i].Name, p[i].Kind)
	}

	return f
}
