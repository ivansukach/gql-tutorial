package validator

import "fmt"

func (v *Validator) EqualToField(field1 string, value1 interface{}, field2 string, value2 interface{}) bool {
	if _, ok := v.Errors[field1]; ok {
		return false
	}
	if value1 != value2 {
		v.Errors[field2] = fmt.Sprintf("Values of %s and %s must match", field1, field2)
		return false
	}
	return true
}
