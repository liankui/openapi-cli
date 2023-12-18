package action

import (
	"log/slog"
	"os"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"github.com/liankui/openapi-cli/internal"
)

type PostmanCovert struct{}

func NewPostmanCovert() *PostmanCovert {
	return &PostmanCovert{}
}

func (p *PostmanCovert) Action(c *cli.Context) error {
	dirFile := c.Args().First()
	if dirFile == "" {
		return errors.New("NOT specified the target file")
	}

	f, err := os.Stat(dirFile)
	if err != nil {
		slog.Error("failed to get the file", "error", err, "path", os.Args[1:])
		return nil
	}

	file, err := os.ReadFile(f.Name())
	if err != nil {
		slog.Error("failed to read file", "error", err, "path", os.Args[1:])
		return nil
	}

	postman := internal.NewPostman()
	if err = jsoniter.Unmarshal(file, &postman); err != nil {
		slog.Error("failed to parse file to postman", "error", err, "path", os.Args[1:])
		return nil
	}

	if _, err := postman.Covert(); err != nil {
		slog.Error("failed to covert postman collection", "error", err)
		return err
	}

	return nil
}
