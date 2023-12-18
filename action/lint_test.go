package action

import (
	"context"
	"os"
	"testing"

	jsoniter "github.com/json-iterator/go"

	"github.com/liankui/openapi-cli/internal"
)

// FIXME The testcase should be not pass
func TestOpenapiLint(t *testing.T) {
	spec, _ := os.ReadFile("../cmd/openapi-cli/testdata/api-docs-cycle-import1.json")

	type args struct {
		spec []byte
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test1",
			args: args{
				spec: spec,
			},
		},
	}

	ctx := context.Background()
	v2 := &internal.Openapi2{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := v2.Lint(ctx, tt.args.spec)
			if err != nil {
				t.Fatalf("get error: %v", err)
			}

			str, err := jsoniter.MarshalIndent(got, "", "    ")
			_ = str
			_ = err
		})
	}
}
