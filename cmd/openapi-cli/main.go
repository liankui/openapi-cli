package main

import (
	"fmt"
	"os"

	"github.com/chaos-io/chaos/core/logs"
	"github.com/urfave/cli/v2"

	"github.com/liankui/openapi-cli/pkg"
)

func main() {
	app := cli.App{
		Name:                   "openapi-cli",
		Usage:                  "openapi command line is a small utility for parsing and validating openapi(swagger) document",
		UsageText:              "openapi-cli [command] [command options] [arguments...]",
		EnableBashCompletion:   true,
		UseShortOptionHandling: true,
		Commands: []*cli.Command{
			{
				Name:    "lint",
				Aliases: []string{"l"},
				Usage:   "lint swagger/openapi document",
				Action:  pkg.NewLint().Action,
			},
			{
				Name:    "upgrade",
				Aliases: []string{"u"},
				Usage:   "upgrade swagger2 to openapi3",
				Action:  pkg.NewUpgrade().Action,
			},
			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "Show version",
				Action: func(ctx *cli.Context) error {
					v, err := os.ReadFile("../../version")
					if err != nil {
						logs.Warnw("failed to get version", "error", err)
						return err
					}
					fmt.Print(string(v))
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		logs.Error(err)
		return
	}
}
