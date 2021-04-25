package openapi

import (
	"net/http"

	"github.com/fluxynet/gocipe/api"
	"github.com/getkin/kin-openapi/openapi3"
)

// Upload defines media uploading paths having a POST and a DELETE
type Upload interface {
	Path() string
	Types() []string
}

func (s *Swagger) AddUpload(u Upload, actions api.ActionSet) *Swagger {
	var path = u.Path()

	if actions.Has(api.ActionCreate) {
		var (
			t     = u.Types()
			media = make(map[string]*openapi3.MediaType, len(t))
		)

		for i := range t {
			media[t[i]] = &openapi3.MediaType{
				Schema: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type:   "string",
						Format: "binary",
					},
				},
			}
		}

		s.AddOperation(path, http.MethodPost, &openapi3.Operation{
			Description: "Create a new upload content",
			RequestBody: &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					ExtensionProps: openapi3.ExtensionProps{},
					Description:    "",
					Required:       false,
					Content:        media,
				},
			},
			Responses: nil,
		})
	}

	if actions.Has(api.ActionDelete) {
		s.AddOperation(path+"/{id}", http.MethodDelete, &openapi3.Operation{
			Description: "Delete an existing asset",
			Parameters: openapi3.Parameters{
				&openapi3.ParameterRef{Ref: "#/components/parameters/id"},
			},
			Responses: nil,
		})
	}

	return s
}
