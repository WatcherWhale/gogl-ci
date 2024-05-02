package gitlab

import (
	"fmt"
	"reflect"

	"github.com/creasty/defaults"
)

const (
	WHEN_ON_SUCCESS = "on_success"
	WHEN_ON_FAILURE = "on_failure"
	WHEN_NEVER      = "never"
	WHEN_ALWAYS     = "always"
	WHEN_MANUAL     = "manual"
	WHEN_DELAYED    = "delayed"
)

type Rule struct {
	If   string
	When string `default:"on_success"`

	Needs        Needs
	Variables    map[string]string
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
		err := parseMap(&value, template.(map[string]any))
		if err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("cannot parse rule")
}
