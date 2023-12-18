package internal

import (
	"log/slog"
	"os"
	"strings"

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
		slog.Warn("failed to read file", "file", filename, "error", err)
		return false
	}

	v2 := Openapi2{Filename: filename}
	_, unmarshaller, err := v2.GetMarshaller()
	if err != nil {
		slog.Warn("failed to get marshaller", "error", err)
		return false
	}

	var vd struct {
		OpenAPI string `json:"openapi" yaml:"openapi"`
		Swagger string `json:"swagger" yaml:"swagger"`
	}

	if err := unmarshaller(data, &vd); err != nil {
		slog.Warn("failed to unmarshaler", "error", err)
		return false
	}

	switch {
	case vd.OpenAPI == "3" || strings.HasPrefix(vd.OpenAPI, "3."):
		OpenapiVersion = vd.OpenAPI
		loader := openapi3.NewLoader()
		loader.IsExternalRefsAllowed = ext

		doc, err := loader.LoadFromFile(filename)
		if err != nil {
			slog.Warn("Loading error", "error", err)
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
			slog.Warn("Validation error", "error", err)
			return false
		}

	case vd.Swagger == "2" || strings.HasPrefix(vd.Swagger, "2."):
		OpenapiVersion = vd.Swagger
		if defaults != defaultDefaults {
			slog.Warn("Flag --defaults is only for OpenAPIv3")
			return false
		}
		if examples != defaultExamples {
			slog.Warn("Flag --examples is only for OpenAPIv3")
			return false
		}
		if ext != defaultExt {
			slog.Warn("Flag --ext is only for OpenAPIv3")
			return false
		}
		if patterns != defaultPatterns {
			slog.Warn("Flag --patterns is only for OpenAPIv3")
			return false
		}

		var doc openapi2.T
		if err := yaml.Unmarshal(data, &doc); err != nil {
			slog.Warn("Loading error", "error", err)
			return false
		}

	default:
		slog.Warn("Missing or incorrect 'openapi' or 'swagger' field")
		return false
	}

	return true
}
