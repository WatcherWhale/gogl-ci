package v1alpha2

import (
	"fmt"
	"path"

	"github.com/creasty/defaults"
	"github.com/watcherwhale/gogl-ci/pkg/format"
	"github.com/watcherwhale/gogl-ci/pkg/gitlab"
	"github.com/watcherwhale/gogl-ci/pkg/graph"
	"github.com/watcherwhale/gogl-ci/pkg/testplan/api/meta"
	"gopkg.in/yaml.v3"
)

type TestPlan struct {
	meta.ApiKind
	Metadata    meta.TestPlanMeta `yaml:"metadata"`
	Tests       []TestSpec        `yaml:"tests"`
	Validations []ValidationSpec  `yaml:"validations"`

	_fileLoc string `yaml:"-"`
}

type TestSpec struct {
	Name string
	Test string `yaml:"test"`
}

type ValidationSpec struct {
	Name          string            `yaml:"name"`
	DefaultBranch string            `yaml:"defaultBranch,omitempty"`
	Branch        string            `yaml:"branch,omitempty"`
	Tag           string            `yaml:"tag,omitempty"`
	MR            *bool             `yaml:"mr"`
	Variables     map[string]string `yaml:"variables"`
}

func LoadPlan(filePath string, yamlSource []byte) (meta.TestPlan, error) {
	var plan TestPlan
	err := defaults.Set(&plan)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlSource, &plan)
	if err != nil {
		return nil, err
	}

	plan._fileLoc = path.Dir(filePath)

	return &plan, nil
}

func (plan *TestPlan) Validate(pipeline *gitlab.Pipeline) format.TestOutput {
	tout := format.TestOutput{
		Name: plan.Metadata.Name,
		SubTests: []format.TestOutput{
			{
				Name:      "Pipeline Validations",
				Succeeded: true,
				SubTests:  make([]format.TestOutput, 0),
			},
		},
		Succeeded: true,
	}

	for _, spec := range plan.Validations {
		outp := spec.validate(*pipeline)
		tout.SubTests[0].SubTests = append(tout.SubTests[0].SubTests, outp)
	}

	for _, spec := range plan.Tests {
		outp := spec.validate(plan._fileLoc, *pipeline)
		tout.SubTests = append(tout.SubTests, outp)
	}

	return tout
}

func (spec *TestSpec) validate(fileLoc string, pipeline gitlab.Pipeline) format.TestOutput {
	tout := format.TestOutput{
		Name:      spec.Name,
		Succeeded: true,
		SubTests:  make([]format.TestOutput, 0),
	}

	testFuncs, err := loadTestFile(path.Join(fileLoc, spec.Test))
	if err != nil {
		tout.Succeeded = false
		tout.Message = fmt.Sprintf("error while opening %s: %v", spec.Test, err)
	}

	for funcName, fn := range testFuncs {
		ok, terr := fn(pipeline)

		tout.SubTests = append(tout.SubTests, format.TestOutput{
			Name:      funcName,
			Succeeded: ok,
			Message:   terr,
		})
	}

	return tout
}

func (pl *ValidationSpec) validate(pipeline gitlab.Pipeline) format.TestOutput {
	vars := pl.BuildVariables()

	g, err := graph.NewGraph(pipeline, vars)
	if err != nil {
		return format.TestOutput{
			Name:      pl.Name,
			Succeeded: false,
			Message:   err.Error(),
		}
	}

	err = g.Validate()
	if err != nil {
		return format.TestOutput{
			Name:      pl.Name,
			Succeeded: false,
			Message:   err.Error(),
		}
	}

	return format.TestOutput{
		Name:      pl.Name,
		Succeeded: true,
	}
}

func (p ValidationSpec) isMr() bool {
	return p.MR != nil && *p.MR
}

func (pipeline *ValidationSpec) BuildVariables() map[string]string {
	variables := pipeline.Variables

	if variables == nil {
		variables = make(map[string]string)
	}

	variables["CI_DEFAULT_BRANCH"] = pipeline.DefaultBranch

	if pipeline.Branch != "" {
		variables["CI_COMMIT_BRANCH"] = pipeline.Branch
	}

	if pipeline.isMr() {
		variables["CI_PIPELINE_SOURCE"] = "merge_request_event"
	}

	if pipeline.Tag != "" {
		variables["CI_COMMIT_TAG"] = pipeline.Tag
	}

	return variables
}
