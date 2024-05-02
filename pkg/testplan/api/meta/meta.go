package meta

import "gopkg.in/yaml.v3"

type ApiKind struct {
	Kind    string `yaml:"kind"`
	Version string `yaml:"apiVersion"`
}

func GetAPiKind(yamlSource []byte) (ApiKind, error) {
	var kind ApiKind

	err := yaml.Unmarshal(yamlSource, &kind)
	if err != nil {
		return kind, err
	}

	return kind, nil
}

type TestPlanMeta struct {
	Name   string            `yaml:"name"`
	Labels map[string]string `yaml:"labels"`
}
