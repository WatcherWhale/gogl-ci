package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/creasty/defaults"
	"gopkg.in/yaml.v2"
)

type Pipeline struct {
	Stages  []string       `default:"[\"build\", \"test\", \"deploy\"]"`
	Include []Include      `default:"[]"`
	Jobs    map[string]Job `default:"{}"`
}

func (pipeline *Pipeline) String() string {
	bytes, err := json.Marshal(pipeline)
	if err != nil {
		return "{}"
	}

	return string(bytes)
}

func (pipeline *Pipeline) Parse(template parsedMap) error {
	err := defaults.Set(pipeline)
	if err != nil {
		return err
	}

	if includes, ok := template["include"]; ok {
		var slice []any
		rIncl := reflect.ValueOf(includes)
		if rIncl.Kind() == reflect.Slice {
			slice = includes.([]any)
		} else {
			slice = []any{includes}
		}

		pipeline.Include = make([]Include, len(slice))
		for i, val := range slice {
			err := pipeline.Include[i].Parse(val)
			if err != nil {
				return err
			}

			template, err := pipeline.Include[i].GetTemplate()
			if err != nil {
				return err
			}
			err = pipeline.Parse(template)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func Parse(file string) (*Pipeline, error) {
	var pipeline Pipeline

	template, err := getFileContents(file)
	if err != nil {
		return nil, err
	}

	err = pipeline.Parse(template)
	if err != nil {
		return nil, err
	}

	return &pipeline, nil
}

func getFileContents(file string) (parsedMap, error) {
	bytes, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("file read error: %v", err)
	}

	pipelineMap := make(parsedMap)

	err = yaml.Unmarshal(bytes, &pipelineMap)
	if err != nil {
		return nil, fmt.Errorf("yaml error: %v", err)
	}

	return pipelineMap, nil
}
