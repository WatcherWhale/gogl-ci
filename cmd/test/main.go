package main

import (
	"fmt"

	"github.com/watcherwhale/gitlabci-test/pkg/gitlab"
)

func main() {
	pl, err := gitlab.Parse("test/pipelines/simpleJob.yaml")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v", pl.String())
}
