package validator

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"net/mail"
)

const dateFormat = "2006-01-02"

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return fmt.Errorf("error while validation data | %w", err)
	}

	return nil
}

func ValidateEmail(fl validator.FieldLevel) bool {
	_, err := mail.ParseAddress(fl.Field().String())
	return err == nil
}

func NewValidator() (*CustomValidator, error) {
	validate := validator.New()
	if err := validate.RegisterValidation("email", ValidateEmail); err != nil {
		return nil, fmt.Errorf("error while register date | %w", err)
	}

	return &CustomValidator{validator: validate}, nil
}
