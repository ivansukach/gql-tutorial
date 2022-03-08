package validator

import (
	"fmt"
	"regexp"
)

var emailRegexp = regexp.MustCompile("^[\\w.\\-]+@[\\w.\\-]+\\.[A-Za-z]{2,}$")

func (v *Validator) ValidateEmail(field, value string) bool {
	if _, ok := v.Errors[field]; ok {
		return false
	}

	if !emailRegexp.MatchString(value) {
		v.Errors[field] = fmt.Sprintf("%s is not valid email", field)
		return false
	}

	return true
}
