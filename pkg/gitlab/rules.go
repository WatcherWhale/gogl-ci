package gitlab

import (
	"fmt"
	"reflect"

	"github.com/watcherwhale/gitlabci-test/internal/util"
)

type Rule struct {
	If           string
	When         string
	AllowFailure AllowFailure `gitlabci:"allow_failure"`
	Changes      []string
}

func (rule *Rule) Parse(template any) error {
	tmplType := reflect.TypeOf(template)

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

func (rule *Rule) Compile(envMap map[string]string) Rule {
	// Copy struct
	cRule := *rule

	output, err := util.SubstituteEnv(cRule.If, envMap)

	// If an error occurred, return untemplated string
	if err == nil {
		cRule.If = output
	}

	return cRule
}
