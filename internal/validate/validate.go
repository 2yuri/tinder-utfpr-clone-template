package validate

import "github.com/go-playground/validator/v10"

var Validate = validator.New()

func Struct(str interface{}) []string {
	var errors []string

	err := Validate.Struct(str)

	if err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			return []string{"something went wrong"}
		}

		for _, e := range errs {
			errors = append(errors, e.Error())
		}
	}

	return errors
}
