package pkg

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/chaos-io/chaos/core/logs"
	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi2conv"
	"github.com/getkin/kin-openapi/openapi3"
	jsoniter "github.com/json-iterator/go"
)

type Openapi2 struct {
	Filename    string
	marshaler   func(v interface{}) ([]byte, error)
	unmarshaler func(data []byte, v interface{}) error
}

func NewOpenapi2(filename string) *Openapi2 {
	o2 := &Openapi2{Filename: filename}
	o2.marshaler, o2.unmarshaler, _ = GetMarshaller(filename)
	return o2
}

func (v2 *Openapi2) GetOpenapi2(ctx context.Context) (*openapi2.T, error) {
	data, err := os.ReadFile(v2.Filename)
	if err != nil {
		logs.Errorw("failed to read file", "error", err, "path", v2.Filename)
		return nil, err
	}

	o2 := &openapi2.T{}
	if err = v2.unmarshaler(data, &o2); err != nil {
		logs.Warnw("failed to parse the openapi", "file", v2.Filename, "error", err)
		return nil, fmt.Errorf("failed to parse the openapi due to %v2", err.Error())
	}

	if len(o2.Schemes) == 0 {
		o2.Schemes = []string{"http"}
	}

	if len(o2.BasePath) > 0 && len(strings.TrimPrefix(o2.BasePath, "/")) > 0 {
		newPaths := make(map[string]*openapi2.PathItem, len(o2.Paths))
		for key, item := range o2.Paths {
			if key == "/" {
				key = o2.BasePath
			} else {
				key = o2.BasePath + key
			}
			newPaths[key] = item
		}
		o2.BasePath = ""
		o2.Paths = newPaths
	}

	if len(o2.Consumes) > 0 && strings.HasPrefix(o2.Consumes[0], "application/json") {

	}

	return o2, nil
}

func (v2 *Openapi2) UpgradeOpenAPI(ctx context.Context) (*openapi3.T, error) {
	logs.Infow("api upgrade", "file", v2.Filename)
	start := time.Now()

	o2, err := v2.GetOpenapi2(ctx)
	if err != nil {
		return nil, err
	}

	v2.RemoveInvalidOperation(ctx, o2)

	if OpenapiVersion == "2" || strings.HasPrefix(OpenapiVersion, "2.") {
		v3, err := openapi2conv.ToV3(o2)
		if err != nil {
			logs.Warnw("failed to convert swagger2 to openapi3", "error", err)
			return nil, err
		}

		o3 := NormalizeV3(v3)
		buffer, err := jsoniter.MarshalIndent(&o3, "", "  ")
		if err != nil {
			logs.Warnw("failed to marshal openapi3", "error", err)
			return nil, err
		}

		newfp := strings.TrimSuffix(v2.Filename, ".json") + "-" + strconv.FormatInt(time.Now().UnixMilli(), 10) + ".json"
		err = os.WriteFile(newfp, buffer, os.ModePerm)
		if err != nil {
			logs.Warnw("failed to write upgraded api", "error", err)
			return nil, err
		}

		logs.Infow("api upgrade successfully", "file", newfp, "version", v3.OpenAPI, "duration", time.Since(start).String())

		return v3, nil
	}

	logs.Infow("skip api upgrade", "openapi version", OpenapiVersion)

	return nil, nil
}

func (v2 *Openapi2) RemoveInvalidOperation(ctx context.Context, o2 *openapi2.T) {
	buffer, _ := v2.marshaler(o2)

	lintResult := OpenapiLint(ctx, buffer)
	if lintResult.Valid {
		return
	}

	for _, operation := range lintResult.Operations {
		if operation.Valid {
			continue
		}

		operations, found := o2.Paths[operation.Path]
		if !found {
			continue
		}

		// remove invalid operations
		switch strings.ToUpper(operation.Method) {
		case "GET":
			operations.Get = nil
		case "POST":
			operations.Post = nil
		case "PUT":
			operations.Put = nil
		case "DELETE":
			operations.Delete = nil
		case "HEAD":
			operations.Head = nil
		case "OPTIONS":
			operations.Options = nil
		case "PATCH":
			operations.Patch = nil
		default:
			continue
		}

		logs.Infow("delete invalid operation", "operation", operation)
	}
}
