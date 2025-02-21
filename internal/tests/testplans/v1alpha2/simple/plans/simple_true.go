package plans_test

import (
	"github.com/watcherwhale/gogl-ci/pkg/gitlab"
	"github.com/watcherwhale/gogl-ci/pkg/graph"
)

func TestMaster(pipeline gitlab.Pipeline) (bool, string) {
	vars := gitlab.NewVariables().WithDefaultBranch("master").WithBranch("master")

	g, err := graph.NewGraph(pipeline, vars)
	if err != nil {
		return false, err.Error()
	}

	err = g.Validate()
	if err != nil {
		return false, err.Error()
	}

	if !g.HasJob("job1") {
		return false, "Job1 should be present"
	}

	return true, ""
}
