package gitlab

import (
	"encoding/json"
	"fmt"
	"reflect"
	"slices"

	"github.com/creasty/defaults"
)

type Job struct {
	Name string

	Image Image `default:"{}"`

	Stage string `default:"test"`

	Script       []string
	BeforeScript []string `default:"[]" gitlabci:"before_script"`
	AfterScript  []string `default:"[]" gitlabci:"after_script"`

	Needs        []Need
	Dependencies []string

	Extends []string

	AllowFailure AllowFailure `gitlabci:"allow_failure"`

	Artifacts Artifacts
	Cache     Cache

	Coverage string

	_keysWithValue []string `default:"[]" parser:"ignore"`
	_filled        bool     `default:"false" parser:"ignore"`
}

func (job *Job) Parse(name string, template map[any]any) error {
	err := defaults.Set(job)

	if err != nil {
		return fmt.Errorf("setting defaults error: %v", err)
	}

	keyMap := getFieldKeys(reflect.TypeOf(*job))

	job.Name = name

	structPtr := reflect.ValueOf(job).Elem()
	for yamlKey, value := range template {
		key, ok := keyMap[yamlKey.(string)]
		if !ok {
			return fmt.Errorf("error parsing job: unknown key %s", yamlKey.(string))
		}

		field := structPtr.FieldByName(key)
		err := parseField(&field, key, value)
		if err != nil {
			return fmt.Errorf("error parsing key %s: %v", key, err)
		}

		job._keysWithValue = append(job._keysWithValue, key)
	}

	return nil
}

func (job *Job) Fill(pipeline *Pipeline) {
	if job._filled {
		return
	}

	job.fill(pipeline.Default)

	for _, extendKey := range job.Extends {
		extendJob := pipeline.Jobs[extendKey]
		extendJob.Fill(pipeline)
		job.fill(extendJob)
	}

	job._filled = true
}

func (job *Job) fill(template Job) {
	jobVal := reflect.ValueOf(job).Elem()
	templateVal := reflect.ValueOf(template)

	for _, fieldName := range template._keysWithValue {
		if slices.Contains(job._keysWithValue, fieldName) {
			continue
		}

		jobVal.FieldByName(fieldName).Set(templateVal.FieldByName(fieldName))
	}
}

func (job *Job) String() string {
	bytes, err := json.Marshal(job)
	if err != nil {
		return "{}"
	}

	return string(bytes)
}
