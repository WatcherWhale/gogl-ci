package commands

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"github.com/watcherwhale/gogl-ci/pkg/gitlab"
)

var ListCommand cli.Command = cli.Command{
	Name:      "list",
	Args:      true,
	ArgsUsage: "JOB",
	Category:  "Query",
	Aliases: []string{
		"jobs",
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
		&cli.BoolFlag{
			Name: "mr",
			Value: false,
		},
	},
	Usage:  "Get all jobs for specific conditions.",
	Action: getJobs,
}

func getJobs(ctx *cli.Context) error {
	pipeline, err := gitlab.Parse(ctx.String("file"))
	if err != nil {
		return err
	}

	varMap := map[string]string{
		"CI_COMMIT_BRANCH":  ctx.String("branch"),
		"CI_DEFAULT_BRANCH": ctx.String("default-branch"),
	}

	if ctx.Bool("mr") {
		varMap["CI_PIPELINE_SOURCE"] = "merge_request_event"
	}

	jobs, err := pipeline.GetActiveJobs(varMap)
	if err != nil {
		return err
	}


	if len(jobs) == 0 {
		fmt.Println("[]")
		return nil
	}

	for _, job := range jobs {
		fmt.Println(job.Name)
	}

	return nil
}
