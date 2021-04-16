package entity

import (
	"github.com/fluxynet/gocipe/types"
	"github.com/fluxynet/gocipe/types/fields"
)

// Entity is that which returns fields
type Entity interface {
	Name() string
	Fields() fields.Fields
}

// Entities is a collection of entities
type Entities []Entity

// Add an entity to the collection
func (e *Entities) Add(entities ...Entity) {
	*e = append(*e, entities...)
}

// From to initialize a slice of entities
// makes it easier to initialize non-homogenous elements,
// all of whom adhering to Entity interface
func From(entities ...Entity) Entities {
	return entities
}

// partial is a simple entity representation containing only specific fields
type partial struct {
	name   string
	fields fields.Fields
}

func (p partial) Name() string {
	return p.name
}

func (p partial) Fields() fields.Fields {
	return p.fields
}

// Partial representation of an entity
func Partial(name string, f fields.Fields) Entity {
	return partial{
		name:   name,
		fields: f,
	}
}

// ID is a partial entity representation containing only the ID
func ID(name string) Entity {
	return partial{
		name: name,
		fields: fields.From(
			fields.Field{
				Name: "id",
				Kind: types.String,
			},
		),
	}
}
