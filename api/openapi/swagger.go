package openapi

import (
	"net/http"
	"strconv"

	"github.com/fluxynet/gocipe/util"
	"github.com/getkin/kin-openapi/openapi3"
)

const (
	paramID     = "#/components/parameters/id"
	paramLimit  = "#/components/parameters/limit"
	paramOffset = "#/components/parameters/offset"
	paramSort   = "#/components/parameters/sort"
)

var (
	contentTypeJSON = "application/json"
)

var (
	statusOK           = strconv.Itoa(http.StatusOK)
	statusBadRequest   = strconv.Itoa(http.StatusBadRequest)
	statusUnauthorized = strconv.Itoa(http.StatusUnauthorized)
	statusForbidden    = strconv.Itoa(http.StatusForbidden)
	statusNotFound     = strconv.Itoa(http.StatusNotFound)
)

// Swagger is an openapi3 schema with added features
type Swagger struct {
	openapi3.Swagger
}

type Info struct {
	Title          string
	Description    string
	TermsOfService string
	ContactName    string
	ContactURL     string
	ContactEmail   string
	LicenseName    string
	LicenseURL     string
	Version        string
}

// New creates a new swagger document with initialized fields
func New(info Info) *Swagger {
	var swagger Swagger

	swagger.OpenAPI = "3.0.0"
	swagger.Info = &openapi3.Info{
		Title:          info.Title,
		Description:    info.Description,
		TermsOfService: info.TermsOfService,
		Version:        info.Version,
	}

	if info.ContactName != "" || info.ContactURL != "" || info.ContactEmail != "" {
		swagger.Info.Contact = &openapi3.Contact{
			Name:  info.ContactName,
			URL:   info.ContactURL,
			Email: info.ContactEmail,
		}
	}

	if info.LicenseName != "" || info.LicenseURL != "" {
		swagger.Info.License = &openapi3.License{
			Name: info.LicenseName,
			URL:  info.LicenseURL,
		}
	}

	swagger.Components.Parameters = openapi3.ParametersMap{
		"id": &openapi3.ParameterRef{
			Value: &openapi3.Parameter{
				ExtensionProps: openapi3.ExtensionProps{},
				Name:           "id",
				In:             "path",
				Description:    "Unique identifier",
				Required:       true,
				Schema: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type:   "string",
						Format: "uuid",
					},
				},
				Example: "00000000-0000-0000-0000-000000000000",
			},
		},

		"limit": &openapi3.ParameterRef{
			Value: &openapi3.Parameter{
				ExtensionProps: openapi3.ExtensionProps{},
				Name:           "__limit",
				In:             "query",
				Description:    "Limit the number of results",
				Schema: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: "integer",
					},
				},
				Example: 100,
			},
		},

		"offset": &openapi3.ParameterRef{
			Value: &openapi3.Parameter{
				ExtensionProps: openapi3.ExtensionProps{},
				Name:           "__offset",
				In:             "query",
				Description:    "Start listing results from this offset",
				Schema: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: "integer",
					},
				},
				Example: 10,
			},
		},

		"sort": &openapi3.ParameterRef{
			Value: &openapi3.Parameter{
				ExtensionProps: openapi3.ExtensionProps{},
				Name:           "__sort",
				In:             "query",
				Description:    "Comma separated list of fields to sort. Prefix fields with - to sort desc. Example sort=name,-age",
				Schema: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: "string",
					},
				},
				Example: "name,-age",
			},
		},
	}

	swagger.Components.Responses = map[string]*openapi3.ResponseRef{
		statusBadRequest: {
			Value: &openapi3.Response{
				Description: util.Str("Bad Request. Request body is incorrect or incomplete"),
			},
		},

		statusUnauthorized: {
			Value: &openapi3.Response{
				Description: util.Str("Unauthorized. Authentication required"),
			},
		},

		statusForbidden: {
			Value: &openapi3.Response{
				Description: util.Str("Forbidden. This action is not allowed for current user"),
			},
		},

		statusNotFound: {
			Value: &openapi3.Response{
				Description: util.Str("Not found"),
			},
		},
	}

	return &swagger
}
