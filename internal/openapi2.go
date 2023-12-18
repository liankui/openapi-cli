package internal

import (
	"context"
	"log/slog"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/daveshanley/vacuum/model"
	"github.com/daveshanley/vacuum/motor"
	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi2conv"
	"github.com/getkin/kin-openapi/openapi3"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/liankui/openapi-cli/action"
)

type Openapi2 struct {
	Filename string
}

func NewOpenapi2(filename string) *Openapi2 {
	return &Openapi2{Filename: filename}
}

func (v2 *Openapi2) Normalize() (*openapi2.T, error) {
	data, err := os.ReadFile(v2.Filename)
	if err != nil {
		slog.Error("failed to read file", "error", err, "path", v2.Filename)
		return nil, err
	}

	o2 := &openapi2.T{}
	_, unmarshaler, _ := v2.GetMarshaller()
	if err = unmarshaler(data, &o2); err != nil {
		slog.Warn("failed to parse the openapi", "file", v2.Filename, "error", err)
		return nil, errors.Errorf("failed to parse the openapi due to %v", err)
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

	return o2, nil
}

func (v2 *Openapi2) Upgrade(ctx context.Context) (*openapi3.T, error) {
	slog.Info("api upgrade", "file", v2.Filename)
	start := time.Now()

	o2, err := v2.Normalize()
	if err != nil {
		return nil, err
	}

	v2.RemoveInvalidOperation(ctx, o2)

	if OpenapiVersion == "2" || strings.HasPrefix(OpenapiVersion, "2.") {
		v3, err := openapi2conv.ToV3(o2)
		if err != nil {
			slog.Warn("failed to convert swagger2 to openapi3", "error", err)
			return nil, err
		}

		o3 := NewOpenapi3(v3)
		buffer, err := jsoniter.MarshalIndent(&o3, "", "  ")
		if err != nil {
			slog.Warn("failed to marshal openapi3", "error", err)
			return nil, err
		}

		newFile := strings.TrimSuffix(v2.Filename, ".json") + "-" + strconv.FormatInt(time.Now().UnixMilli(), 10) + ".json"
		err = os.WriteFile(newFile, buffer, os.ModePerm)
		if err != nil {
			slog.Warn("failed to write upgraded api", "error", err)
			return nil, err
		}

		slog.Info("successfully", "newFile", newFile, "duration", time.Since(start).String())

		return v3, nil
	}

	slog.Info("skip api upgrade", "openapi version", OpenapiVersion)

	return nil, nil
}

func (v2 *Openapi2) RemoveInvalidOperation(ctx context.Context, o2 *openapi2.T) {
	marshaller, _, _ := v2.GetMarshaller()
	buffer, _ := marshaller(o2)

	lintResult, err := v2.Lint(ctx, buffer)
	if err != nil {
		slog.Warn("[RemoveInvalidOperation] openapi lint error", "error", err)
		return
	}
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

		slog.Info("delete invalid operation", "operation", operation)
	}
}

func (v2 *Openapi2) GetMarshaller() (marshaller func(v interface{}) ([]byte, error), unmarshaller func(data []byte, v interface{}) error, err error) {
	if strings.HasSuffix(v2.Filename, ".yaml") || strings.HasSuffix(v2.Filename, ".yml") {
		marshaller, unmarshaller = yaml.Marshal, yaml.Unmarshal
	} else if strings.HasSuffix(v2.Filename, ".json") {
		marshaller, unmarshaller = jsoniter.Marshal, jsoniter.Unmarshal
	} else {
		err = errors.Errorf("filename's type is invalid, file:%v", v2.Filename)
	}
	return
}

func (v2 *Openapi2) Lint(ctx context.Context, spec []byte) (*LintResult, error) {
	result := motor.ApplyRulesToRuleSet(&motor.RuleSetExecution{
		RuleSet:         action.LintRules,
		Spec:            spec,
		CustomFunctions: map[string]model.RuleFunction{},
	})

	if result.Index == nil {
		return &LintResult{Valid: false}, nil
	}

	if len(result.Errors) > 0 {
		return &LintResult{Valid: false}, errors.Errorf("apply rule get errors, errors: %v", result.Errors)
	}

	operations := result.Index.GetAllPaths()

	lintResult := &LintResult{
		Operations: make([]*LintOperationResult, 0, len(operations)),
		Valid:      true,
	}

	for path, operation := range operations {
		for method := range operation {
			operationResult := &LintOperationResult{
				Path:   path,
				Method: method,
				Valid:  true,
			}

			operationPath := strings.Join([]string{"$.paths", path, method, "parameters"}, ".")
			for _, _result := range result.Results {
				if _result.Path != operationPath {
					continue
				}

				operationResult.Valid = false
				operationResult.Description = _result.Rule.Description
				operationResult.HowToFix = _result.Rule.HowToFix
				if _result.StartNode != nil {
					operationResult.StartLine = int32(_result.StartNode.Line)
				}
				if _result.EndNode != nil {
					operationResult.EndLine = int32(_result.EndNode.Line)
				}

				lintResult.Valid = false
				break
			}

			lintResult.Operations = append(lintResult.Operations, operationResult)
		}
	}

	sort.Sort(LintOperationResults(lintResult.Operations))

	return lintResult, nil
}

func (v2 *Openapi2) Valid() bool {
	file, err := os.ReadFile(v2.Filename)
	if err != nil {
		slog.Warn("failed to read file", "file", v2.Filename, "error", err)
		return false
	}

	var vd struct {
		OpenAPI string `json:"openapi" yaml:"openapi"`
		Swagger string `json:"swagger" yaml:"swagger"`
	}

	_, unmarshaller, _ := v2.GetMarshaller()
	if err := unmarshaller(file, &vd); err != nil {
		slog.Warn("failed to get unmarshaler", "error", err)
		return false
	}

	switch {
	case vd.OpenAPI == "3" || strings.HasPrefix(vd.OpenAPI, "3."):
		OpenapiVersion = vd.OpenAPI
	case vd.Swagger == "2" || strings.HasPrefix(vd.Swagger, "2."):
		OpenapiVersion = vd.Swagger
	default:
		slog.Warn("missing or incorrect 'openapi' or 'swagger' field")
		return false
	}

	return true
}
