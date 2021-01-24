package fields

import (
	"github.com/fluxynet/gocipe"
	"strings"
)

// Field as part of a set
type Field struct {
	Name     string
	Kind     gocipe.Type
	Required bool
	prev *Field
	next     *Field
}

// Fields representing a field set
type Fields struct {
	head  *Field
	tail  *Field
	items map[string]*Field
}

// IsEmpty checks if there are any fields in the set
func (f Fields) IsEmpty() bool {
	return f.head == nil
}

// Length of the field set
func (f Fields) Length() int {
	return len(f.items)
}

// Unset removes an item
func (f *Fields) Unset(name string) *Fields {
	var d, ok = f.items[name]
	if !ok {
		return nil
	}

	if d.prev == nil { // only one in list
		f.head = nil
		f.tail = nil
		delete(f.items, name)
		return nil
	}

	d.prev.next = d.next
	d.prev = nil
	d.next = nil

	return f
}

// Set a named Field kind
func (f *Fields) Set(name string, kind gocipe.Type) *Fields {
	var required bool

	if strings.HasPrefix(name, "+") {
		required = true
		name = name[1:]
	}

	var node = &Field{Name: name, Kind: kind, Required: required}

	if f.head == nil { // list is empty
		f.head = node
		f.tail = f.head
		f.items = map[string]*Field{name: f.head}
		return nil
	}

	if _, ok := f.items[name]; ok { // list contains item, must replace
		f.items[name].Kind = kind
		return nil
	}

	node.prev = f.tail
	f.tail.next = node
	f.tail = node
	f.items[name] = node

	return f
}

// GetEmptyValues for fields consisting of pointers whereunto the data can be placed
func (f Fields) GetEmptyValues() map[string]interface{} {
	var (
		it  = f.Iterator()
		dst = make(map[string]interface{}, f.Length())
	)

	for it.Next() {
		var i = it.Field()
		dst[i.Name] = gocipe.DefaultValue(i.Kind)
	}

	return dst
}

// TypeOf returns typeof a field
func (f Fields) TypeOf(name string) gocipe.Type {
	if f.head == nil {
		return gocipe.Undefined
	}

	return f.items[name].Kind
}

// String representation (for debugging mainly)
func (f Fields) String() string {
	var (
		s  []string
		it = f.Iterator()
	)

	for it.Next() {
		var req string
		f := it.Field()

		if f.Required {
			req = "!"
		}

		s = append(s, f.Name+":"+string(f.Kind)+req)
	}

	return strings.Join(s, ", ")
}

// Iterator returns an iterator to loop though fields based on order
func (f Fields) Iterator() Iterator {
	return &iterator{head: f.head}
}

// Iterator allows moving through a list
type Iterator interface {
	// Reset the iteration
	Reset()

	// Next moves the iterator if next is available and returns true; if end has been reached returns false
	Next() bool

	// Field returns the current field
	Field() *Field
}

type iterator struct {
	head    *Field
	current *Field
}

func (i *iterator) Reset() {
	i.current = nil
}

func (i *iterator) Next() bool {
	if i.head == nil {
		return false
	}

	if i.current == nil {
		i.current = i.head
	} else if i.current.next == nil {
		return false
	} else {
		i.current = i.current.next
	}

	return true
}

func (i *iterator) Field() *Field {
	return i.current
}

// FromMap creates a field set from map
func FromMap(m map[string]gocipe.Type) Fields {
	var f Fields

	for k, v := range m {
		f.Set(k, v)
	}

	return f
}

// FromPairs sets Values from a slice of Value
func FromPairs(p []Field) Fields {
	var f Fields
	for i := range p {
		f.Set(p[i].Name, p[i].Kind)
	}

	return f
}
