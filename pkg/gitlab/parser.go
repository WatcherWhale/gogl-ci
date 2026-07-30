package gitlab

import (
	"reflect"
	"strings"

	"github.com/rs/zerolog/log"
)

func getFieldKeys(structVal reflect.Type) map[string]string {
	valMap := make(map[string]string)

	for i := 0; i < structVal.NumField(); i++ {
		val := structVal.Field(i)

		if instruction, ok := val.Tag.Lookup("parser"); ok {
			if instruction == "ignore" {
				continue
			}
		}

		if key, ok := val.Tag.Lookup("gitlabci"); ok {
			valMap[key] = val.Name
		} else {
			valMap[strings.ToLower(val.Name)] = val.Name
		}
	}

	return valMap
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
		err := parseSlice(field, key, value)
		if err != nil {
			return err
		}
	}

	return nil
}

func parseMap(structPtr *reflect.Value, template map[string]any) error {
	keyMap := getFieldKeys(structPtr.Type())

	for yamlKey, value := range template {
		key := keyMap[yamlKey]
		if key == "" {
			log.Trace().Msgf("parseMap %s: unknown keyword %q (skipped)", structPtr.Type().Name(), yamlKey)
			continue
		}
		field := structPtr.FieldByName(key)

		err := parseField(&field, key, value)
		if err != nil {
			return err
		}
	}
	return nil
}

func parseSlice(field *reflect.Value, key, value any) error {
	if value == nil {
		log.Debug().Msgf("parseSlice: nil value for field %q, treating as empty (check your YAML for a key with no value)", key)
		return nil
	}

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
