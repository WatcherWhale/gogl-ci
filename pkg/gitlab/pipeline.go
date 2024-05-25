package gitlab

import (
	"encoding/json"
	"fmt"
	"reflect"
	"slices"

	"github.com/creasty/defaults"
	"github.com/watcherwhale/gogl-ci/pkg/gitlab/file"
)

type Pipeline struct {
	Stages []string `default:"[\"build\", \"test\", \"deploy\"]"`

	Include         []Include `default:"[]"`
	_parsedIncludes []string  `default:"[]"`

	Jobs map[string]Job `default:"{}"`

	Variables map[string]string

	WorkFlow WorkFlow

	Default Job
}

func (pipeline *Pipeline) GetJobs() map[string]Job {
	jobMap := make(map[string]Job)

	for name, job := range pipeline.Jobs {
		if name[0:1] != "." {
			jobMap[name] = job
		}
	}

	return jobMap
}

func (pipeline *Pipeline) GetActiveJobs(variables map[string]string) (map[string]Job, error) {
	jobMap := make(map[string]Job)

	for name, job := range pipeline.Jobs {
		if name[0:1] != "." {
			aj, err := job.GetActiveJob(variables)
			if err != nil {
				return nil, err
			}
			jobMap[name] = aj
		}
	}

	return jobMap, nil
}

func (pipeline *Pipeline) GetJobsByStage(stage string) map[string]Job {
	jobMap := make(map[string]Job)

	for name, job := range pipeline.Jobs {
		if name[0:1] != "." && job.Stage == stage {
			jobMap[name] = job
		}
	}

	return jobMap
}

func (pipeline *Pipeline) String() string {
	bytes, err := json.Marshal(pipeline)
	if err != nil {
		return "{}"
	}

	return string(bytes)
}

func (pipeline *Pipeline) parse(template map[any]any, recursive bool, parentInclude *Include) error {
	err := defaults.Set(pipeline)
	if err != nil {
		return err
	}

	pipeline.Default._filled = true

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

		for _, val := range slice {
			var newInclude Include

			// Parse include
			err := newInclude.Parse(val)
			if err != nil {
				return err
			}

			// Set parent
			newInclude._parent = parentInclude

			// Skip already parsed includes
			if slices.Contains(pipeline._parsedIncludes, newInclude.hash()) {
				continue
			}

			// Add include
			pipeline.Include = append(pipeline.Include, newInclude)
			pipeline._parsedIncludes = append(pipeline._parsedIncludes, newInclude.hash())

			// Get include pipeline templates
			templates, err := newInclude.GetTemplate()
			if err != nil {
				return err
			}

			// Parse include templates into pipeline
			i := len(pipeline.Include) - 1
			for _, template := range templates {
				err = pipeline.parse(template, true, &pipeline.Include[i])
				if err != nil {
					return err
				}
			}
		}
	}

	if defaultTmpl, ok := template["default"]; ok {
		pipeline.Default = Job{}
		err := pipeline.Default.Parse("default", defaultTmpl.(map[string]any))
		if err != nil {
			return err
		}
	}

	structPtr := reflect.ValueOf(pipeline).Elem()
	for yamlKey, value := range template {
		if yamlKey.(string) == "include" || yamlKey.(string) == "default" {
			continue
		}

		key, ok := keyMap[yamlKey.(string)]

		// If key is not known assume a job is found
		if !ok {
			var job Job
			err := job.Parse(yamlKey.(string), value.(map[string]any))
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
		err := job.Fill(pipeline)
		if err != nil {
			return err
		}
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

	err = pipeline.parse(template, false, nil)
	if err != nil {
		return nil, err
	}

	return &pipeline, nil
}
