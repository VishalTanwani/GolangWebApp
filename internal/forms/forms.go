package forms

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"net/http"
	"net/url"
	"strings"
)

//Form creates a custom structure
type Form struct {
	url.Values
	Errors errors
}

//New intialize the new structure
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

//Required will add Required error in all fields
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be null")
		}
	}
}

//Has checks if form field is in post or not empty
func (f *Form) Has(field string, r *http.Request) bool {
	x := r.Form.Get(field)
	if x == "" {
		return false
	}
	return true
}

//Valid check wheather we have errors or not in form data
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

//MinLength will check mini mum length of string
func (f *Form) MinLength(field string, length int, r *http.Request) bool {
	x := r.Form.Get(field)
	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d character long", length))
		return false
	}
	return true
}

//IsEmail to check valid email
func (f *Form) IsEmail(field string, r *http.Request) {
	if !govalidator.IsEmail(r.Form.Get(field)) {
		f.Errors.Add(field, "Invalid email address")
	}
}
