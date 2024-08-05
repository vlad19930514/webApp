package util

import (
	"errors"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestGetValidationErrors(t *testing.T) {
	validate := validator.New()
	type TestStruct struct {
		FirstName string `validate:"required,alpha"`
		Email     string `validate:"required,email"`
		Age       int    `validate:"required,min=1,max=130"`
	}

	testCases := []struct {
		name     string
		input    TestStruct
		expected []ErrorMsg
	}{
		{
			name:     "valid input",
			input:    TestStruct{FirstName: "John", Email: "john.doe@example.com", Age: 30},
			expected: nil,
		},
		{
			name:  "missing required fields",
			input: TestStruct{FirstName: "", Email: "", Age: 0},
			expected: []ErrorMsg{
				{Field: "FirstName", Message: "Это обязательное поле"},
				{Field: "Email", Message: "Это обязательное поле"},
				{Field: "Age", Message: "Это обязательное поле"},
			},
		},
		{
			name:  "invalid email and age",
			input: TestStruct{FirstName: "John", Email: "invalid-email", Age: -1},
			expected: []ErrorMsg{
				{Field: "Email", Message: "Это не email - invalid-email?"},
				{Field: "Age", Message: "Маловато будет - -1"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validate.Struct(tc.input)
			customErrors, errorsExist := GetValidationErrors(&err)

			if tc.expected == nil {
				assert.False(t, errorsExist)
				assert.Nil(t, customErrors)
			} else {
				assert.True(t, errorsExist)
				assert.ElementsMatch(t, tc.expected, customErrors)
			}
		})
	}
}

func TestGetErrorMsg(t *testing.T) {
	validate := validator.New()
	type TestStruct struct {
		FirstName string `validate:"required,alpha"`
		Email     string `validate:"required,email"`
		Age       int    `validate:"required,min=1,max=130"`
		UUID      string `validate:"required,uuid"`
	}

	testCases := []struct {
		name     string
		input    TestStruct
		expected []ErrorMsg
	}{
		{
			name:  "missing required fields",
			input: TestStruct{FirstName: "", Email: "", Age: 0, UUID: ""},
			expected: []ErrorMsg{
				{Field: "FirstName", Message: "Это обязательное поле"},
				{Field: "Email", Message: "Это обязательное поле"},
				{Field: "Age", Message: "Это обязательное поле"},
				{Field: "UUID", Message: "Это обязательное поле"},
			},
		},
		{
			name:  "invalid alpha and email",
			input: TestStruct{FirstName: "John123", Email: "invalid-email", Age: 30, UUID: "123"},
			expected: []ErrorMsg{
				{Field: "FirstName", Message: "Передаем только буквы - John123"},
				{Field: "Email", Message: "Это не email - invalid-email?"},
				{Field: "UUID", Message: "Это что - 123?"},
			},
		},
		{
			name:  "age out of range",
			input: TestStruct{FirstName: "John", Email: "john.doe@example.com", Age: 150, UUID: "123e4567-e89b-12d3-a456-426614174000"},
			expected: []ErrorMsg{
				{Field: "Age", Message: "Многовато будет - 150"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validate.Struct(tc.input)
			var ve validator.ValidationErrors
			errors.As(err, &ve)

			var customErrors []ErrorMsg
			for _, fe := range ve {
				customErrors = append(customErrors, ErrorMsg{Field: fe.Field(), Message: getErrorMsg(fe)})
			}

			assert.ElementsMatch(t, tc.expected, customErrors)
		})
	}
}
