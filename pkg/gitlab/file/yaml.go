package file

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

func parseYaml(bytes []byte) (map[any]any, error) {
	pipelineMap := make(map[any]any)

	err := yaml.Unmarshal(bytes, &pipelineMap)
	if err != nil {
		return nil, fmt.Errorf("yaml error: %v", err)
	}

	return pipelineMap, nil
}
