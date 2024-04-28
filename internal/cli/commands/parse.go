package commands

import (
	"github.com/urfave/cli/v2"
)

var ParseCommand cli.Command = cli.Command{
	Name: "parse",
	Action: func(ctx *cli.Context) error {

		return nil
	},
}
