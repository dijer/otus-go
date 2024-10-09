package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	errs := make([]string, len(v))
	for i, err := range v {
		errs[i] = fmt.Sprintf("%s: %s", err.Field, err.Err)
	}

	return strings.Join(errs, ", ")
}

var (
	ErrNoStruct = errors.New("no struct")
	ErrNoParts  = errors.New("no name and value in tag")

	ErrValidationStringLengthNotEqual = errors.New("string length not equal")
	ErrValidationRegExpNotMatch       = errors.New("regexp not match")
	ErrValidationNotIncludesString    = errors.New("not includes string")
	ErrValidationIntNotMin            = errors.New("int not in min")
	ErrValidationIntNotMax            = errors.New("int not in max")
	ErrValidationIntNotIncludes       = errors.New("int not includes")
)

func Validate(v interface{}) error {
	rValue := reflect.ValueOf(v)
	rType := reflect.TypeOf(v)

	if rValue.Kind() != reflect.Struct {
		return ErrNoStruct
	}

	var errs ValidationErrors

	for i := 0; i < rType.NumField(); i++ {
		field := rValue.Type().Field(i)
		fValue := rValue.Field(i)
		tag := field.Tag.Get("validate")

		if tag == "" {
			continue
		}

		rules := strings.Split(tag, "|")
		for _, rule := range rules {
			err := validateStruct(fValue, rule)
			if err != nil {
				errs = append(errs, ValidationError{
					Field: field.Name,
					Err:   err,
				})
			}
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

func validateStruct(field reflect.Value, tag string) error {
	parts := strings.Split(tag, ":")
	if len(parts) != 2 {
		return ErrNoParts
	}
	ruleType, ruleVal := parts[0], parts[1]

	switch field.Kind() {
	case reflect.String:
		return validateString(field, ruleType, ruleVal)
	case reflect.Int:
		return validateInt(field, ruleType, ruleVal)
	case reflect.Slice:
		return validateSlice(field, tag)
	default:
		return nil
	}
}

func validateString(field reflect.Value, ruleType, ruleVal string) error {
	value := field.String()
	switch ruleType {
	case "len":
		expectedLen, err := strconv.Atoi(ruleVal)
		if err != nil {
			return err
		}

		if len(value) != expectedLen {
			return ErrValidationStringLengthNotEqual
		}
	case "regexp":
		regex, err := regexp.Compile(ruleVal)
		if err != nil {
			return err
		}

		if !regex.MatchString(value) {
			return ErrValidationRegExpNotMatch
		}
	case "in":
		vals := strings.Split(ruleVal, ",")
		for _, val := range vals {
			if value == val {
				return nil
			}
		}

		return ErrValidationNotIncludesString
	}

	return nil
}

func validateInt(field reflect.Value, ruleType, ruleVal string) error {
	value := int(field.Int())

	switch ruleType {
	case "min":
		minVal, err := strconv.Atoi(ruleVal)
		if err != nil {
			return err
		}

		if value < minVal {
			return ErrValidationIntNotMin
		}
	case "max":
		maxVal, err := strconv.Atoi(ruleVal)
		if err != nil {
			return err
		}

		if value > maxVal {
			return ErrValidationIntNotMax
		}
	case "in":
		nums := strings.Split(ruleVal, ",")
		for _, num := range nums {
			num, err := strconv.Atoi(num)
			if err != nil {
				return err
			}

			if num == value {
				return nil
			}
		}

		return ErrValidationIntNotIncludes
	}

	return nil
}

func validateSlice(field reflect.Value, tag string) error {
	for i := 0; i < field.Len(); i++ {
		err := validateStruct(field.Index(i), tag)
		if err != nil {
			return err
		}
	}

	return nil
}
