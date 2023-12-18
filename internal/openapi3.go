package internal

import (
	"github.com/getkin/kin-openapi/openapi3"
)

type Openapi3 struct {
	Extensions map[string]interface{} `json:"-" yaml:"-"`

	OpenAPI      string                        `json:"openapi" yaml:"openapi"` // Required
	Info         *openapi3.Info                `json:"info" yaml:"info"`       // Required
	Servers      openapi3.Servers              `json:"servers,omitempty" yaml:"servers,omitempty"`
	Tags         openapi3.Tags                 `json:"tags,omitempty" yaml:"tags,omitempty"`
	Paths        openapi3.Paths                `json:"paths" yaml:"paths"` // Required
	Components   *openapi3.Components          `json:"components,omitempty" yaml:"components,omitempty"`
	Security     openapi3.SecurityRequirements `json:"security,omitempty" yaml:"security,omitempty"`
	ExternalDocs *openapi3.ExternalDocs        `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
}

func NewOpenapi3(doc *openapi3.T) *Openapi3 {
	return &Openapi3{
		Extensions:   doc.Extensions,
		OpenAPI:      doc.OpenAPI,
		Components:   doc.Components,
		Info:         doc.Info,
		Paths:        doc.Paths,
		Security:     doc.Security,
		Servers:      doc.Servers,
		Tags:         doc.Tags,
		ExternalDocs: doc.ExternalDocs,
	}
}
