package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/watcherwhale/gogl-ci/pkg/gitlab"
)

var (
	TEST_PIPELINE gitlab.Pipeline = gitlab.Pipeline{
		Stages: []string{
			"start",
			"middle",
			"end",
		},
		Jobs: map[string]gitlab.Job{
			".ignored": {
				Name:  ".ignored",
				Needs: nil,
				Rules: []gitlab.Rule{
					{When: "always"},
				},
			},
			"test": {
				Name:  "test",
				Stage: "start",
				Needs: nil,
				Rules: []gitlab.Rule{
					{When: "always"},
				},
			},
			"test2": {
				Name:  "test2",
				Stage: "middle",
				Needs: nil,
				Rules: []gitlab.Rule{
					{When: "always"},
				},
			},
			"test3": {
				Name:  "test3",
				Stage: "end",
				Needs: nil,
				Rules: []gitlab.Rule{
					{When: "always"},
				},
			},
			"test4": {
				Name:  "test4",
				Stage: "end",
				Needs: nil,
				Rules: []gitlab.Rule{
					{When: "never"},
				},
			},
			"test5": {
				Name:  "test5",
				Stage: "end",
				Needs: []gitlab.Need{
					{
						Job: "test",
					},
				},
				Rules: []gitlab.Rule{
					{When: "always"},
				},
			},
			"test6": {
				Name:  "test6",
				Stage: "end",
				Needs: []gitlab.Need{
					{
						Job: "test5",
					},
				},
				Rules: []gitlab.Rule{
					{When: "always"},
				},
			},
		},
	}
)

func TestGraphBuild(t *testing.T) {
	var jg JobGraph
	jg.New(TEST_PIPELINE, make(map[string]string))

	assert.False(t, jg.HasJob(".ignored"))
	assert.True(t, jg.HasJob("test"))
	assert.True(t, jg.HasJob("test2"))
	assert.True(t, jg.HasJob("test3"))
	assert.False(t, jg.HasJob("test4"))
	assert.True(t, jg.HasJob("test5"))
	assert.True(t, jg.HasJob("test6"))
}

func TestGraphDependencies(t *testing.T) {
	var jg JobGraph
	jg.New(TEST_PIPELINE, make(map[string]string))

	assert.True(t, jg.HasDependency("test", "test2"))
	assert.True(t, jg.HasDependency("test", "test3"))
	assert.True(t, jg.HasDependency("test2", "test3"))
	assert.True(t, jg.HasDependency("test", "test6"))

	assert.ElementsMatch(t, jg.GetDependencies("test3"), []string{
		"test",
		"test2",
	})
	assert.ElementsMatch(t, jg.GetDependencies("test2"), []string{
		"test",
	})
	assert.Empty(t, jg.GetDependencies("test"))
	assert.ElementsMatch(t, jg.GetDependencies("test5"), []string{
		"test",
	})
	assert.ElementsMatch(t, jg.GetDependencies("test6"), []string{
		"test",
		"test5",
	})
}
