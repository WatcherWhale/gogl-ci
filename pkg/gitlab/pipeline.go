package gitlab

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/creasty/defaults"
	"gopkg.in/yaml.v2"
)

type Pipeline struct {
	Stages []string `default:"[\"build\", \"test\", \"deploy\"]"`

	Include []Include `default:"[]"`

	Jobs map[string]Job `default:"{}"`

	Default Job

	Artifacts Artifacts
	Cache     Cache

	AllowFailure AllowFailure

	BeforeScript []string `default:"[]" gitlabci:"before_script"`
	AfterScript  []string `default:"[]" gitlabci:"after_script"`
}

func (pipeline *Pipeline) String() string {
	bytes, err := json.Marshal(pipeline)
	if err != nil {
		return "{}"
	}

	return string(bytes)
}

func (pipeline *Pipeline) Parse(template parsedMap, recursive bool) error {
	err := defaults.Set(pipeline)
	if err != nil {
		return err
	}

	keyMap := getFieldKeys(reflect.TypeOf(*pipeline))

	// Parse includes first, so overwriting works
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
			err = pipeline.Parse(template, true)
			if err != nil {
				return err
			}
		}
	}

	structPtr := reflect.ValueOf(pipeline).Elem()
	for yamlKey, value := range template {
		if yamlKey.(string) == "include" {
			continue
		}

		key, ok := keyMap[yamlKey.(string)]

		// If key is not known assume a job is found
		if !ok {
			var job Job
			err := job.Parse(yamlKey.(string), value.(parsedMap))
			if err != nil {
				return err
			}

			pipeline.Jobs[yamlKey.(string)] = job
			continue
		}

		// Parse field
		field := structPtr.FieldByName(key)
		err := parseField(&field, key, value)
		if err != nil {
			return fmt.Errorf("error parsing key %s: %v", key, err)
		}
	}

	// Append .pre and .post stages
	if !recursive {
		pipeline.Stages = append(append([]string{".pre"}, pipeline.Stages...), ".post")
	}

	return nil
}

func Parse(file string) (*Pipeline, error) {
	var pipeline Pipeline

	template, err := getFileContents(file)
	if err != nil {
		return nil, err
	}

	err = pipeline.Parse(template, false)
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
