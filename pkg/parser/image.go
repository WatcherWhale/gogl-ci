package parser

import (
	"fmt"
	"reflect"
)

var _ parseable = &Image{}
var _ parseable = (*Image)(nil)

type Image struct {
	Name string
	EntryPoint []string `default:"[]"`
}

func (image *Image) Parse(template any) error {
	tmplType := reflect.TypeOf(template)

	if tmplType.Kind() == reflect.String {
		image.Name = template.(string)
		return nil
	}

	if tmplType.Kind() == reflect.Map {
		tmplMap := template.(map[string]interface{})
		image.Name = tmplMap["name"].(string)
		image.EntryPoint = tmplMap["entrypoint"].([]string)
	}

	return fmt.Errorf("cannot parse image field")
}
