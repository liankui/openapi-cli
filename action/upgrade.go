package action

import (
	"log/slog"
	"os"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"github.com/liankui/openapi-cli/internal"
)

type Upgrade struct{}

func NewUpgrade() *Upgrade {
	return &Upgrade{}
}

func (u *Upgrade) Action(c *cli.Context) error {
	dirFile := c.Args().First()
	if dirFile == "" {
		return errors.New("NOT specified the target file")
	}

	f, err := os.Stat(dirFile)
	if err != nil {
		slog.Error("failed to read file", "error", err, "path", os.Args[1:])
		return nil
	}

	if f.IsDir() {
		return errors.Errorf("can't support directory")
	} else {
		v2 := internal.NewOpenapi2(f.Name())
		if v2.Valid() {
			if _, err := v2.Upgrade(c.Context); err != nil {
				slog.Warn("api upgrade failed", "file", f.Name(), "error", err)
				return err
			}
		}
	}

	return nil
}
