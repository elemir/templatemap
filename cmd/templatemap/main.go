package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/elemir/templatemap/pkg/plugin"
)

func runPlugin(log *logrus.Logger) cli.ActionFunc {
	return func(c *cli.Context) error {
		plugin, err := plugin.NewPlugin(c.String("nodeid"), c.String("endpoint"), log)
		if err != nil {
			return fmt.Errorf("cannot create plugin: %w", err)
		}
		return plugin.Run()
	}
}

func main() {
	log := logrus.New()
	log.Level = logrus.DebugLevel
	log.Formatter = &logrus.JSONFormatter{}

	app := &cli.App{
		Name: "templatemap",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "nodeid",
				Aliases: []string{"n"},
				Usage:   "node id",
			},
			&cli.StringFlag{
				Name:    "endpoint",
				Value:   "unix:///csi/csi.sock",
				Aliases: []string{"e"},
				Usage:   "CSI endpoint",
			},
		},
		Action: runPlugin(log),
	}

	err := app.Run(os.Args)
	if err != nil {
		log.WithError(err).Fatal("cannot run plugin")
	}
}
