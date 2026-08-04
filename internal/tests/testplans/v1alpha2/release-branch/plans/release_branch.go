package plans

import (
	"github.com/watcherwhale/gogl-ci/pkg/gitlab"
	"github.com/watcherwhale/gogl-ci/pkg/graph"
)

// TestReleaseBranchJobs verifies that rc and promote jobs are active on a
// release branch and that the tag job is absent.
func TestReleaseBranchJobs(pipeline gitlab.Pipeline) (bool, string) {
	vars := gitlab.NewVariables().
		WithDefaultBranch("master").
		WithBranch("prefix/1.0.x")

	g, err := graph.NewGraph(pipeline, vars)
	if err != nil {
		return false, err.Error()
	}

	if err := g.Validate(); err != nil {
		return false, err.Error()
	}

	if !g.HasJob("release:rc") {
		return false, "release:rc should be present on release branch"
	}
	if !g.HasJob("release:promote") {
		return false, "release:promote should be present on release branch"
	}
	if g.HasJob("release:tag") {
		return false, "release:tag should NOT be present on release branch"
	}

	return true, ""
}

// TestMasterNoReleaseJobs verifies that rc and promote jobs are absent on master.
func TestMasterNoReleaseJobs(pipeline gitlab.Pipeline) (bool, string) {
	vars := gitlab.NewVariables().
		WithDefaultBranch("master").
		WithBranch("master")

	g, err := graph.NewGraph(pipeline, vars)
	if err != nil {
		return false, err.Error()
	}

	if g.HasJob("release:rc") {
		return false, "release:rc should NOT be present on master"
	}
	if g.HasJob("release:promote") {
		return false, "release:promote should NOT be present on master"
	}

	return true, ""
}
