package gitlab

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/creasty/defaults"
	"github.com/rs/zerolog/log"
	"github.com/watcherwhale/gitlabci-test/pkg/gitlab/file"
)

type Pipeline struct {
	Stages []string `default:"[\"build\", \"test\", \"deploy\"]"`

	Include []Include `default:"[]"`

	Jobs map[string]Job `default:"{}"`

	Variables map[string]string

	Default Job
}

func (pipeline *Pipeline) String() string {
	bytes, err := json.Marshal(pipeline)
	if err != nil {
		return "{}"
	}

	return string(bytes)
}

func (pipeline *Pipeline) Parse(template map[any]any, recursive bool) error {
	err := defaults.Set(pipeline)
	if err != nil {
		return err
	}

	pipeline.Default._filled = true

	keyMap := getFieldKeys(reflect.TypeOf(*pipeline))

	log.Logger.Trace().Msg("parsing includes")
	// Parse includes first, so overwriting works
	if includes, ok := template["include"]; ok {
		var slice []any
		rIncl := reflect.ValueOf(includes)
		if rIncl.Kind() == reflect.Slice {
			slice = includes.([]any)
		} else {
			slice = []any{includes}
		}

		for _, val := range slice {
			log.Logger.Trace().Msgf("parsing include %v", val)
			var newInclude Include
			err := newInclude.Parse(val)
			if err != nil {
				return err
			}

			pipeline.Include = append(pipeline.Include, newInclude)
			i := len(pipeline.Include) - 1

			templates, err := pipeline.Include[i].GetTemplate()
			if err != nil {
				return err
			}
			for _, template := range templates {
				err = pipeline.Parse(template, true)
				if err != nil {
					return err
				}
			}
		}
	}

	structPtr := reflect.ValueOf(pipeline).Elem()
	for yamlKey, value := range template {
		if yamlKey.(string) == "include" {
			continue
		}

		key, ok := keyMap[yamlKey.(string)]
		log.Logger.Trace().Msgf("parsing %s", yamlKey.(string))

		// If key is not known assume a job is found
		if !ok {
			var job Job
			err := job.Parse(yamlKey.(string), value.(map[any]any))
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

	for key, job := range pipeline.Jobs {
		job.Fill(pipeline)
		pipeline.Jobs[key] = job
	}

	return nil
}

func Parse(fileName string) (*Pipeline, error) {
	var pipeline Pipeline

	template, err := file.GetTemplateFile(fileName)
	if err != nil {
		return nil, err
	}

	err = pipeline.Parse(template, false)
	if err != nil {
		return nil, err
	}

	return &pipeline, nil
}
