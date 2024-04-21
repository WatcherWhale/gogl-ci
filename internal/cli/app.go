package cli

import (
	"os"

	"github.com/urfave/cli/v2"
	"github.com/watcherwhale/gitlabci-test/internal/cli/commands"
)

func InitCli() error {
	app := &cli.App{
		Name:    "gitlabci-test",
		Version: "0.0.0",

		Commands: []*cli.Command{
			&commands.ParseCommand,
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
		},
	}

	return app.Run(os.Args)
}
