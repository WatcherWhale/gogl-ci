package gitlab

import (
	"fmt"
	"reflect"
)

type Need struct {
	Job       string
	Artifacts bool
}

func (need *Need) Parse(template any) error {
	tmplVal := reflect.ValueOf(template)

	if tmplVal.Kind() == reflect.String {
		need.Job = template.(string)
		need.Artifacts = true
		return nil
	}

	return fmt.Errorf("cannot parse need")
}
