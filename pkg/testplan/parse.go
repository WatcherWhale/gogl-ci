package testplan

import (
	"fmt"
	"os"

	"github.com/watcherwhale/gogl-ci/pkg/testplan/api/meta"
	v1alpha1 "github.com/watcherwhale/gogl-ci/pkg/testplan/api/v1alpha1"
	v1alpha2 "github.com/watcherwhale/gogl-ci/pkg/testplan/api/v1alpha2"
)

func ParseFile(file string) (meta.TestPlan, error) {
	yamlSrc, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return ParseSpec(file, yamlSrc)
}

func ParseSpec(file string, yamlSource []byte) (meta.TestPlan, error) {
	kind, err := meta.GetAPiKind(yamlSource)
	if err != nil {
		return nil, err
	}

	if kind.Kind != "TestPlan" {
		return nil, fmt.Errorf("yaml is not a TestPlan")
	}

	switch kind.Version {
	case "test.gogl.ci/v1alpha1":
		return v1alpha1.LoadPlan(yamlSource)
	case "test.gogl.ci/v1alpha2":
		return v1alpha2.LoadPlan(file, yamlSource)
	default:
		return nil, fmt.Errorf("unknown api version '%s'", kind.Version)
	}
}
