package api

import (
	"fmt"

	"github.com/creasty/defaults"
	"github.com/rs/zerolog/log"
	"github.com/watcherwhale/gogl-ci/pkg/gitlab"
	"github.com/watcherwhale/gogl-ci/pkg/graph"
	"github.com/watcherwhale/gogl-ci/pkg/testplan/api/meta"
	"gopkg.in/yaml.v2"
)

type TestPlan struct {
	meta.ApiKind
	Metadata meta.TestPlanMeta `yaml:"metdata"`
	Spec     TestPlanSpec      `yaml:"spec"`
}

type TestPlanSpec struct {
	Pipeline Pipeline   `yaml:"pipeline"`
	Tests    []TestCase `yaml:"tests"`
}

type Pipeline struct {
	Branch    string            `yaml:"branch,omitempty"`
	Tag       string            `yaml:"tag,omitempty"`
	Variables map[string]string `yaml:"variables"`
}

type TestCase struct {
	Name string `yaml:"name"`
	Job  string `yaml:"job"`

	Present bool `yaml:"present" default:"true"`

	// If empty array, check that job depends on nothing
	// If nil, ignore
	DependsOn []string `yaml:"dependsOn"`
}

func LoadPlan(yamlSource []byte) (*TestPlan, error) {
	var plan TestPlan
	err := defaults.Set(&plan)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlSource, &plan)
	if err != nil {
		return nil, err
	}

	log.Trace().Msgf("%s", plan.Spec.Tests[0].DependsOn)

	return &plan, nil
}

func (plan *TestPlan) Validate(pipeline *gitlab.Pipeline) (bool, string) {
	g := graph.CalculateJobGraph(*pipeline)

	status := true
	message := ""

	logger := log.Logger

	for _, tc := range plan.Spec.Tests {
		logger.Trace().Msgf(tc.Name)
		logger.Trace().Msgf("%v", g.HasJob(tc.Job) != tc.Present)
		if g.HasJob(tc.Job) != tc.Present {
			status = false
			if tc.Present {
				message += fmt.Sprintf("- %s: '%s' not found in pipeline\n", tc.Name, tc.Job)
			} else {
				message += fmt.Sprintf("- %s: '%s' has been found in pipeline\n", tc.Name, tc.Job)
			}

			continue
		}

		logger.Trace().Msgf("%v",!tc.Present && tc.Present == g.HasJob(tc.Job))
		if !tc.Present && tc.Present == g.HasJob(tc.Job) {
			continue
		}

		logger.Trace().Msgf("%v", true )
		logger.Trace().Msgf("%v",!tc.Present)
		if !tc.Present {
			status = false
			message += fmt.Sprintf("- %s: cannot validate dependencies while job is not present\n", tc.Name)
			continue
		}

		logger.Trace().Msgf("%v", true )
		logger.Trace().Msgf("%v", tc.DependsOn != nil )
		logger.Trace().Msgf("%v", len(tc.DependsOn) == 0 )
		logger.Trace().Msgf("%v", tc.DependsOn )

		if tc.DependsOn != nil {
			for _, dep := range tc.DependsOn {
				logger.Trace().Msgf("%v", dep )
				if !g.HasJob(dep) {
					logger.Trace().Msgf("%v", false )
					status = false
					message += fmt.Sprintf("- %s: %s is not present in pipeline\n", tc.Name, dep)
					continue
				}

				if !g.HasDependency(tc.Job, dep) {
					logger.Trace().Msgf("%v", false )
					status = false
					message += fmt.Sprintf("- %s: %s does not depend on %s\n", tc.Name, tc.Job, dep)
					continue
				}
				logger.Trace().Msgf("%v", true )
			}
		} else if len(tc.DependsOn) == 0 {
			if len(g.GetJob(tc.Job).Needs) != 0 {
				status = false
				message += fmt.Sprintf("- %s: %s has dependencies defined\n", tc.Name, tc.Job)
				continue
			}
		}
	}

	return status, message
}
