package gitlab

import (
	"fmt"
	"reflect"

	"github.com/creasty/defaults"
)

type Need struct {
	Job       string
	Ref       string
	Project   string
	Pipeline  string
	Artifacts bool `default:"true"`
	Optional  bool `default:"false"`

	// TODO: parallel matrix
}

func (need *Need) Parse(template any) error {
	err := defaults.Set(need)
	if err != nil {
		return err
	}

	tmplVal := reflect.ValueOf(template)

	if tmplVal.Kind() == reflect.String {
		need.Job = template.(string)
		return nil
	}

	if tmplVal.Kind() == reflect.Map {
		value := reflect.ValueOf(need).Elem()
		err := parseMap(&value, template.(parsedMap))
		if err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("cannot parse need")
}
