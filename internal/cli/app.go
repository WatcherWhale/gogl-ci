package cli

import (
	"os"

	"github.com/urfave/cli"
)

func InitCli() error {
	app := &cli.App{
		Name:    "gitlabci-test",
		Version: "0.0.0",

		Commands: []cli.Command{},
	}

	return app.Run(os.Args)
}
