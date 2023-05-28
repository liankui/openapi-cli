package pkg

import (
	"fmt"
	"strings"
	
	jsoniter "github.com/json-iterator/go"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

type OpenapiCli interface {
	Action(c *cli.Context) error
}

func GetMarshaller(filename string) (marshaller func(v interface{}) ([]byte, error), unmarshaller func(data []byte, v interface{}) error, err error) {
	if strings.HasSuffix(filename, ".yaml") || strings.HasSuffix(filename, ".yml") {
		marshaller, unmarshaller = yaml.Marshal, yaml.Unmarshal
	} else if strings.HasSuffix(filename, ".json") {
		marshaller, unmarshaller = jsoniter.Marshal, jsoniter.Unmarshal
	} else {
		err = fmt.Errorf("filename's type is invalid type, filename:%v", filename)
	}
	return
}
