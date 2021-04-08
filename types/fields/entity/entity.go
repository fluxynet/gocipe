package entity

import "github.com/fluxynet/gocipe/types/fields"

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
