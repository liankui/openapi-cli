package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/liankui/openapi-converter/action"
)

// start
// stop
func main() {
	app := cli.App{
		Name:                   "openapi-cli",
		Usage:                  "openapi command line is a small utility for parsing and validating openapi(swagger) document",
		UsageText:              "openapi-cli command [command options] [arguments...]",
		EnableBashCompletion:   true,
		UseShortOptionHandling: true,
		Commands: []*cli.Command{
			{
				Name:   "lint",
				Usage:  "lint openapi/swagger document",
				Action: action.NewLint().Action(),
			},
			{
				Name:    "from",
				Aliases: []string{"f"},
				Usage:   "Specifies format to convert",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "dir",
						Aliases: []string{"d"},
						Usage:   "the target dir to start, default is the current work dir",
					}},
			},
			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "Show version",
				Action: func(ctx *cli.Context) error {
					fmt.Println("v0.1.0")
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Printf("err=%v\n", err)
		return
	}
}
