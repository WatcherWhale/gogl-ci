package v1alpha1

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/watcherwhale/gogl-ci/pkg/format"
	"github.com/watcherwhale/gogl-ci/pkg/gitlab"
	"github.com/watcherwhale/gogl-ci/pkg/testplan"
)

// TestReleaseBranch_SlashEscape verifies that rule conditions containing \/
// (escaped forward slash) are parsed and evaluated correctly.
func TestReleaseBranch_SlashEscape(t *testing.T) {
	pipeline, err := gitlab.Parse("release-branch/.gitlab-ci.yml")

	require.NoError(t, err)
	require.NotNil(t, pipeline)

	plan, err := testplan.ParseFile("release-branch/plans/release_branch.yaml")

	require.NoError(t, err)
	require.NotNil(t, plan)

	tout := plan.Validate(pipeline)

	t.Logf("\n%s", format.SprintTests(tout))

	assert.True(t, tout.IsTreeSucceeded())
}
