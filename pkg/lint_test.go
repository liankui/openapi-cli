package pkg

import (
	"context"
	"os"
	"testing"

	jsoniter "github.com/json-iterator/go"
)

// FIXME The testcase should be not pass
func TestOpenapiLint(t *testing.T) {
	spec, _ := os.ReadFile("../cmd/openapi-cli/testdata/api-docs-cycle-import1.json")

	type args struct {
		ctx  context.Context
		spec []byte
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test1",
			args: args{
				ctx:  context.Background(),
				spec: spec,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := OpenapiLint(tt.args.ctx, tt.args.spec)
			if err != nil {
				t.Fatalf("get error: %v", err)
			}

			str, err := jsoniter.MarshalIndent(got, "", "    ")
			_ = str
			_ = err
		})
	}
}
