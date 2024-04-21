package parser

import (
	"reflect"
)

type parsedMap map[any]any

func parseSlice(field *reflect.Value, key, value any) error {
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

	return nil
}

func parseStruct(field *reflect.Value, _, value any) error {
	method := field.Addr().MethodByName("Parse")

	in := reflect.ValueOf(value)

	errs := method.Call([]reflect.Value{in})
	errI := errs[0].Interface()
	if errI != nil {
		return errI.(error)
	}

	return nil
}

