package parser

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/creasty/defaults"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Job struct {
	Name string

	Image Image `default:"{}"`

	Stage string `default:"test"`

	Script       []string
	BeforeScript []string
	AfterScript  []string

	Needs        []Need
	Dependencies []string

	Extends []string

	AllowFailure AllowFailure

	Artifacts Artifacts
	Cache     Cache

	Coverage string
}

func (job *Job) Parse(name string, template parsedMap) error {
	err := defaults.Set(job)

	if err != nil {
		return fmt.Errorf("setting defaults error: %v", err)
	}

	job.Name = name

	structPtr := reflect.ValueOf(job).Elem()
	for key, value := range template {
		field := structPtr.FieldByName(cases.Title(language.English, cases.Compact).String(key.(string)))
		err := parseField(&field, key, value)
		if err != nil {
			return fmt.Errorf("error parsing key %s: %v", key.(string), err)
		}
	}

	return nil
}

func parseField(field *reflect.Value, key, value any) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value.(string))
	case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64:
		field.SetInt(value.(int64))
	case reflect.Bool:
		field.SetBool(value.(bool))
	case reflect.Struct:
		err := parseStruct(field, key, value)
		if err != nil {
			return err
		}
	case reflect.Slice:
		err := parseSlice(field, key, value)
		if err != nil {
			return err
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
