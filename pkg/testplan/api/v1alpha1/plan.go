package api

import (
	"github.com/creasty/defaults"
	"github.com/watcherwhale/gogl-ci/pkg/gitlab"
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
	DependsOn []string `yaml:"dependsOn" default:"nil"`
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

	return &plan, nil
}

func (plan *TestPlan) Validate(pipeline *gitlab.Pipeline) (bool, error) {
	return false, nil
}
