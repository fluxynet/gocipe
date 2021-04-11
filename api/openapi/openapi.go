package openapi

import (
	"net/http"

	"github.com/fluxynet/gocipe/api"
	"github.com/fluxynet/gocipe/types"
	"github.com/fluxynet/gocipe/util"
	"github.com/getkin/kin-openapi/openapi3"
)

// HandlerFunc returns an http.HandlerFunc for the openApi json representation
// it is not live, resources added afterwards are not added
func (s *Swagger) HandlerFunc() http.HandlerFunc {
	var b, err = s.MarshalJSON()
	if err != nil {
		panic(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(b)
	}
}

func (s *Swagger) AddResources(res ...Resource) *Swagger {
	for i := range res {
		s.AddResource(res[i])
	}

	return s
}

func (s *Swagger) AddResource(res Resource) *Swagger {
	var (
		name     = res.Name()
		ref      = "#/components/schemas/" + name
		path     = res.Path()
		pathID   = path + "/{id}"
		actions  = res.Actions()
		paramsID = openapi3.Parameters{
			&openapi3.ParameterRef{Ref: paramID},
		}
	)

	// no actions, skip
	if actions&api.ActionAll == 0 {
		return s
	}

	if s.Components.Schemas == nil {
		s.Components.Schemas = make(openapi3.Schemas)
	}

	s.Components.Schemas[name] = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type:                 "object",
			Description:          res.Description(),
			Enum:                 nil,
			Required:             requiredProps(res.Properties()),
			Properties:           propsToSchemas(res.Properties()),
			MinProps:             0,
			MaxProps:             nil,
			AdditionalProperties: nil,
			Discriminator:        nil,
		},
	}

	if s.Paths == nil {
		s.Paths = make(map[string]*openapi3.PathItem)
	}

	if actions.Has(api.ActionList | api.ActionCreate) {
		s.Paths[path] = new(openapi3.PathItem)
	}

	if actions.Has(api.ActionRead | api.ActionReplace | api.ActionUpdate | api.ActionDelete) {
		s.Paths[pathID] = new(openapi3.PathItem)
	}

	if actions.Has(api.ActionList) {
		s.Paths[path].Get = &openapi3.Operation{
			Description: "Get a list of " + name + " items",
			Parameters:  paramsList(res.props),
			Responses: responsesWithErrors(openapi3.Responses{
				statusOK: &openapi3.ResponseRef{
					Value: &openapi3.Response{
						Description: util.Str("OK - List of items"),
						Content: map[string]*openapi3.MediaType{
							contentTypeJSON: {
								Schema: &openapi3.SchemaRef{
									Ref: ref,
								},
							},
						},
					},
				},
			}, actions),
		}
	}

	if actions.Has(api.ActionRead) {
		s.Paths[pathID].Get = &openapi3.Operation{
			Description: "Get a single " + name + " by id",
			Parameters:  paramsID,
			Responses: responsesWithErrors(openapi3.Responses{
				statusOK: &openapi3.ResponseRef{
					Ref: ref,
				},
			}, actions),
		}
	}

	if actions.Has(api.ActionCreate) {
		s.Paths[path].Post = &openapi3.Operation{
			Description: "Create a new " + name,
			RequestBody: &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Content: map[string]*openapi3.MediaType{
						contentTypeJSON: {
							Schema: &openapi3.SchemaRef{Ref: ref},
						},
					},
				},
			},
			Responses: responsesWithErrors(openapi3.Responses{
				statusOK: &openapi3.ResponseRef{Ref: ref},
			}, actions),
		}
	}

	if actions.Has(api.ActionReplace) {
		s.Paths[pathID].Put = &openapi3.Operation{
			Description: "Replace an existing " + name,
			Parameters:  paramsID,
			RequestBody: &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Content: map[string]*openapi3.MediaType{
						contentTypeJSON: {
							Schema: &openapi3.SchemaRef{Ref: ref},
						},
					},
				},
			},
			Responses: responsesWithErrors(openapi3.Responses{
				statusOK: &openapi3.ResponseRef{Ref: ref},
			}, actions),
		}
	}

	if actions.Has(api.ActionUpdate) {
		s.Paths[pathID].Patch = &openapi3.Operation{
			Description: "Update an existing " + name,
			Parameters:  paramsID,
			RequestBody: &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Description: "Fields and values that need to be updated only",
					Content: map[string]*openapi3.MediaType{
						contentTypeJSON: {
							Schema: &openapi3.SchemaRef{
								Value: &openapi3.Schema{
									Description: "Partial or complete definition of " + name,
									Properties:  propsToSchemas(res.props),
								},
							},
						},
					},
				},
			},
			Responses: responsesWithErrors(openapi3.Responses{
				statusOK: &openapi3.ResponseRef{
					Value: &openapi3.Response{
						Description: util.Str("Updated successfully"),
					},
				},
			}, actions),
		}
	}

	if actions.Has(api.ActionDelete) {
		s.Paths[pathID].Delete = &openapi3.Operation{
			Description: "Delete an existing " + name + " by id",
			Parameters:  paramsID,
			Responses: responsesWithErrors(openapi3.Responses{
				statusOK: &openapi3.ResponseRef{
					Value: &openapi3.Response{
						Description: util.Str("Deleted successfully"),
					},
				},
			}, actions),
		}
	}

	return s
}

func paramsList(props Properties) []*openapi3.ParameterRef {
	var p = []*openapi3.ParameterRef{
		{Ref: paramLimit},
		{Ref: paramOffset},
		{Ref: paramSort},
	}

	for i := range props {
		p = append(p, &openapi3.ParameterRef{
			Value: &openapi3.Parameter{
				ExtensionProps: openapi3.ExtensionProps{},
				Name:           props[i].Name,
				Description:    props[i].Description,
				In:             "query",
				Schema: &openapi3.SchemaRef{
					Value: propToSchema(props[i]),
				},
			},
		})
	}

	return p
}

func propsToSchemas(props Properties) openapi3.Schemas {
	var m = make(openapi3.Schemas, len(props))

	for i := range props {

		if props[i].Ref == "" {
			m[props[i].Name] = &openapi3.SchemaRef{
				Value: propToSchema(props[i]),
			}
		} else {
			m[props[i].Name] = &openapi3.SchemaRef{
				Ref: props[i].Ref,
			}
		}

	}

	return m
}

func requiredProps(props Properties) []string {
	var r []string

	for i := range props {
		if props[i].Required {
			r = append(r, props[i].Name)
		}
	}

	return r
}

func propToSchema(prop Property) *openapi3.Schema {
	var p openapi3.Schema

	switch prop.Kind {
	case types.Bool:
		p.Type = "boolean"
	case types.String:
		p.Type = "string"
	case types.Int64:
		p.Type = "integer"
		p.Format = "int64"
	case types.Float64:
		p.Type = "number"
		p.Format = "double"
	}

	p.Description = prop.Description
	p.Example = prop.Example
	p.Enum = prop.Enum

	if prop.Maximum != 0 {
		var v = float64(prop.Maximum)
		p.Max = &v
	}

	if prop.Minimum != 0 {
		var v = float64(prop.Minimum)
		p.Min = &v
	}

	if prop.MaxLength != 0 {
		var v = uint64(prop.MaxLength)
		p.MaxLength = &v
	}

	if prop.MinLength != 0 {
		p.MinLength = uint64(prop.MinLength)
	}

	// todo
	//if len(prop.Items) != 0 {
	//	p.Items = &openapi3.SchemaRef{
	//		//Value: Schema,
	//	}
	//}

	return &p
}

func responsesWithErrors(r openapi3.Responses, actions api.ActionSet) openapi3.Responses {
	r[statusUnauthorized] = &openapi3.ResponseRef{
		Ref: "#/components/responses/" + statusUnauthorized,
	}

	r[statusForbidden] = &openapi3.ResponseRef{
		Ref: "#/components/responses/" + statusForbidden,
	}

	r[statusNotFound] = &openapi3.ResponseRef{
		Ref: "#/components/responses/" + statusNotFound,
	}

	if actions.Has(api.ActionCreate | api.ActionReplace | api.ActionUpdate) {
		r[statusBadRequest] = &openapi3.ResponseRef{
			Ref: "#/components/responses/" + statusBadRequest,
		}
	}

	return r
}
