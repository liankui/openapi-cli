package main

import (
	_ "embed"
	"os"
	"strings"

	"github.com/chaos-io/chaos/core/logs"
	"github.com/urfave/cli/v2"

	"github.com/liankui/openapi-cli/pkg"
)

//go:embed version
var version string

func main() {
	app := cli.App{
		Name:            "openapi-cli",
		Usage:           "openapi command line is a small utility for parsing and validating openapi(swagger) document",
		UsageText:       "openapi-cli [command] file/directory",
		HideHelp:        true,
		HideHelpCommand: false,
		Version:         strings.TrimSpace(version),
		Commands: []*cli.Command{
			{
				Name:    "lint",
				Aliases: []string{"l"},
				Usage:   "lint swagger/openapi",
				Action:  pkg.NewLint().Action,
			},
			{
				Name:    "upgrade",
				Aliases: []string{"u"},
				Usage:   "upgrade swagger2 to openapi3",
				Action:  pkg.NewUpgrade().Action,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		logs.Error(err)
		return
	}
}
