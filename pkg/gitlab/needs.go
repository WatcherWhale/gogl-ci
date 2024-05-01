package gitlab

import (
	"fmt"
	"reflect"

	"github.com/creasty/defaults"
)

type Needs struct {
	// A list of needs
	Needs []Need

	// Is true when relying on stages to define dependencies
	NoNeeds bool
}

func (needs *Needs) Parse(template any) error {
	structPtr := reflect.ValueOf(needs).Elem()
	needsPtr := structPtr.FieldByName("Needs")
	err := parseField(&needsPtr, "needs", template)
	if err != nil {
		return err
	}

	needs.NoNeeds = false

	return nil
}

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
		err := parseMap(&value, template.(map[any]any))
		if err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("cannot parse need")
}
