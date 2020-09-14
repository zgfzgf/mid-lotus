// +build !nodaemon

package daemon

import (
	"context"

	"github.com/zgfzgf/mid-lotus/node"
	"github.com/zgfzgf/mid-lotus/node/config"

	"gopkg.in/urfave/cli.v2"
)

// Cmd is the `mid-lotus daemon` command
var Cmd = &cli.Command{
	Name:  "daemon",
	Usage: "Start a lotus daemon process",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "api",
			Value: ":1234",
		},
	},
	Action: func(cctx *cli.Context) error {
		ctx := context.Background()

		cfg, err := config.FromFile("./config.toml")
		if err != nil {
			return err
		}

		api, err := node.New(ctx, node.Online(), node.Config(cfg))
		if err != nil {
			return err
		}

		return serveRPC(api, cctx.String("api"))
	},
}
