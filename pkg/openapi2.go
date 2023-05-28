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
	"gopkg.in/yaml.v3"
)

type Openapi2 struct {
	Filename     string
	marshaller   func(v interface{}) ([]byte, error)
	unmarshaller func(data []byte, v interface{}) error
}

func NewOpenapi2(filename string) *Openapi2 {
	o2 := &Openapi2{Filename: filename}
	if strings.HasSuffix(filename, ".yaml") || strings.HasSuffix(filename, ".yml") {
		o2.marshaller, o2.unmarshaller = yaml.Marshal, yaml.Unmarshal
	} else if strings.HasSuffix(filename, ".json") {
		o2.marshaller, o2.unmarshaller = jsoniter.Marshal, jsoniter.Unmarshal
	}
	return o2
}

func (o2 *Openapi2) GetOpenapi2(ctx context.Context) (*openapi2.T, error) {
	data, err := os.ReadFile(o2.Filename)
	if err != nil {
		logs.Errorw("failed to read file", "error", err, "path", o2.Filename)
		return nil, err
	}
	
	v2 := &openapi2.T{}
	if err = o2.unmarshaller(data, &v2); err != nil {
		logs.Warnw("failed to parse the openapi", "file", o2.Filename, "error", err)
		return nil, fmt.Errorf("failed to parse the openapi due to %v2", err.Error())
	}
	
	if len(v2.Schemes) == 0 {
		v2.Schemes = []string{"http"}
	}
	
	if len(v2.BasePath) > 0 && len(strings.TrimPrefix(v2.BasePath, "/")) > 0 {
		newPaths := make(map[string]*openapi2.PathItem, len(v2.Paths))
		for key, item := range v2.Paths {
			if key == "/" {
				key = v2.BasePath
			} else {
				key = v2.BasePath + key
			}
			newPaths[key] = item
		}
		v2.BasePath = ""
		v2.Paths = newPaths
	}
	
	return v2, nil
}

func (o2 *Openapi2) UpgradeOpenAPI(ctx context.Context) (*openapi3.T, error) {
	logs.Infow("api upgrade", "file", o2.Filename)
	start := time.Now()
	
	v2, err := o2.GetOpenapi2(ctx)
	if err != nil {
		return nil, err
	}
	
	o2.RemoveInvalidOperation(ctx, v2)
	
	if OpenapiVersion == "2" || strings.HasPrefix(OpenapiVersion, "2.") {
		v3, err := openapi2conv.ToV3(v2)
		if err != nil {
			logs.Warnw("failed to convert swagger2 to openapi3", "error", err)
			return nil, err
		}
		
		buffer, err := jsoniter.Marshal(&v3)
		if err != nil {
			logs.Warnw("failed to marshal openapi3", "error", err)
			return nil, err
		}
		
		newfp := strings.TrimSuffix(o2.Filename, ".json") + strconv.Itoa(time.Now().Nanosecond()) + ".json"
		err = os.WriteFile(newfp, buffer, os.ModePerm)
		if err != nil {
			logs.Warnw("failed to write upgraded api", "error", err)
			return nil, err
		}
		
		logs.Infow("api upgrade successfully", "upgraded version", v3.OpenAPI, "duration", time.Since(start).String())
		
		return v3, nil
	}
	
	logs.Infow("skip api upgrade", "openapi version", OpenapiVersion, "duration", time.Since(start).String())
	
	return nil, nil
}

func (o2 *Openapi2) RemoveInvalidOperation(ctx context.Context, v2 *openapi2.T) {
	buffer, _ := o2.marshaller(v2)
	
	lintResult := OpenapiLint(ctx, buffer)
	if lintResult.Valid {
		return
	}
	
	for _, operation := range lintResult.Operations {
		if operation.Valid {
			continue
		}
		
		operations, found := v2.Paths[operation.Path]
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
		
		logs.Infow("invalid operation has been deleted", "operation", operation)
	}
}
