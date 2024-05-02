package file

import (
	"fmt"
	"regexp"

	"gopkg.in/yaml.v3"
)

func parseYaml(bytes []byte) (map[any]any, error) {
	pipelineMap := make(map[any]any)

	replacer, err := regexp.Compile(`!reference\s+\[([^\]]+)\]`)
	if err != nil {
		return nil, fmt.Errorf("yaml error: %v", err)
	}

	bytes = replacer.ReplaceAll(bytes, []byte("\"!reference [$1]\""))

	err = yaml.Unmarshal(bytes, &pipelineMap)
	if err != nil {
		return nil, fmt.Errorf("yaml error: %v", err)
	}

	return pipelineMap, nil
}
