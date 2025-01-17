package plans

import (
	"github.com/watcherwhale/gogl-ci/pkg/gitlab"
	"github.com/watcherwhale/gogl-ci/pkg/graph"
)

func TestJob2DependsOnJob1(p gitlab.Pipeline, g graph.JobGraph, vars map[string]string) (bool, string) {
	if !g.HasDependency("job1", "job2") {
		return false, "job2 does not depend on job1"
	}
	return true, ""
}

func TestJob3Depends(p gitlab.Pipeline, g graph.JobGraph, vars map[string]string) (bool, string) {
	if !g.HasDependency("job2", "job3") {
		return false, "job3 does not depend on job2"
	}

	if !g.HasDependency("job1", "job3") {
		return false, "job3 does not depend on job1"
	}

	return true, ""
}

func TestJob4Needs(p gitlab.Pipeline, g graph.JobGraph, vars map[string]string) (bool, string) {
	if len(g.GetDependencies("job4")) != 0 {
		return false, "Job4 has needs, while we expect none"
	}

	return true, ""
}

func TestJob5Absent(p gitlab.Pipeline, g graph.JobGraph, vars map[string]string) (bool, string) {
	if g.HasJob("job5") {
		return false, "job5 is present"
	}
	return true, ""
}
