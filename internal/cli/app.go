package cli

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"github.com/watcherwhale/gogl-ci/internal/cli/commands"
	"github.com/watcherwhale/gogl-ci/internal/token"
	"github.com/watcherwhale/gogl-ci/pkg/api"
)

func InitCli() error {
	app := &cli.App{
		Name:        "Gogl",
		Description: "A cli tool for getting insight into your gitlab pipelines",
		Version:     "0.0.0",

		EnableBashCompletion: true,

		Commands: []*cli.Command{
			&commands.TestCommand,

			&commands.DependsCommand,

			&commands.CacheCommad,

			&commands.LoginCommand,

			&commands.ListCommand,
		},

		Before: func(ctx *cli.Context) error {
			switch strings.ToUpper(ctx.String("log-level")) {
			case "TRACE":
				zerolog.SetGlobalLevel(zerolog.TraceLevel)
			case "DEBUG":
				zerolog.SetGlobalLevel(zerolog.DebugLevel)
			case "WARN":
				zerolog.SetGlobalLevel(zerolog.WarnLevel)
			case "ERROR":
				zerolog.SetGlobalLevel(zerolog.ErrorLevel)
			default:
				zerolog.SetGlobalLevel(zerolog.InfoLevel)
			}

			tokenStr := ctx.String("token")
			if tokenStr == "" {
				var err error
				tokenStr, err = token.GetToken(ctx.String("instance"))
				if err != nil {
					tokenStr = ""
				}
			}

			if tokenStr != "" {
				scheme := "https://"
				if !ctx.Bool("https") {
					scheme = "http://"
				}

				api.GitlabUrl = scheme + ctx.String("instance") + "/api/v4"
				log.Logger.Debug().Msgf("Using gitlab instance: %s", api.GitlabUrl)

				return api.Login(tokenStr)
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
				Usage:    "The gitlab api token, can also be set with environment variable GITLAB_TOKEN",
				Value:    os.Getenv("GITLAB_TOKEN"),
			},
			&cli.StringFlag{
				Name:  "log-level",
				Usage: "Set the log level, can also be set with environment variable LOG_LEVEL",
				Value: envOrDefault("LOG_LEVEL", "INFO"),
			},
		},
	}

	return app.Run(os.Args)
}

func envOrDefault(key string, defaultStr string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultStr
	}
	return value
}
