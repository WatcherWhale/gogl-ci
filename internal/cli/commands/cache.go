package commands

import (
	"github.com/urfave/cli/v2"
	"github.com/watcherwhale/gogl-ci/internal/cache"
)

var CacheCommad cli.Command = cli.Command{
	Name:     "cache",
	Category: "Cache",
	Subcommands: []*cli.Command{
		&cleanCacheCommad,
	},
}

var cleanCacheCommad cli.Command = cli.Command{
	Name:    "clean",
	Aliases: []string{"c"},
	Action: func(ctx *cli.Context) error {
		return cache.CleanCache()
	},
}
