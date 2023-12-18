package action

import (
	"github.com/urfave/cli/v2"
)

type OpenapiCli interface {
	Action(c *cli.Context) error
}
