package pkg

import (
	"os"
	"path"

	"github.com/chaos-io/chaos/core/logs"
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
		return logs.NewError("NOT specified the target file or dir")
	}

	f, err := os.Stat(dirFile)
	if err != nil {
		logs.Errorw("failed to read file", "error", err, "path", os.Args[1:])
		return nil
	}

	if f.IsDir() {
		files, err := os.ReadDir(dirFile)
		if err != nil {
			logs.Errorw("failed to read dir", "error", err, "path", os.Args[1])
			return err
		}

		for _, file := range files {
			if !file.IsDir() && Valid2(path.Join(dirFile, file.Name())) {
				filePath := path.Join(dirFile, file.Name())
				openapi2 := NewOpenapi2(filePath)
				if _, err := openapi2.UpgradeOpenAPI(c.Context); err != nil {
					logs.Warnw("api upgrade failed", "file", filePath, "error", err)
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
