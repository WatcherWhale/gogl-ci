package meta

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetApiKind(t *testing.T) {
	yamlSource := `
kind: TestPlan
apiVersion: v1alpha1`

	kind, err := GetAPiKind([]byte(yamlSource))
	require.NoError(t, err)

	assert.Equal(t, "TestPlan", kind.Kind)
	assert.Equal(t, "v1alpha1", kind.Version)
}
