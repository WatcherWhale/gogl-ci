package commands

import (
	"github.com/urfave/cli/v2"
	"github.com/watcherwhale/gogl-ci/pkg/gitlab"
	"github.com/watcherwhale/gogl-ci/pkg/graph"
)

var ParseCommand cli.Command = cli.Command{
	Name: "parse",
	Action: func(ctx *cli.Context) error {
		pipeline, err := gitlab.Parse(ctx.String("file"))

		if err != nil {
			return err
		}

		_ := graph.CalculateJobGraph(*pipeline)

		return nil
	},
}
