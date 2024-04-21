package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/watcherwhale/gitlabci-test/internal/cli"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	err := cli.InitCli()

	if err != nil {
		log.Logger.Error().Err(err).Msg("Exited with an error")
		os.Exit(1)
		return
	}
}
