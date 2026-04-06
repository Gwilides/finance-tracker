package req

import "github.com/go-playground/validator/v10"

var validate = validator.New()

func IsValid(payload any) error {
	return validate.Struct(payload)
}
