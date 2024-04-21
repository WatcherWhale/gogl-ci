package gitlab

import (
	"fmt"
	"reflect"
	"regexp"

	"github.com/creasty/defaults"
)

type Include struct {
	File    string
	Local   bool `default:"true"`
	Web     bool `default:"false"`
	Project string
	Rules   []Rule `default:"[]"`
}

func (include *Include) Parse(template any) error {
	err := defaults.Set(include)
	if err != nil {
		return err
	}

	rTmpl := reflect.ValueOf(template)

	switch rTmpl.Kind() {
	case reflect.String:
		include.File = template.(string)
		webMatch, err := regexp.Compile(`https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()!@:%_\+.~#?&\/\/=]*)`)
		if err != nil {
			return err
		}

		if webMatch.Match([]byte(include.File)) {
			include.Web = true
			include.Local = false
		}

	default:
		return fmt.Errorf("error parsing include: invalid type %s", rTmpl.Kind().String())
	}

	return nil
}

func (include *Include) GetTemplate() (parsedMap, error) {
	if include.Local {
		return include.getLocalTemplate()
	}

	return nil, fmt.Errorf("error paring include: invalid syntax")
}

func (include *Include) getLocalTemplate() (parsedMap, error) {
	file := include.File[1:]

	return getFileContents(file)
}
