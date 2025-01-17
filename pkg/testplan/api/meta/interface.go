package meta

import (
	"github.com/watcherwhale/gogl-ci/pkg/format"
	"github.com/watcherwhale/gogl-ci/pkg/gitlab"
)

type TestPlan interface {
	Validate(pipeline *gitlab.Pipeline) format.TestOutput
}
