package hw09structvalidator

import (
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrLessThanMinimum     = errors.New("the value is less than the minimum allowed")
	ErrGreaterThanMaximum  = errors.New("the value is greater than the maximum allowed")
	ErrValueIsNotAllowed   = errors.New("the value is not allowed")
	ErrInvalidStringLength = errors.New("invalid string length")
	ErrWrongFormat         = errors.New("wrong format")
)

func parseValidators(validationRules string) (map[string]string, error) {
	validators := strings.Split(validationRules, "|")
	result := make(map[string]string)
	for _, validator := range validators {
		parts := strings.Split(validator, ":")
		if len(parts) < 2 || parts[1] == "" {
			return nil, ErrWrongValidatorFormat
		}
		result[parts[0]] = parts[1]
	}
	return result, nil
}

func validateSlice(fieldName string, value reflect.Value, rule string, errors *ValidationErrors) error {
	var err error
	switch value.Interface().(type) {
	case []int:
		for _, item := range value.Interface().([]int) {
			err = validateInt(fieldName, item, rule, errors)
		}
	case []string:
		for _, item := range value.Interface().([]string) {
			err = validateString(fieldName, item, rule, errors)
		}
	}
	return err
}

func validateString(fieldName string, value string, validationRules string, errors *ValidationErrors) error {
	validators, err := parseValidators(validationRules)
	if err != nil {
		return err
	}
	for key, val := range validators {
		switch key {
		case "len":
			length, err := strconv.Atoi(val)
			if err != nil {
				return err
			}
			if len(value) != length {
				*errors = append(*errors, ValidationError{
					Field: fieldName,
					Err:   ErrInvalidStringLength,
				})
			}
		case "regexp":
			isValid, err := regexp.Match(val, []byte(value))
			if err != nil {
				return err
			}
			if !isValid {
				*errors = append(*errors, ValidationError{
					Field: fieldName,
					Err:   ErrWrongFormat,
				})
			}

		case "in":
			allowedValues := strings.Split(val, ",")
			isAllowed := false
			for _, item := range allowedValues {
				if value == item {
					isAllowed = true
					break
				}
			}
			if !isAllowed {
				*errors = append(*errors, ValidationError{
					Field: fieldName,
					Err:   ErrValueIsNotAllowed,
				})
			}
		}
	}
	return nil
}

func validateInt(fieldName string, value int, validationRules string, errors *ValidationErrors) error {
	validators, err := parseValidators(validationRules)
	if err != nil {
		return err
	}
	for key, val := range validators {
		switch key {
		case "min":
			minValue, err := strconv.Atoi(val)
			if err != nil {
				return err
			}
			if value < minValue {
				*errors = append(*errors, ValidationError{
					Field: fieldName,
					Err:   ErrLessThanMinimum,
				})
			}
		case "max":
			maxValue, err := strconv.Atoi(val)
			if err != nil {
				return err
			}
			if value > maxValue {
				*errors = append(*errors, ValidationError{
					Field: fieldName,
					Err:   ErrGreaterThanMaximum,
				})
			}
		case "in":
			allowedValues, err := splitToIntSlice(val)
			if err != nil {
				return err
			}
			isAllowed := false
			for _, item := range allowedValues {
				if value == item {
					isAllowed = true
					break
				}
			}
			if !isAllowed {
				*errors = append(*errors, ValidationError{
					Field: fieldName,
					Err:   ErrValueIsNotAllowed,
				})
			}
		}
	}
	return nil
}

func splitToIntSlice(input string) ([]int, error) {
	values := strings.Split(input, ",")
	result := make([]int, 0, len(values))
	for _, item := range values {
		allowedValue, err := strconv.Atoi(item)
		if err != nil {
			return nil, err
		}
		result = append(result, allowedValue)
	}
	return result, nil
}
