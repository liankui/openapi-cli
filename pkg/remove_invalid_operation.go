package pkg

import (
	"context"
	"strings"

	"github.com/chaos-io/chaos/core/logs"
	"github.com/getkin/kin-openapi/openapi2"
)

func RemoveInvalidOperation(ctx context.Context, buffer []byte, marshaler func(v interface{}) ([]byte, error), unmarshaler func(data []byte, v interface{}) error) ([]byte, error) {
	lintResult := OpenapiLint(ctx, buffer)
	if lintResult.Valid {
		return buffer, nil
	}

	v2 := &openapi2.T{}
	err := unmarshaler(buffer, v2)
	if err != nil {
		logs.Warnw("failed to unmarshal openapi2 document", "error", err, "buffer", string(buffer))
		return nil, err
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

		logs.Debugw("invalid operation has been deleted", "operation", operation)
	}

	buffer, err = marshaler(v2)
	if err != nil {
		logs.Warnw("failed to marshal openapi2 document", "error", err, "v2", v2)
		return nil, err
	}

	return buffer, nil
}
