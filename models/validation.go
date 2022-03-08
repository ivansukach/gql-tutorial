package models

import "github.com/ivansukach/gql-tutorial/validator"

func (r RegisterInput) Validate() (bool, map[string]string) {
	v := validator.New()
	v.Required("username", r.Username)
	v.ValidateMinLength("username", r.Username, 5)

	v.Required("email", r.Email)
	v.ValidateEmail("email", r.Email)

	v.Required("password", r.Password)
	v.ValidateMinLength("password", r.Password, 8)

	v.Required("firstName", r.FirstName)
	v.ValidateMinLength("firstName", r.FirstName, 3)

	v.Required("lastName", r.LastName)
	v.ValidateMinLength("lastName", r.LastName, 3)

	v.Required("confirmPassword", r.ConfirmPassword)
	v.EqualToField("password", r.Password, "confirmPassword", r.Password)

	return v.Validate(), v.Errors
}

func (l LoginInput) Validate() (bool, map[string]string) {
	v := validator.New()

	v.Required("email", l.Email)
	v.ValidateEmail("email", l.Email)

	v.Required("password", l.Password)
	v.ValidateMinLength("password", l.Password, 8)

	return v.Validate(), v.Errors
}
