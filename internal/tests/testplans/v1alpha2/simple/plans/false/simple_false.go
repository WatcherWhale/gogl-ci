package plans_test

import (
	"github.com/watcherwhale/gogl-ci/pkg/gitlab"
	"github.com/watcherwhale/gogl-ci/pkg/graph"
)

func TestMaster(pipeline gitlab.Pipeline) (bool, string) {
	g, err := graph.NewGraph(pipeline, map[string]string{
		"CI_COMMIT_BRANCH":  "master",
		"CI_DEFAULT_BRANCH": "master",
	})
	if err != nil {
		return false, err.Error()
	}

	err = g.Validate()
	if err != nil {
		return false, err.Error()
	}

	if g.HasJob("job2") {
		return false, "Job2 should not be present"
	}

	return true, ""
}
