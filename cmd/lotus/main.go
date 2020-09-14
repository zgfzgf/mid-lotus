package main

import (
	"log"
	"os"

	"gopkg.in/urfave/cli.v2"

	"github.com/zgfzgf/mid-lotus/api"
	"github.com/zgfzgf/mid-lotus/api/client"
	"github.com/zgfzgf/mid-lotus/build"
	lcli "github.com/zgfzgf/mid-lotus/cli"
	"github.com/zgfzgf/mid-lotus/daemon"
)

func main() {
	local := []*cli.Command{
		daemon.Cmd,
	}

	app := &cli.App{
		Name:    "lotus",
		Usage:   "Filecoin decentralized storage network client",
		Version: build.Version,
		Metadata: map[string]interface{}{
			"api": lcli.ApiConnector(func() api.API {
				// TODO: get this from repo
				return client.NewRPC("http://127.0.0.1:1234/rpc/v0")
			}),
		},

		Commands: append(local, lcli.Commands...),
	}

	if err := app.Run(os.Args); err != nil {
		log.Println(err)
		return
	}
}
