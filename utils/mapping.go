package utils

import "reflect"

func StructToMap(item interface{}) interface{} {
	var response interface{}
	value := reflect.ValueOf(item)

	if value.Kind() == reflect.Struct {
		response = map[string]interface{}{}
		valueType := reflect.TypeOf(item)
		for i := 0; i < value.NumField(); i++ {
			v := value.Field(i)
			vt := valueType.Field(i)
			key := LowerFirstLetter(vt.Name)
			tagName := vt.Tag.Get("name")
			if tagName != "" {
				key = tagName
			}
			response.(map[string]interface{})[key] = StructToMap(v.Interface())
		}
	} else {
		response = value.Interface()
	}

	return response
}
