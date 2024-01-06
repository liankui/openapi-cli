package internal

import (
	"net/url"
	"strings"

	"github.com/chaos-io/postman/go/pkg/postman/v2"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/pkg/errors"
)

type Postman postman.Collection

func NewPostman() *Postman {
	return &Postman{}
}

func (p *Postman) Normalize() {

}

func (p *Postman) Covert() (*Openapi3, error) {
	v3 := &Openapi3{
		OpenAPI: "3.0.0",
		Info: &openapi3.Info{
			Title:       p.Info.Name,
			Description: p.Info.Description.String(),
			Version:     "1.0.0",
		},
		Servers: nil,
	}

	var err error

	for _, item := range p.Item {
		if item != nil {
			if v3.Servers == nil {
				v3.Servers, err = p.scrapeURL(item.GetRequest().GetUrl())
				if err != nil {
					continue
				}
			}
		}
	}

	return v3, nil
}

func (p *Postman) scrapeURL(u *postman.Url) (openapi3.Servers, error) {
	if u == nil || u.GetRaw() == "" {
		return nil, errors.New("url is empty")
	}

	fixedURL := u.GetRaw()
	if strings.HasPrefix(u.GetRaw(), "{{") {
		fixedURL = "http://" + fixedURL
	}

	URL, err := url.Parse(fixedURL)
	if err != nil {
		return nil, errors.Errorf("failed to parse url, error: %v", err)
	}

	vars := make(map[string]*openapi3.ServerVariable, len(p.Variable))
	for _, variable := range p.Variable {
		vars[variable.Name] = &openapi3.ServerVariable{
			Enum:    []string{variable.Value},
			Default: variable.Value,
		}
	}

	server := &openapi3.Server{
		URL: URL.String(),
		// Description:,
		Variables: vars,
	}

	return openapi3.Servers{server}, nil
}
