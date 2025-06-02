package handlers

import (
	"github.com/go-playground/validator/v10"
)

type RequestValidator struct {
	Validator *validator.Validate
}

func NewRequestValidator() *RequestValidator {
	v := validator.New()
	return &RequestValidator{Validator: v}
}

func (v *RequestValidator) Validate(i interface{}) error {
	return v.Validator.Struct(i)
}
