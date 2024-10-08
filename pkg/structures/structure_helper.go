package structHelper

import (
	"reflect"
	"strings"
)

func GetFieldsAsMapStructureTags(str interface{}) []string {
	val := reflect.ValueOf(str).Elem()
	t := val.Type()

	result := make([]string, t.NumField())

	for i := 0; i < t.NumField(); i++ {
		result[i] = t.Field(i).Tag.Get("mapstructure")
	}

	return result
}

func GetFieldsAsJsonTags(str interface{}) []string {
	val := reflect.ValueOf(str).Elem()
	t := val.Type()

	result := make([]string, t.NumField())

	for i := 0; i < t.NumField(); i++ {
		result[i] = t.Field(i).Tag.Get("json")
	}

	return result
}

func GetFieldsAsUpperSnake(str interface{}) []string {
	fields := GetFieldsAsJsonTags(str)

	result := make([]string, len(fields))

	for i, v := range fields {
		result[i] = strings.ToUpper(v)
	}

	return result
}

// GetStructName - Возвращает название структуры
func GetStructName(structure interface{}) string {
	if t := reflect.TypeOf(structure); t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	} else {
		return t.Name()
	}
}
