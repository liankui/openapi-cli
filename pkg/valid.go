package pkg

import (
	"os"
	"strings"
	
	"github.com/chaos-io/chaos/core/logs"
	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi3"
	"gopkg.in/yaml.v3"
)

const (
	defaultExt      = false
	defaultDefaults = true
	defaultExamples = true
	defaultPatterns = true
)

var OpenapiVersion string

// TODO valid error
func Valid(filename string, ext, defaults, examples, patterns bool) bool {
	data, err := os.ReadFile(filename)
	if err != nil {
		logs.Warn(err)
		return false
	}
	
	_, unmarshaller, err := GetMarshaller(filename)
	if err != nil {
		logs.Warnw("failed to get marshaller", "error", err)
		return false
	}
	
	var vd struct {
		OpenAPI string `json:"openapi" yaml:"openapi"`
		Swagger string `json:"swagger" yaml:"swagger"`
	}
	if err := unmarshaller(data, &vd); err != nil {
		logs.Warn(err)
		return false
	}
	
	switch {
	case vd.OpenAPI == "3" || strings.HasPrefix(vd.OpenAPI, "3."):
		OpenapiVersion = vd.OpenAPI
		loader := openapi3.NewLoader()
		loader.IsExternalRefsAllowed = ext
		
		doc, err := loader.LoadFromFile(filename)
		if err != nil {
			logs.Warnw("Loading error", "error", err)
			return false
		}
		
		var opts []openapi3.ValidationOption
		// when false, disables schemas' default field validation
		if !defaults {
			opts = append(opts, openapi3.DisableSchemaDefaultsValidation())
		}
		// when false, disables all example schema validation
		if !examples {
			opts = append(opts, openapi3.DisableExamplesValidation())
		}
		// enables visiting other files
		if !patterns {
			opts = append(opts, openapi3.DisableSchemaPatternValidation())
		}
		
		if err = doc.Validate(loader.Context, opts...); err != nil {
			logs.Warnw("Validation error", "error", err)
			return false
		}
	
	case vd.Swagger == "2" || strings.HasPrefix(vd.Swagger, "2."):
		OpenapiVersion = vd.Swagger
		if defaults != defaultDefaults {
			logs.Warn("Flag --defaults is only for OpenAPIv3")
			return false
		}
		if examples != defaultExamples {
			logs.Warn("Flag --examples is only for OpenAPIv3")
			return false
		}
		if ext != defaultExt {
			logs.Warn("Flag --ext is only for OpenAPIv3")
			return false
		}
		if patterns != defaultPatterns {
			logs.Warn("Flag --patterns is only for OpenAPIv3")
			return false
		}
		
		var doc openapi2.T
		if err := yaml.Unmarshal(data, &doc); err != nil {
			logs.Warnw("Loading error", "error", err)
			return false
		}
	
	default:
		logs.Warn("Missing or incorrect 'openapi' or 'swagger' field")
		return false
	}
	
	return true
}

func Valid2(filename string) bool {
	data, err := os.ReadFile(filename)
	if err != nil {
		logs.Warn(err)
		return false
	}
	
	_, unmarshaller, err := GetMarshaller(filename)
	if err != nil {
		logs.Warnw("failed to get marshaller", "error", err)
		return false
	}
	
	var vd struct {
		OpenAPI string `json:"openapi" yaml:"openapi"`
		Swagger string `json:"swagger" yaml:"swagger"`
	}
	if err := unmarshaller(data, &vd); err != nil {
		logs.Warn(err)
		return false
	}
	
	switch {
	case vd.OpenAPI == "3" || strings.HasPrefix(vd.OpenAPI, "3."):
		OpenapiVersion = vd.OpenAPI
	
	case vd.Swagger == "2" || strings.HasPrefix(vd.Swagger, "2."):
		OpenapiVersion = vd.Swagger
	
	default:
		logs.Warn("Missing or incorrect 'openapi' or 'swagger' field")
		return false
	}
	
	return true
}
