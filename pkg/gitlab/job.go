package gitlab

import (
	"encoding/json"
	"fmt"
	"reflect"

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
}

func (job *Job) Parse(name string, template parsedMap) error {
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
	}

	return nil
}

func (job *Job) String() string {
	bytes, err := json.Marshal(job)
	if err != nil {
		return "{}"
	}

	return string(bytes)
}
