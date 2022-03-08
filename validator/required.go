package validator

import (
	"fmt"
	"reflect"
)

func (v *Validator) Required(field string, value interface{}) bool {
	if _, ok := v.Errors[field]; ok {
		return false
	}
	typedValue := reflect.ValueOf(value)

	if Empty(typedValue) {
		v.Errors[field] = fmt.Sprintf("%s is required", field)
		return false
	}

	return true
}

func Empty(typedValue reflect.Value) bool {
	switch typedValue.Kind() {
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map:
		return typedValue.Len() == 0
		//case reflect.Int, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64:
		//	return typedValue. != 0
	}
	return false
}
