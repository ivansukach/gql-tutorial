package validator

import "fmt"

func (v *Validator) ValidateMinLength(field, value string, minLength int) bool {
	if _, ok := v.Errors[field]; ok {
		return false
	}
	if len(value) < minLength {
		v.Errors[field] = fmt.Sprintf("%s must consist of at least %d characters", field, minLength)
		return false
	}
	return true
}
