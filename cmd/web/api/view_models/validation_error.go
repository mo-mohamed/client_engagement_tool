/*
Holds ValidationError struct used to show invalid requests
*/
package viewModels

import "gopkg.in/go-playground/validator.v9"

type ValidationError struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

func convertErrors(validationError validator.ValidationErrors) []*ValidationError {
	var errors []*ValidationError
	for _, err := range validationError {
		el := ValidationError{
			Field: err.Field(),
			Tag:   err.Tag(),
			Value: err.Param(),
			Type:  err.Type().String(),
		}
		errors = append(errors, &el)
	}
	return errors
}
