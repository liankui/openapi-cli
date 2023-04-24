package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

// start
// stop
func main() {
	app := cli.App{
		UseShortOptionHandling: true,
		Commands: []*cli.Command{{
			Name:    "from",
			Aliases: []string{"f"},
			Usage:   "Specifies format to convert",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "dir",
					Aliases: []string{"d"},
					Usage:   "the target dir to start, default is the current work dir",
				}},
		}},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Printf("err=%v\n", err)
		return
	}
}
