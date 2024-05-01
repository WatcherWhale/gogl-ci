package v1alpha1

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/watcherwhale/gogl-ci/pkg/gitlab"
	"github.com/watcherwhale/gogl-ci/pkg/testplan"
)

func TestSimple_SimpleTrue(t *testing.T) {
	pipeline, err := gitlab.Parse("simple/.gitlab-ci.yml", false)

	require.NoError(t, err)
	require.NotNil(t, pipeline)

	plan, err := testplan.ParseFile("simple/plans/simple_true.yaml")

	require.NoError(t, err)
	require.NotNil(t, plan)

	ok, errorMsg := plan.Validate(pipeline)

	assert.True(t, ok)
	assert.Empty(t, errorMsg)
}

func TestSimple_SimpleFalse(t *testing.T) {
	pipeline, err := gitlab.Parse("simple/.gitlab-ci.yml", false)

	require.NoError(t, err)
	require.NotNil(t, pipeline)

	plan, err := testplan.ParseFile("simple/plans/simple_false.yaml")

	require.NoError(t, err)
	require.NotNil(t, plan)

	ok, errorMsg := plan.Validate(pipeline)

	assert.False(t, ok)
	assert.NotEmpty(t, errorMsg)
}
