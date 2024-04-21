package parser

import (
	"fmt"
	"reflect"
)

type Image struct {
	Name       string
	EntryPoint []string `default:"[]"`
}

func (image *Image) Parse(template any) error {
	tmplType := reflect.TypeOf(template)

	if tmplType.Kind() == reflect.String {
		image.Name = template.(string)
		return nil
	}

	if tmplType.Kind() == reflect.Map {
		tmplMap := template.(parsedMap)
		image.Name = tmplMap["name"].(string)

		eps := tmplMap["entrypoint"].([]interface{})
		image.EntryPoint = make([]string, len(eps))
		for i, ep := range eps {
			image.EntryPoint[i] = ep.(string)
		}

		return nil
	}

	return fmt.Errorf("cannot parse image field")
}
