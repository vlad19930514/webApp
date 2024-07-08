package util

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type ErrorMsg struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func GetValidationErrors(err *error) ([]ErrorMsg, bool) {
	var ve validator.ValidationErrors
	errors.As(*err, &ve)
	fmt.Println(ve)

	if errors.As(*err, &ve) {

		out := make([]ErrorMsg, len(ve))
		for i, fe := range ve {
			out[i] = ErrorMsg{Field: fe.Field(), Message: getErrorMsg(fe)}
		}
		return out, true
	}
	return nil, false
}
func getErrorMsg(fe validator.FieldError) string {

	switch fe.Tag() {
	case "required":
		return "Это обязательное поле"
	case "alpha":
		return fmt.Sprintf("Передаем только буквы - %v", fe.Value())
	case "email":
		return fmt.Sprintf("Это не email - %v?", fe.Value())
	case "age":
		return fmt.Sprintf("Это не возраст - %v?", fe.Value())
	case "max":
		return fmt.Sprintf("Многовато будет - %v", fe.Value())
	case "min":
		return fmt.Sprintf("Маловато будет - %v", fe.Value())
	case "uuid":
		return fmt.Sprintf("Это что - %v?", fe.Value())
	}
	return "Unknown error"
}
