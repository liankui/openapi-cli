package internal

import "github.com/chaos-io/postman/go/pkg/postman/v2"

type Postman postman.Postman

func NewPostman() *Postman {
	return &Postman{}
}

func (p *Postman) Normalize() {

}

func (p *Postman) Covert() (*Openapi3, error) {
	var v3 *Openapi3
	for _, items := range p.Item {
		if items != nil {

		}
	}

	return v3, nil
}

func (p *Postman) scrapeURL(url *postman.Url) bool {
	if url == nil || url.GetRaw() == "" {
		return false
	}

	return true
}
