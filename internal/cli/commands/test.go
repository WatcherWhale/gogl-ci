package commands

import (
	"os"
	"path"
	"regexp"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"github.com/watcherwhale/gogl-ci/pkg/format"
	"github.com/watcherwhale/gogl-ci/pkg/gitlab"
	"github.com/watcherwhale/gogl-ci/pkg/testplan"
)

var TestCommand cli.Command = cli.Command{
	Name:      "test",
	Args:      true,
	Usage:     "Test a TestPlan or a directory of TestPlans against a pipeline.",
	ArgsUsage: "[TestPlan | Dir 1] [TestPlan | Dir 2] ...",
	Category:  "Testing",
	Action: func(ctx *cli.Context) error {
		pipeline, err := gitlab.Parse(ctx.String("file"))
		if err != nil {
			return err
		}

		files := make([]string, 0)
		for _, file := range ctx.Args().Slice() {
			stat, err := os.Stat(file)
			if err != nil {
				return err
			}

			if stat.IsDir() {
				dirFiles, err := readDir(file)
				if err != nil {
					return err
				}
				files = append(files, dirFiles...)
			} else {
				files = append(files, file)
			}
		}

		exitCode := 0
		for _, file := range files {
			plan, err := testplan.ParseFile(file)
			if err != nil {
				log.Debug().Err(err).Msgf("Failed to parse %s", file)
				continue
			}

			output := plan.Validate(pipeline)
			format.PrintTests(output)

			if !output.IsTreeSucceeded() {
				exitCode = 1
			}
		}

		if exitCode == 0 {
			log.Info().Msg("All tests succeeded")
		}
		os.Exit(exitCode)
		return nil
	},
}

func readDir(dir string) ([]string, error) {
	fileEntries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	files := make([]string, 0)
	yamlRegex := regexp.MustCompile(`.*\.ya?ml$`)
	for _, file := range fileEntries {
		if yamlRegex.Match([]byte(file.Name())) {
			files = append(files, path.Join(dir, file.Name()))
		}
	}

	return files, nil
}
