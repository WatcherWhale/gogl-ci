package file

import (
	"fmt"
	"os"
)

func GetTemplateFile(fileName string) (map[any]any, error) {
	bytes, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("file read error: %v", err)
	}

	return parseYaml(bytes)
}
