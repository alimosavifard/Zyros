package requests

import (
	"github.com/go-playground/validator/v10"
	"errors"
)

type Validatable interface {
	Validate() error
}

var validate = validator.New()

func ValidateStruct(s interface{}) error {
	if err := validate.Struct(s); err != nil {
		return err
	}
	return nil
}