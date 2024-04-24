package commands

import (
	"github.com/urfave/cli/v2"
	// "github.com/watcherwhale/gogl-ci/pkg/gitlab"
	// "github.com/watcherwhale/gogl-ci/pkg/graph"
)

var TestCommand cli.Command = cli.Command{
	Name: "test",
	Action: func(ctx *cli.Context) error {

		return nil
	},
}
