package commands

import (
	"github.com/urfave/cli/v2"
	"github.com/watcherwhale/gogl-ci/internal/token"
)

var LoginCommand cli.Command = cli.Command{
	Name:     "login",
	Category: "GitLab",
	Args:     true,
	Action:   loginAction,
	Usage:    "gogl login <instance> <token>",
}

func loginAction(ctx *cli.Context) error {
	return token.SaveToken(ctx.Args().Get(0), ctx.Args().Get(1))
}
