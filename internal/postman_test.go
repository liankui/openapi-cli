package internal

import (
	_ "embed"
	"fmt"
	"reflect"
	"testing"

	"github.com/chaos-io/postman/go/pkg/postman/v2"
	"github.com/getkin/kin-openapi/openapi3"
	jsoniter "github.com/json-iterator/go"
)

var (
	//go:embed testdata/postman/collection.json
	quizFile []byte
	quiz     postman.Collection
)

func init() {
	err := jsoniter.Unmarshal(quizFile, &quiz)
	if err != nil {
		fmt.Errorf("parse postman error: %v", err)
	}
}

func TestPostman_scrapeURL(t *testing.T) {
	type args struct {
		u *postman.Url
	}
	tests := []struct {
		name    string
		p       Postman
		args    args
		want    openapi3.Servers
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.scrapeURL(tt.args.u)
			if (err != nil) != tt.wantErr {
				t.Errorf("scrapeURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("scrapeURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}
