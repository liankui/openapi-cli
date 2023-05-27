package pkg

import (
	"context"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
	
	"github.com/chaos-io/chaos/core/logs"
	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi2conv"
	"github.com/getkin/kin-openapi/openapi3"
	jsoniter "github.com/json-iterator/go"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

type Upgrade struct {
}

func NewUpgrade() *Upgrade {
	return &Upgrade{}
}

func (u *Upgrade) Action() func(c *cli.Context) error {
	return func(c *cli.Context) error {
		folder := c.String("dir")
		if folder == "" {
			folder = "./testdata"
		}
		logs.Debugw("get the folder", "folder", folder)
		
		files, err := os.ReadDir(folder)
		if err != nil {
			logs.Errorw("failed to read dir", "error", err, "path", os.Args[1])
			return err
		}
		
		var marshaller func(v interface{}) ([]byte, error)
		var unmarshaller func(data []byte, v interface{}) error
		
		for _, file := range files {
			if !file.IsDir() && (strings.HasSuffix(file.Name(), ".yaml") || strings.HasSuffix(file.Name(), ".yml") || strings.HasSuffix(file.Name(), ".json")) {
				filePath := path.Join(folder, file.Name())
				data, err := os.ReadFile(filePath)
				if err != nil {
					logs.Errorw("failed to read file", "error", err, "path", filePath)
					return err
				}
				
				if strings.HasSuffix(filePath, ".yaml") || strings.HasSuffix(filePath, ".yml") {
					marshaller, unmarshaller = yaml.Marshal, yaml.Unmarshal
				} else if strings.HasSuffix(filePath, ".json") {
					marshaller, unmarshaller = jsoniter.Marshal, jsoniter.Unmarshal
				}
				
				oa2 := openapi2.T{}
				if err = unmarshaller(data, &oa2); err != nil {
					logs.Warnw("failed to parse the openapi", "file", file.Name(), "error", err)
					return fmt.Errorf("failed to parse the openapi due to %v", err.Error())
				}
				
				if !strings.HasPrefix(oa2.Info.Version, "3.0") {
					// try to upgrade openapi to latest version
					logs.Info("try to upgrade openapi to latest version due to invalid openapi version!")
					_, err := UpgradeOpenAPI(c.Context, filePath, data, marshaller, unmarshaller)
					if err != nil {
						logs.Warnw("failed to upgrade openAPI to latest version", "error", err)
						return fmt.Errorf("failed to upgrade openAPI to latest version due to %v", err.Error())
					}
					
				}
			}
		}
		
		return nil
	}
}

func UpgradeOpenAPI(ctx context.Context, filePath string, buffer []byte, marshaler func(v interface{}) ([]byte, error), unmarshaler func(data []byte, v interface{}) error) (*openapi3.T, error) {
	start := time.Now()
	var err error
	
	buffer, err = RemoveInvalidOperation(ctx, buffer, marshaler, unmarshaler)
	if err != nil {
		return nil, err
	}
	
	v2 := &openapi2.T{}
	err = unmarshaler(buffer, v2)
	if err != nil {
		logs.Warnw("failed to unmarshal openapi2 document", "error", err, "buffer", string(buffer))
		return nil, err
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
	
	v3, err := openapi2conv.ToV3(v2)
	if err != nil {
		logs.Warnw("failed to convert swagger2 to openapi3", "error", err)
		return nil, err
	}
	
	buffer, err = jsoniter.Marshal(&v3)
	if err != nil {
		logs.Warnw("failed to marshal openapi3", "error", err)
		return nil, err
	}
	
	newfp := strings.TrimSuffix(filePath, ".json") + strconv.Itoa(time.Now().Nanosecond()) + ".json"
	err = os.WriteFile(newfp, buffer, os.ModePerm)
	if err != nil {
		logs.Warnw("failed to write upgraded api", "error", err)
		return nil, err
	}
	
	logs.Infow("api upgrade successfully", "upgraded version", v3.Info.Version, "duration", time.Since(start).String())
	
	return v3, nil
}
