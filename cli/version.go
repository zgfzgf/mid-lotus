package cli

import (
	"gopkg.in/urfave/cli.v2"
)

var versionCmd = &cli.Command{
	Name:  "version",
	Usage: "Print version",
	Action: func(context *cli.Context) error {
		cli.VersionPrinter(context)
		return nil
	},
}
