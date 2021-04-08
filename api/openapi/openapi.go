package openapi

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/fluxynet/gocipe/api"
)

const (
	paramID     = "#/parameters/id"
	paramLimit  = "#/parameters/limit"
	paramOffset = "#/parameters/offset"
	paramSort   = "#/parameters/sort"
)

var (
	statusOK           = strconv.Itoa(http.StatusOK)
	statusCreated      = strconv.Itoa(http.StatusCreated)
	statusBadRequest   = strconv.Itoa(http.StatusBadRequest)
	statusUnauthorized = strconv.Itoa(http.StatusUnauthorized)
	statusForbidden    = strconv.Itoa(http.StatusForbidden)
	statusNotFound     = strconv.Itoa(http.StatusNotFound)
)

var (
	responseBadRequest = Response{
		Description: "Bad Request. Request body is incorrect or incomplete",
	}

	responseUnauthorized = Response{
		Description: "Unauthorized. Authentication required",
	}

	responseForbidden = Response{
		Description: "Forbidden. This action is not allowed for current user",
	}

	responseNotFound = Response{
		Description: "Not found",
	}
)

type DocumentInformation struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Version     string `json:"version,omitempty"`
}

// Definitions name -> definition map
type Definitions map[string]Definition

// PropMeta is a key => value list of metadata for a property
type PropMeta map[string]interface{}

// Definition of a resource
type Definition struct {
	Description string              `json:"description,omitempty"`
	Properties  map[string]PropMeta `json:"properties,omitempty"`
}

type Paths map[string]Path

type Path map[string]PathDetails

type PathDetails struct {
	Description string              `json:"description,omitempty"`
	Parameters  []Parameter         `json:"parameters,omitempty"`
	Responses   map[string]Response `json:"responses,omitempty"`
}

type Parameter struct {
	Ref         string              `json:"$ref,omitempty"`
	Name        string              `json:"name,omitempty"`
	Description string              `json:"description,omitempty"`
	Type        string              `json:"type,omitempty"`
	In          string              `json:"in,omitempty"`
	Required    bool                `json:"required,omitempty"`
	Schema      map[string]PropMeta `json:"schema,omitempty"`
}

type Response struct {
	Description string `json:"description,omitempty"`
	Schema      Schema `json:"schema,omitempty"`
}

type Schema struct {
	Type  string            `json:"type,omitempty"`
	Items map[string]string `json:"items,omitempty"`
	Ref   string            `json:"$ref,omitempty"`
}

type Document struct {
	Swagger     string              `json:"swagger,omitempty"`
	Info        DocumentInformation `json:"info,omitempty"`
	Host        string              `json:"host,omitempty"`
	BasePath    string              `json:"base_path,omitempty"`
	Consumes    []string            `json:"consumes,omitempty"`
	Produces    []string            `json:"produces,omitempty"`
	Definitions Definitions         `json:"definitions,omitempty"`
	Paths       Paths               `json:"paths,omitempty"`
}

// HandlerFunc returns an http.HandlerFunc for the openApi json representation
// it is not live, resources added afterwards are not added
func (d Document) HandlerFunc() http.HandlerFunc {
	var b, err = json.Marshal(d)
	if err != nil {
		panic(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(b)
	}
}

func (d *Document) AddResources(res ...Resource) *Document {
	for i := range res {
		d.AddResource(res[i])
	}

	return d
}

func (d *Document) AddResource(res Resource) *Document {
	var (
		name    = res.Name()
		ref     = "#/definitions/" + name
		nameid  = name + "/{id}"
		actions = res.Actions()
	)

	if d.Definitions == nil {
		d.Definitions = make(map[string]Definition)
	}

	d.Definitions[name] = Definition{
		Description: res.Description(),
		Properties:  propsAsMap(res.props, false),
	}

	if d.Paths == nil {
		d.Paths = make(map[string]Path)
	}

	if actions.Has(api.ActionList | api.ActionCreate) {
		d.Paths[name] = make(map[string]PathDetails)
	}

	if actions.Has(api.ActionRead | api.ActionReplace | api.ActionUpdate | api.ActionDelete) {
		d.Paths[nameid] = make(map[string]PathDetails)
	}

	if actions.Has(api.ActionList) {
		d.Paths[name][http.MethodGet] = PathDetails{
			Description: "Get a list of " + name + " items",
			Parameters:  paramsList(res.props),
			Responses: map[string]Response{
				statusOK: {
					Description: "OK",
					Schema: Schema{
						Type: "array",
						Items: map[string]string{
							"$ref": "#/definitions/" + name,
						},
					},
				},
				statusUnauthorized: responseUnauthorized,
				statusForbidden:    responseForbidden,
				statusNotFound:     responseNotFound,
			},
		}
	}

	if actions.Has(api.ActionRead) {
		d.Paths[nameid][http.MethodGet] = PathDetails{
			Description: "Get a single " + name + " by id",
			Parameters: []Parameter{{
				Required: true,
				In:       "query",
				Ref:      paramID,
			}},
			Responses: map[string]Response{
				statusOK: {
					Description: "",
					Schema: Schema{
						Ref: ref,
					},
				},
				statusUnauthorized: responseUnauthorized,
				statusForbidden:    responseForbidden,
				statusNotFound:     responseNotFound,
			},
		}
	}

	if actions.Has(api.ActionCreate) {
		d.Paths[name][http.MethodPost] = PathDetails{
			Description: "Create a new " + name,
			Parameters: []Parameter{
				{
					Required: true,
					In:       "body",
					Ref:      ref,
				},
			},
			Responses: map[string]Response{
				statusCreated: {
					Description: "Created",
					Schema: Schema{
						Ref: ref,
					},
				},
				statusBadRequest:   responseBadRequest,
				statusUnauthorized: responseUnauthorized,
				statusForbidden:    responseForbidden,
				statusNotFound:     responseNotFound,
			},
		}
	}

	if actions.Has(api.ActionReplace) {
		d.Paths[nameid][http.MethodPut] = PathDetails{
			Description: "Replace an existing " + name,
			Parameters: []Parameter{
				{
					Required: true,
					In:       "body",
					Ref:      ref,
				},
			},
			Responses: map[string]Response{
				statusOK: {
					Description: "Updated",
				},
				statusBadRequest:   responseBadRequest,
				statusUnauthorized: responseUnauthorized,
				statusForbidden:    responseForbidden,
				statusNotFound:     responseNotFound,
			},
		}
	}

	if actions.Has(api.ActionUpdate) {
		d.Paths[nameid][http.MethodPatch] = PathDetails{
			Description: "Update an existing " + name,
			Parameters: []Parameter{
				{
					Name:        "Partial " + name,
					Description: "Fields and values that need to be updated only",
					Required:    true,
					In:          "body",
					Schema:      propsAsMap(res.props, true),
				},
			},
			Responses: map[string]Response{
				statusOK: {
					Description: "Updated",
				},
				statusBadRequest:   responseBadRequest,
				statusUnauthorized: responseUnauthorized,
				statusForbidden:    responseForbidden,
				statusNotFound:     responseNotFound,
			},
		}
	}

	if actions.Has(api.ActionDelete) {
		d.Paths[nameid][http.MethodDelete] = PathDetails{
			Description: "Delete an existing " + name + " by id",
			Parameters: []Parameter{{
				Required: true,
				In:       "query",
				Ref:      paramID,
			}},
			Responses: map[string]Response{
				statusOK: {
					Description: "Deleted",
				},
				statusUnauthorized: responseUnauthorized,
				statusForbidden:    responseForbidden,
				statusNotFound:     responseNotFound,
			},
		}
	}

	return d
}

func paramsList(props Properties) []Parameter {
	var p = []Parameter{
		{Ref: paramLimit},
		{Ref: paramOffset},
		{Ref: paramSort},
	}

	for i := range props {
		p = append(p, Parameter{
			Name:        props[i].Name,
			Description: props[i].Description,
			Type:        string(props[i].Kind),
			In:          "query",
		})
	}

	return p
}

func propsAsMap(props Properties, allowPartial bool) map[string]PropMeta {
	var m = make(map[string]PropMeta, len(props))

	for i := range props {
		m[props[i].Name] = propMeta(props[i], allowPartial)
	}

	return m
}

func propMeta(prop Property, allowPartial bool) PropMeta {
	var p = make(PropMeta)

	if prop.Description != "" {
		p["description"] = prop.Description
	}

	if prop.Example != nil {
		p["example"] = prop.Example
	}

	if len(prop.Enum) != 0 {
		p["enum"] = prop.Enum
	}

	if prop.Maximum != 0 {
		p["maximum"] = prop.Maximum
	}

	if prop.Minimum != 0 {
		p["minimum"] = prop.Minimum
	}

	if prop.MaxLength != 0 {
		p["maxlength"] = prop.MaxLength
	}

	if prop.MinLength != 0 {
		p["minlength"] = prop.MinLength
	}

	if allowPartial {
		p["required"] = false
	} else {
		p["required"] = prop.Required
	}

	if len(prop.Items) != 0 {
		p["items"] = propsAsMap(prop.Items, false)
	}

	if prop.Ref != "" {
		p["$ref"] = prop.Ref
	}

	return p
}
