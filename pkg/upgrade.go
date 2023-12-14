package pkg

import (
	"log/slog"
	"os"
	"path"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

type Upgrade struct {
}

func NewUpgrade() *Upgrade {
	return &Upgrade{}
}

func (u *Upgrade) Action(c *cli.Context) error {
	dirFile := c.Args().First()
	if dirFile == "" {
		return errors.New("NOT specified the target file or dir")
	}

	f, err := os.Stat(dirFile)
	if err != nil {
		slog.Error("failed to read file", "error", err, "path", os.Args[1:])
		return nil
	}

	if f.IsDir() {
		files, err := os.ReadDir(dirFile)
		if err != nil {
			slog.Error("failed to read dir", "error", err, "path", os.Args[1])
			return err
		}

		for _, file := range files {
			if !file.IsDir() && Valid2(path.Join(dirFile, file.Name())) {
				filePath := path.Join(dirFile, file.Name())
				openapi2 := NewOpenapi2(filePath)
				if _, err := openapi2.UpgradeOpenAPI(c.Context); err != nil {
					slog.Warn("api upgrade failed", "file", filePath, "error", err)
				}
			}
		}
	} else {
		if Valid2(dirFile) {
			openapi2 := NewOpenapi2(dirFile)
			if _, err := openapi2.UpgradeOpenAPI(c.Context); err != nil {
				return err
			}
		}
	}

	return nil
}
