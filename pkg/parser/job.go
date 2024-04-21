package parser

import (
	"fmt"
	"encoding/json"
	"reflect"

	"github.com/creasty/defaults"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Job struct {
	Name string

	Image Image `default:"{}"`

	Stage string `default:"test"`

    Script []string
    BeforeScript []string
    AfterScript []string

    Needs []Need
    Dependencies []string

    Extends []string

	AllowFailure AllowFailure

    Artifacts Artifacts
    Cache Cache

    Coverage string
}

func (job *Job) Parse(name string, template parsedMap) error {
	err := defaults.Set(job)

	if err != nil {
		return fmt.Errorf("setting defaults error: %v", err)
	}

	job.Name = name

	structPtr := reflect.ValueOf(job).Elem()
	for key, value := range template {
		field := structPtr.FieldByName(cases.Title(language.English, cases.Compact).String(key.(string)))
		err := parseField(&field, key, value)
		if err != nil {
			return fmt.Errorf("error parsing key %s: %v", key.(string), err)
		}
	}
	
	return nil
}

func parseField(field *reflect.Value, key, value any) error {
	switch field.Kind() {
		case reflect.String:
			field.SetString(value.(string))
		case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64:
			field.SetInt(value.(int64))
		case reflect.Bool:
			field.SetBool(value.(bool))
		case reflect.Struct:
			err := parseStruct(field, key, value)
			if err != nil {
				return err
			}
		case reflect.Slice:
			var valSlice []interface{}

			rVal := reflect.ValueOf(value)
			if rVal.Kind() == reflect.Slice {
				valSlice = value.([]interface{})
			} else {
				valSlice = []interface{}{
					value,
				}
			}

			ref := reflect.New(field.Type())
			ref.Elem().Set(reflect.MakeSlice(field.Type(), len(valSlice), len(valSlice)))

			for i, val := range valSlice {
				elem := ref.Elem().Index(i)
				if elem.Kind() == reflect.Struct {
					err := parseStruct(&elem, key, val)
					if err != nil {
						return err
					}
				} else {
					elem.Set(reflect.ValueOf(val))
				}
			}

			field.Set(ref.Elem().Convert(field.Type()))
	}

	return nil
}

func parseStruct(field *reflect.Value, _, value any) error {
	method := field.Addr().MethodByName("Parse")

	in := reflect.ValueOf(value)

	errs := method.Call([]reflect.Value{ in })
	errI := errs[0].Interface()
	if errI != nil {
		return errI.(error)
	}

	return nil
}

func (job *Job) String() string {
	bytes, err := json.Marshal(&job)
	if err != nil {
		return "{}"
	}

	return string(bytes)
}
