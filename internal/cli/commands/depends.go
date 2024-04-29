package commands

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"github.com/watcherwhale/gogl-ci/pkg/gitlab"
	"github.com/watcherwhale/gogl-ci/pkg/graph"
)

var DependsCommand cli.Command = cli.Command{
	Name:      "depends",
	Args:      true,
	ArgsUsage: "JOB",
	Category:  "Query",
	Aliases: []string{
		"deps",
	},
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "branch",
			Value: "master",
		},
		&cli.StringFlag{
			Name:  "default-branch",
			Value: "master",
		},
	},
	Usage:  "Get all dependencies for a specific job.",
	Action: getDeps,
}

func getDeps(ctx *cli.Context) error {
	pipeline, err := gitlab.Parse(ctx.String("file"))
	if err != nil {
		return err
	}

	varMap := map[string]string{
		"CI_COMMIT_BRANCH":  ctx.String("branch"),
		"CI_DEFAULT_BRANCH": ctx.String("default-branch"),
	}

	var depGraph graph.JobGraph
	depGraph.New(*pipeline, varMap)

	if !depGraph.HasJob(ctx.Args().First()) {
		return fmt.Errorf("%s does not exist in these conditions", ctx.Args().First())
	}

	deps := depGraph.GetDependencies(ctx.Args().First())

	if len(deps) == 0 {
		fmt.Println("[]")
		return nil
	}

	for _, dep := range deps {
		fmt.Println(dep)
	}

	return nil
}
