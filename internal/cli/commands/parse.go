package commands

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"github.com/watcherwhale/gitlabci-test/pkg/gitlab"
)

var ParseCommand cli.Command = cli.Command{
	Name: "parse",
	Action: func(ctx *cli.Context) error {
		pipeline, err := gitlab.Parse(ctx.String("file"))

		if err != nil {
			return err
		}

		fmt.Println(pipeline.String())

		return nil
	},
}
