package requests

import (
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