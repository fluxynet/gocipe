package gocipe

import "strings"

type Entity struct {
	Name        string
	Description string
	Fields      []Field
}

func (e Entity) String() string {
	var (
		s []string
		description string
	)

	if e.Description != "" {
		description = " // " + description
	}

	s = append(s, "type " + e.Name + " { " + description)

	for _, f := range e.Fields {
		s = append(s, "\t"+f.String())
	}

	s = append(s, "}")

	return strings.Join(s, "\n") + "\n"
}

type Field struct {
	Name        string
	Description string
	Type        string
	IsMany      bool
	IsComplex   bool
	IsEmbedded  bool
	IsPointer   bool
}

func (f Field) String() string {
	var s string

	if !f.IsEmbedded {
		s = f.Name + " "
	}

	if f.IsMany {
		s += "[]"
	}

	if f.IsPointer {
		s += "*"
	}

	s += f.Type

	if f.IsComplex {
		s += "*"
	}

	if f.Description != "" {
		s += " // " + f.Description
	}

	return s
}
