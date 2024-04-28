package gitlab

import (
	"fmt"
	"reflect"
)

type Rule struct {
	If           string
	When         string
	AllowFailure AllowFailure `gitlabci:"allow_failure"`
	//Changes      []string
	_reference string
}

func (rule *Rule) Parse(template any) error {
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
