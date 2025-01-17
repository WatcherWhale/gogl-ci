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
	Metadata meta.TestPlanMeta `yaml:"metadata"`
	Tests    []TestPlanSpec    `yaml:"tests"`

	_fileLoc string `yaml:"-"`
}

type TestPlanSpec struct {
	Name     string
	Pipeline Pipeline          `yaml:"pipeline"`
	Tests    map[string]string `yaml:"tests"`
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

func (spec *TestPlanSpec) BuildVariables() map[string]string {
	variables := spec.Pipeline.Variables

	if variables == nil {
		variables = make(map[string]string)
	}

	variables["CI_DEFAULT_BRANCH"] = spec.Pipeline.DefaultBranch

	if spec.Pipeline.Branch != "" {
		variables["CI_COMMIT_BRANCH"] = spec.Pipeline.Branch
	}

	if spec.Pipeline.isMr() {
		variables["CI_PIPELINE_SOURCE"] = "merge_request_event"
	}

	if spec.Pipeline.Tag != "" {
		variables["CI_COMMIT_TAG"] = spec.Pipeline.Tag
	}

	return variables
}

func (plan *TestPlan) Validate(pipeline *gitlab.Pipeline) format.TestOutput {
	tout := format.TestOutput{
		Name:      plan.Metadata.Name,
		SubTests:  make([]format.TestOutput, 0),
		Succeeded: true,
	}

	for _, spec := range plan.Tests {
		outp := spec.validate(plan._fileLoc, *pipeline)
		tout.SubTests = append(tout.SubTests, outp)
	}

	return tout
}

func (spec *TestPlanSpec) validate(fileLoc string, pipeline gitlab.Pipeline) format.TestOutput {
	vars := spec.BuildVariables()
	g, err := graph.NewGraph(pipeline, vars)
	if err != nil {
		return format.TestOutput{
			Name:      spec.Name,
			Succeeded: false,
			Message:   err.Error(),
		}
	}

	err = g.Validate()
	if err != nil {
		return format.TestOutput{
			Name:      spec.Name,
			Succeeded: false,
			Message:   fmt.Sprintf("error while validating job graph: %v", err),
		}
	}

	outp := format.TestOutput{
		Name:      spec.Name,
		Succeeded: true,
		SubTests:  make([]format.TestOutput, 0),
	}

	for ns, file := range spec.Tests {
		tout := format.TestOutput{
			Name:      ns,
			Succeeded: true,
			SubTests:  make([]format.TestOutput, 0),
		}

		testFuncs, err := loadTestFile(path.Join(fileLoc, file))
		if err != nil {
			tout.Succeeded = false
			tout.Message = fmt.Sprintf("error while opening %s: %v", file, err)

			outp.SubTests = append(outp.SubTests, tout)
			continue
		}

		for funcName, fn := range testFuncs {
			ok, terr := fn(pipeline, *g, vars)

			tout.SubTests = append(tout.SubTests, format.TestOutput{
				Name:      funcName,
				Succeeded: ok,
				Message:   terr,
			})
		}

		outp.SubTests = append(outp.SubTests, tout)
	}

	return outp
}
