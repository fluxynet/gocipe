package gocipe

import (
	"errors"
	"strconv"
)

var (
	// ErrInvalidValue is when a value sent is invalid
	ErrInvalidValue = errors.New("invalid value provided")
)

// Type represents a variable type
type Type string

const (
	// Undefined type duh
	Undefined = Type("")

	// Bool indicates native bool
	Bool = Type("bool")

	// String indicates native string
	String = Type("string")

	// Int64 indicates native int64
	Int64 = Type("int64")

	// Float64 indicates native float64
	Float64 = Type("float64")
)

// DefaultValue for types
func DefaultValue(t Type) interface{} {
	switch t {
	case Bool:
		return true
	case String:
		return ""
	case Int64:
		return 0
	case Float64:
		return float64(0)
	}

	return nil
}

// DefaultPointer for types
func DefaultPointer(t Type) interface{} {
	switch t {
	case Bool:
		return new(bool)
	case String:
		return new(string)
	case Int64:
		return new(int64)
	case Float64:
		return new(float64)
	}

	return nil
}

// BoolFromString parses Bool
func BoolFromString(s string) (bool, error) {
	switch s {
	case "true":
		return true, nil
	case "false":
		return false, nil
	}

	return false, ErrInvalidValue
}

// Int64FromString parses Int64
func Int64FromString(s string) (int64, error) {
	return strconv.ParseInt(string(s), 10, 64)
}

// Float64FromString parses Float64
func Float64FromString(s string) (float64, error) {
	return strconv.ParseFloat(string(s), 64)
}

// Parser introspects a file and returns defined
type Parser interface {
	Parse(name string, src interface{}) error
	Entities() Entities
}

// Entities is a reference of entity_name => Entity
type Entities map[string]Entity

//
//// ResolveEntityEmbeds resolves all embedded types in entities
//func ResolveEntityEmbeds(entities Entities, ent Entity) Entity {
//	var (
//		entity   Entity
//		resolved bool
//	)
//
//	entity.Name = ent.Name
//
//	for _, f := range ent.Fields {
//		if f.IsEmbedded {
//			embedded, ok := entities[f.Type]
//			if !ok {
//				continue
//			}
//
//			entity.Fields = append(entity.Fields, embedded.Fields...)
//			resolved = true
//		} else {
//			entity.Fields = append(entity.Fields, f)
//		}
//	}
//
//	if resolved {
//		return ResolveEntityEmbeds(entities, entity)
//	}
//
//	return entity
//}
//
//func EntitiesToMap(entities []Entity) Entities {
//	var entMap = make(Entities)
//	for _, entity := range entities {
//		entMap[entity.Name] = entity
//	}
//
//	return entMap
//}
