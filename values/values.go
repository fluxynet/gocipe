package values

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fluxynet/gocipe"
	"github.com/fluxynet/gocipe/fields"
	"io"
	"strconv"
	"strings"
)

var (
	// ErrValueMissing is when a value is expected but not found from http.Request
	ErrValueMissing = errors.New("mandatory value missing")

	// ErrInvalidValue is when a value sent is invalid
	ErrInvalidValue = errors.New("invalid value provided")
)

var (
	boolTrue  = []byte("true")
	boolFalse = []byte("false")
)

// Value is a named data value
type Value struct {
	Name  string
	Value interface{}
	prev  *Value
	next  *Value
}

// IsBool returns true if value held is of type bool
func (v Value) IsBool() bool {
	var _, ok = v.Value.(bool)
	return ok
}

// Bool returns the value as a bool if it is of type bool or an empty bool
func (v Value) Bool() bool {
	var x = v.Value.(bool)
	return x
}

// IsString returns true if value held is of type string
func (v Value) IsString() bool {
	var _, ok = v.Value.(string)
	return ok
}

// String returns the value as a string if it is of type string or an empty string
func (v Value) String() string {
	var x = v.Value.(string)
	return x
}

// IsInt returns true if value held is of type int
func (v Value) IsInt64() bool {
	var _, ok = v.Value.(int)
	return ok
}

// Int64 returns the value as a int if it is of type int or an empty int
func (v Value) Int64() int64 {
	var x = v.Value.(int64)
	return x
}

// IsFloat64 returns true if value held is of type float64
func (v Value) IsFloat64() bool {
	var _, ok = v.Value.(float64)
	return ok
}

// Float64 returns the value as a float64 if it is of type float64 or an empty float64
func (v Value) Float64() float64 {
	var x = v.Value.(float64)
	return x
}

// Values is an ordered list of Value items that can be randomly accessed. Not thread-safe.
type Values struct {
	head  *Value
	tail  *Value
	items map[string]*Value
}

// IsEmpty checks if there are any values in the set
func (v Values) IsEmpty() bool {
	return v.head == nil
}

// Length of the value set
func (v Values) Length() int {
	return len(v.items)
}

// Set a named Value key
func (v *Values) Set(name string, value interface{}) {
	var node = &Value{Name: name, Value: value}

	if v.head == nil { // list is empty
		v.head = node
		v.tail = v.head
		v.items = map[string]*Value{name: v.head}
		return
	}

	if n := v.items[name]; n != nil { // list contains item, must replace
		n.Value = value
		return
	}

	node.prev = v.tail
	v.tail.next = node
	v.tail = node
	v.items[name] = node
}

// Get returns raw Value
func (v Values) Get(name string) *Value {
	if v.head == nil {
		return nil
	}

	return v.items[name]
}

// Unset removes a value item
func (v Values) Unset(name string) Values {
	var d, ok = v.items[name]
	if !ok {
		 return v
	}

	if d.prev == nil { // only one in list
		v.head = nil
		v.tail = nil
		delete(v.items, name)
		return v
	}

	d.prev.next = d.next
	d.prev = nil
	d.next = nil

	return v
}

// ToMap Returns a map[string]interface{} representation of the list
func (v Values) ToMap() map[string]interface{} {
	var m = make(map[string]interface{}, len(v.items))
	var it = v.Iterator()

	for it.Next() {
		var c = it.Value()
		m[c.Name] = c.Value
	}

	return m
}

// FromMap sets values from a map into the Values set
func (v Values) FromMap(m map[string]interface{}) Values {
	for key, val := range m {
		v.Set(key, val)
	}

	return v
}

// FromPairs sets Values from a slice of Value
func (v Values) FromPairs(p []Value) Values {
	for i := range p {
		v.Set(p[i].Name, p[i].Value)
	}

	return v
}

// FromMap initiates a Values structure from a map
func FromMap(m map[string]interface{}) Values {
	var values Values

	for key, val := range m {
		values.Set(key, val)
	}

	return values
}

// FromPairs initiates a Values structure from a slice of Value
func FromPairs(p []Value) Values {
	var vals Values
	for i := range p {
		vals.Set(p[i].Name, p[i].Value)
	}

	return vals
}

// FromJSON returns values from an http.Request
func FromJSON(r io.ReadCloser, f fields.Fields) (Values, error) {
	var b, err = io.ReadAll(r)
	defer gocipe.Closed(r, &err)

	if err != nil {
		return Values{}, err
	}

	var m = make(map[string]json.RawMessage, f.Length())

	err = json.Unmarshal(b, &m)
	if err != nil {
		return Values{}, err
	}

	var vals Values

	var it = f.Iterator()
	for it.Next() {
		var i = it.Field()
		var v, ok = m[i.Name]
		var x interface{}

		if i.Required && !ok {
			return vals, ErrValueMissing
		}

		switch i.Kind {
		case gocipe.Bool:
			x = bytes.Equal(v, boolTrue)
		case gocipe.String:
			var s = string(v)
			if strings.HasPrefix(s, `"`) && strings.HasSuffix(s, `"`) {
				x = strings.TrimSuffix(strings.TrimPrefix(s, `"`), `"`)
			} else {
				err = ErrInvalidValue
			}
		case gocipe.Int64:
			x, err = strconv.Atoi(string(v))
		case gocipe.Float64:
			x, err = strconv.ParseFloat(string(v), 64)
		}

		if err != nil {
			return vals, err
		}

		vals.Set(i.Name, x)
	}

	return vals, nil
}

func (v Values) String() string {
	var (
		s  []string
		it = v.Iterator()
	)

	for it.Next() {
		f := it.Value()
		s = append(s, fmt.Sprintf("%s:%#v", f.Name, f.Value))
	}

	return strings.Join(s, ", ")
}

// Iterator returns an iterator to loop though values based on order
func (v Values) Iterator() Iterator {
	return &iterator{head: v.head}
}

// Iterator allows moving through a list
type Iterator interface {
	// Next moves the iterator if next is available and returns true; if end has been reached returns false
	Next() bool

	// Value returns the current value
	Value() *Value
}

type iterator struct {
	head    *Value
	current *Value
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

func (i *iterator) Value() *Value {
	return i.current
}
