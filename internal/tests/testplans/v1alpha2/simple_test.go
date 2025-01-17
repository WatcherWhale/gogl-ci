package v1alpha1

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/watcherwhale/gogl-ci/pkg/format"
	"github.com/watcherwhale/gogl-ci/pkg/gitlab"
	"github.com/watcherwhale/gogl-ci/pkg/testplan"
)

func TestSimple_SimpleTrue(t *testing.T) {
	pipeline, err := gitlab.Parse("simple/.gitlab-ci.yml")

	require.NoError(t, err)
	require.NotNil(t, pipeline)

	plan, err := testplan.ParseFile("simple/plans/simple_true.yaml")

	require.NoError(t, err)
	require.NotNil(t, plan)

	tout := plan.Validate(pipeline)

	t.Logf("\n%s", format.SprintTests(tout))

	assert.True(t, tout.IsTreeSucceeded())
}

func TestSimple_SimpleFalse(t *testing.T) {
	pipeline, err := gitlab.Parse("simple/.gitlab-ci.yml")

	require.NoError(t, err)
	require.NotNil(t, pipeline)

	plan, err := testplan.ParseFile("simple/plans/simple_false.yaml")

	require.NoError(t, err)
	require.NotNil(t, plan)

	tout := plan.Validate(pipeline)
	t.Logf("\n%s", format.SprintTests(tout))

	assert.False(t, tout.IsTreeSucceeded())
}
