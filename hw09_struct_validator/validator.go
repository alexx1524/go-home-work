package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

const validationTag = "validate"

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (errors ValidationErrors) Error() string {
	var sb strings.Builder
	for _, err := range errors {
		sb.WriteString(fmt.Sprintf("Field:%s Error:%s \n", err.Field, err.Err.Error()))
	}
	return sb.String()
}

var (
	ErrNotStruct                 = errors.New("input parameter is not struct")
	ErrUnsupportedValidationType = errors.New("unsupported validation type")
	ErrWrongValidatorFormat      = errors.New("wrong validation format")
)

func Validate(v interface{}) error {
	value := reflect.ValueOf(v)
	valueType := value.Type()

	if value.Kind() != reflect.Struct {
		return ErrNotStruct
	}

	result := make(ValidationErrors, 0)

	for i := 0; i < value.Type().NumField(); i++ {
		field := valueType.Field(i)
		validation := field.Tag.Get(validationTag)

		if validation == "" {
			continue
		}

		var err error
		switch field.Type.Kind() {
		case reflect.Int:
			err = validateInt(field.Name, int(value.Field(i).Int()), validation, &result)
		case reflect.String:
			err = validateString(field.Name, value.Field(i).String(), validation, &result)
		case reflect.Slice:
			err = validateSlice(field.Name, value.Field(i), validation, &result)
		case reflect.Struct:
			err = validateNestedStruct(value.Field(i), validation, &result)
		default:
			return ErrUnsupportedValidationType
		}

		if err != nil {
			return err
		}
	}

	if len(result) == 0 {
		return nil
	}

	return result
}

func validateNestedStruct(value reflect.Value, validationRules string, validationErrors *ValidationErrors) error {
	if validationRules == "nested" {
		err := Validate(value.Interface())
		var errs ValidationErrors
		if errors.As(err, &errs) {
			*validationErrors = append(*validationErrors, errs...)
		} else {
			return err
		}
	}
	return nil
}
