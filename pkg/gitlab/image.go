package gitlab

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
		value := reflect.ValueOf(image).Elem()
		err := parseMap(&value, template.(parsedMap))
		if err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("cannot parse image")
}
