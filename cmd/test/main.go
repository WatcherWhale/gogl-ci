package main

import (
	"fmt"

	"github.com/watcherwhale/gitlabci-test/pkg/parser"
)

func main() {
	pl, err := parser.Parse("test/pipelines/simpleJob.yaml")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v", pl.Jobs[0].String())
}
