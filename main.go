package main

import (
	"fmt"
	"os"

	"github.com/mritd/logger"
	"github.com/urfave/cli/v2"
)

var (
	version   string
	buildDate string
	commitID  string
)

func main() {
	app := &cli.App{
		Name:    "bark-server",
		Usage:   "Push Server For Bark",
		Version: fmt.Sprintf("%s %s %s", version, commitID, buildDate),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "addr",
				Usage:   "Server Listen Address",
				EnvVars: []string{"BARK_SERVER_ADDRESS"},
				Value:   "0.0.0.0:8080",
			},
			&cli.StringFlag{
				Name:    "data",
				Usage:   "Server Data Storage Dir",
				EnvVars: []string{"BARK_SERVER_DATA_DIR"},
				Value:   "/data",
			},
			&cli.BoolFlag{
				Name:    "debug",
				Usage:   "Enable Debug Level Log",
				EnvVars: []string{"BARK_SERVER_DEBUG"},
				Value:   false,
			},
		},
		Authors: []*cli.Author{
			{Name: "mritd", Email: "mritd@linux.com"},
			{Name: "Finb", Email: "to@day.app"},
		},
		Action: func(c *cli.Context) error {
			if c.Bool("debug") {
				logger.SetDevelopment()
			}

			databaseSetup(c.String("data"))
			apns2Setup()
			routerSetup(c.Bool("debug"))
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
}
