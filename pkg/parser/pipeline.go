package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/creasty/defaults"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v2"
)

type Pipeline struct {
	Stages       []string       `default:"[\"build\", \"test\", \"deploy\"]"`
	Include      []Include      `default:"[]"`
	Jobs         map[string]Job `default:"{}"`
	Default      Job
	Artifacts    Artifacts
	Cache        Cache
	AllowFailure AllowFailure
	BeforeScript []string `default:"[]"`
	AfterScript  []string `default:"[]"`
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
	for key, value := range template {
		if key.(string) == "include" {
			continue
		}

		field := structPtr.FieldByName(cases.Title(language.English, cases.Compact).String(key.(string)))

		if !field.IsValid() {
			var job Job
			err := job.Parse(key.(string), value.(parsedMap))
			if err != nil {
				return err
			}

			pipeline.Jobs[key.(string)] = job
			continue
		}

		err := parseField(&field, key, value)
		if err != nil {
			return fmt.Errorf("error parsing key %s: %v", key.(string), err)
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
