package plans

import (
	"github.com/watcherwhale/gogl-ci/pkg/gitlab"
	"github.com/watcherwhale/gogl-ci/pkg/graph"
)

func TestJob5Depends(p gitlab.Pipeline, g graph.JobGraph, vars map[string]string) (bool, string) {
	if !g.HasDependency("job1", "job5") {
		return false, "job5 does not depend on job1"
	}

	return true, ""
}

func TestJob5Present(p gitlab.Pipeline, g graph.JobGraph, vars map[string]string) (bool, string) {
	if !g.HasJob("job5") {
		return false, "job5 is not present"
	}
	return true, ""
}

func TestJob1Depends(p gitlab.Pipeline, g graph.JobGraph, vars map[string]string) (bool, string) {
	if !g.HasDependency("job1", "job5") {
		return false, "job5 does not depend on job1"
	}

	if !g.HasDependency("job5", "job1") {
		return false, "job1 does not depend on job5"
	}

	return true, ""
}
