package api

import (
	"fmt"

	"github.com/creasty/defaults"
	"github.com/watcherwhale/gogl-ci/pkg/format"
	"github.com/watcherwhale/gogl-ci/pkg/gitlab"
	"github.com/watcherwhale/gogl-ci/pkg/graph"
	"github.com/watcherwhale/gogl-ci/pkg/testplan/api/meta"
	"gopkg.in/yaml.v3"
)

type TestPlan struct {
	meta.ApiKind
	Metadata meta.TestPlanMeta `yaml:"metadata"`
	Spec     TestPlanSpec      `yaml:"spec"`
}

type TestPlanSpec struct {
	Pipeline Pipeline   `yaml:"pipeline"`
	Tests    []TestCase `yaml:"tests"`
}

type Pipeline struct {
	DefaultBranch string            `yaml:"defaultBranch,omitempty"`
	Branch        string            `yaml:"branch,omitempty"`
	Tag           string            `yaml:"tag,omitempty"`
	MR            *bool             `yaml:"mr"`
	Variables     map[string]string `yaml:"variables"`
}

func (p Pipeline) isMr() bool {
	return p.MR != nil && *p.MR
}

type TestCase struct {
	Name string `yaml:"name"`
	Job  string `yaml:"job"`

	Present *bool `yaml:"present" default:"true"`

	// If empty array, check that job depends on nothing
	// If nil, ignore
	DependsOn []string `yaml:"dependsOn"`
}

func (tc TestCase) isPresent() bool {
	return tc.Present == nil || *tc.Present
}

func LoadPlan(yamlSource []byte) (meta.TestPlan, error) {
	var plan TestPlan
	err := defaults.Set(&plan)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlSource, &plan)
	if err != nil {
		return nil, err
	}

	return &plan, nil
}

func (plan *TestPlan) BuildVariables() map[string]string {
	variables := plan.Spec.Pipeline.Variables

	if variables == nil {
		variables = make(map[string]string)
	}

	variables["CI_DEFAULT_BRANCH"] = plan.Spec.Pipeline.DefaultBranch

	if plan.Spec.Pipeline.Branch != "" {
		variables["CI_COMMIT_BRANCH"] = plan.Spec.Pipeline.Branch
	}

	if plan.Spec.Pipeline.isMr() {
		variables["CI_PIPELINE_SOURCE"] = "merge_request_event"
	}

	if plan.Spec.Pipeline.Tag != "" {
		variables["CI_COMMIT_TAG"] = plan.Spec.Pipeline.Tag
	}

	return variables
}

func (plan *TestPlan) Validate(pipeline *gitlab.Pipeline) format.TestOutput {
	tout := format.TestOutput{
		Name:      plan.Metadata.Name,
		Succeeded: true,
		SubTests:  make([]format.TestOutput, len(plan.Spec.Tests)),
	}

	var g graph.JobGraph
	err := g.New(*pipeline, plan.BuildVariables())
	if err != nil {
		tout.Succeeded = false
		tout.Message = err.Error()
		return tout
	}

	err = g.Validate()
	if err != nil {
		tout.Succeeded = false
		tout.Message = fmt.Sprintf("error while validating job needs: %v", err)
		return tout
	}

	for i, tc := range plan.Spec.Tests {
		if g.HasJob(tc.Job) != tc.isPresent() {
			if tc.isPresent() {
				tout.SubTests[i] = format.TestOutput{
					Name:      tc.Name,
					Succeeded: false,
					Message:   fmt.Sprintf("%s: '%s' not found in pipeline", tc.Name, tc.Job),
				}
			} else {
				tout.SubTests[i] = format.TestOutput{
					Name:      tc.Name,
					Succeeded: false,
					Message:   fmt.Sprintf("%s: '%s' has been found in pipeline", tc.Name, tc.Job),
				}
			}

			continue
		}

		if !tc.isPresent() && tc.isPresent() == g.HasJob(tc.Job) {
			tout.SubTests[i] = format.TestOutput{
				Name:      tc.Name,
				Succeeded: true,
			}
			continue
		}

		if !tc.isPresent() {
			tout.SubTests[i] = format.TestOutput{
				Name:      tc.Name,
				Succeeded: false,
				Message:   "cannot validate dependencies while job is not present",
			}
			continue
		}

		if tc.DependsOn != nil {
			for _, dep := range tc.DependsOn {
				if !g.HasJob(dep) {
					tout.SubTests[i] = format.TestOutput{
						Name:      tc.Name,
						Succeeded: false,
						Message:   fmt.Sprintf("%s is not present in pipeline", dep),
					}
					continue
				}

				if !g.HasDependency(dep, tc.Job) {
					tout.SubTests[i] = format.TestOutput{
						Name:      tc.Name,
						Succeeded: false,
						Message:   fmt.Sprintf("%s does not depend on %s", tc.Job, dep),
					}
					continue
				}
			}
		} else if len(tc.DependsOn) == 0 {
			job, err := g.GetJob(tc.Job)

			if err != nil {
				tout.SubTests[i] = format.TestOutput{
					Name:      tc.Name,
					Succeeded: false,
					Message:   fmt.Sprintf("%s does not exist", tc.Job),
				}
				continue
			}

			if len(job.Needs.Needs) != 0 {
				tout.SubTests[i] = format.TestOutput{
					Name:      tc.Name,
					Succeeded: false,
					Message:   fmt.Sprintf("%s has dependencies defined", tc.Job),
				}
				continue
			}
		}

		tout.SubTests[i] = format.TestOutput{
			Name:      tc.Name,
			Succeeded: true,
		}
	}

	return tout
}
