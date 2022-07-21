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

		err := switchValidator(field, value.Field(i), validation, &result)
		if err != nil {
			return err
		}
	}

	if len(result) == 0 {
		return nil
	}

	return result
}

func switchValidator(field reflect.StructField, value reflect.Value, rules string, result *ValidationErrors) error {
	switch field.Type.Kind() {
	case reflect.Int:
		return validateInt(field.Name, int(value.Int()), rules, result)
	case reflect.String:
		return validateString(field.Name, value.String(), rules, result)
	case reflect.Slice:
		return validateSlice(field.Name, value, rules, result)
	case reflect.Struct:
		return validateNestedStruct(value, rules, result)
	case reflect.Array:
	case reflect.Bool:
	case reflect.Chan:
	case reflect.Complex128:
	case reflect.Complex64:
	case reflect.Float32:
	case reflect.Float64:
	case reflect.Func:
	case reflect.Int16:
	case reflect.Int32:
	case reflect.Int64:
	case reflect.Int8:
	case reflect.Interface:
	case reflect.Invalid:
	case reflect.Map:
	case reflect.Pointer:
	case reflect.Uint:
	case reflect.Uint16:
	case reflect.Uint32:
	case reflect.Uint64:
	case reflect.Uint8:
	case reflect.Uintptr:
	case reflect.UnsafePointer:
	default:
		return ErrUnsupportedValidationType
	}
	return nil
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
