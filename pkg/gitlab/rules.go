package gitlab

import (
	"fmt"
	"reflect"

	"github.com/creasty/defaults"
)

type Rule struct {
	If           string
	When         string       `default:"always"`
	AllowFailure AllowFailure `gitlabci:"allow_failure"`
	//Changes      []string
	_reference string
}

func (rule *Rule) Parse(template any) error {
	err := defaults.Set(rule)
	if err != nil {
		return err
	}

	tmplType := reflect.TypeOf(template)

	if tmplType.Kind() == reflect.String {
		rule._reference = template.(string)
		return nil
	}

	if tmplType.Kind() == reflect.Map {
		value := reflect.ValueOf(rule).Elem()
		err := parseMap(&value, template.(map[any]any))
		if err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("cannot parse rule")
}
