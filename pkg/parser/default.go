package parser

import "reflect"

func initDefault(obj *struct{}) {
	objType := reflect.TypeOf(*obj)
	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)

		if _, ok := field.Tag.Lookup("default"); ok {
			if field.Type.Kind() == reflect.String {

			}
		}
	}
}
