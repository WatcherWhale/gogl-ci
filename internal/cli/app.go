package cli

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"github.com/watcherwhale/gitlabci-test/internal/cli/commands"
	"github.com/watcherwhale/gitlabci-test/pkg/api"
)

func InitCli() error {
	app := &cli.App{
		Name:    "gitlabci-test",
		Version: "0.0.0",

		Commands: []*cli.Command{
			&commands.ParseCommand,
		},

		Before: func(ctx *cli.Context) error {
			if ctx.String("token") != "" {
				scheme := "https://"
				if !ctx.Bool("https") {
					scheme = "http://"
				}

				api.GitlabUrl = scheme + ctx.String("instance") + "/api/v4"
				log.Logger.Debug().Msgf("Using gitlab instance: %s", api.GitlabUrl)

				return api.Login(ctx.String("token"))
			}
			return nil
		},

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "file",
				Usage: "The ci root file",
				Value: ".gitlab-ci.yml",
			},
			&cli.StringFlag{
				Name:  "dir",
				Usage: "Set the project directory, to search files",
				Action: func(_ *cli.Context, dir string) error {
					return os.Chdir(dir)
				},
			},
			&cli.StringFlag{
				Name:     "instance",
				Category: "GitLab",
				Usage:    "The gitlab instance url",
				Value:    "gitlab.com",
			},
			&cli.BoolFlag{
				Name:     "https",
				Category: "GitLab",
				Usage:    "Use https to connect to the GitLab Instance",
				Value:    true,
			},
			&cli.StringFlag{
				Name:     "token",
				Category: "GitLab",
				Usage:    "The gitlab api token",
			},
		},
	}

	return app.Run(os.Args)
}
