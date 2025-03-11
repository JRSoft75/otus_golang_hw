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
	result := ""
	for i, validationError := range v {
		if i > 0 {
			result += ", "
		}
		result += validationError.Err.Error()
	}
	return result
}

func Validate(v interface{}) error {
	var validationErrors ValidationErrors

	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("expected a struct, got %s", val.Kind())
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Получаем тэги валидации
		validateTag := fieldType.Tag.Get("validate")
		if validateTag == "" {
			continue // Игнорируем поля без тэгов
		}

		// Разбиваем тэги на отдельные валидаторы
		validators := strings.Split(validateTag, "|")
		for _, validator := range validators {
			err := validateField(field, validator)
			if err != nil {
				validationErrors = append(validationErrors, ValidationError{
					Field: fieldType.Name,
					Err:   err,
				})
			}
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}
	return nil
}

func validateField(field reflect.Value, validator string) error {
	switch field.Kind() {
	case reflect.String:
		return validateString(field.String(), validator)
	case reflect.Int:
		return validateInt(int(field.Int()), validator)
	case reflect.Slice:
		return validateSlice(field, validator)

	case reflect.Bool:
	case reflect.Int8:
	case reflect.Int16:
	case reflect.Int32:
	case reflect.Int64:
	case reflect.Uint:
	case reflect.Uint8:
	case reflect.Uint16:
	case reflect.Uint32:
	case reflect.Uint64:
	case reflect.Float32:
	case reflect.Float64:
	case reflect.Invalid:
	case reflect.Uintptr:
	case reflect.Complex64:
	case reflect.Complex128:
	case reflect.Array:
	case reflect.Chan:
	case reflect.Func:
	case reflect.Interface:
	case reflect.Map:
	case reflect.Ptr:
	case reflect.Struct:
	case reflect.UnsafePointer:
	default:
	}
	return nil
}

func validateString(value, validator string) error {
	switch {
	case strings.HasPrefix(validator, "len:"):
		length, _ := strconv.Atoi(validator[4:])
		if len(value) != length {
			return fmt.Errorf("длина строки должна быть %d символов", length)
		}
	case strings.HasPrefix(validator, "regexp:"):
		re := regexp.MustCompile(validator[7:])
		if !re.MatchString(value) {
			return fmt.Errorf("строка не соответствует регулярному выражению %s", validator[7:])
		}
	case strings.HasPrefix(validator, "in:"):
		options := strings.Split(validator[3:], ",")
		for _, option := range options {
			if value == option {
				return nil
			}
		}
		return fmt.Errorf("строка должна входить в множество: %s", strings.Join(options, ", "))
	}
	return nil
}

func validateInt(value int, validator string) error {
	switch {
	case strings.HasPrefix(validator, "min:"):
		minimum, _ := strconv.Atoi(validator[4:])
		if value < minimum {
			return fmt.Errorf("число не может быть меньше %d", minimum)
		}
	case strings.HasPrefix(validator, "max:"):
		maximum, _ := strconv.Atoi(validator[4:])
		if value > maximum {
			return fmt.Errorf("число не может быть больше %d", maximum)
		}
	case strings.HasPrefix(validator, "in:"):
		options := strings.Split(validator[3:], ",")
		for _, option := range options {
			num, _ := strconv.Atoi(option)
			if value == num {
				return nil
			}
		}
		return fmt.Errorf("число должно входить в множество: %s", strings.Join(options, ", "))
	}
	return nil
}

func validateSlice(field reflect.Value, validator string) error {
	var fieldValidationErrors []string
	for i := 0; i < field.Len(); i++ {
		if err := validateField(field.Index(i), validator); err != nil {
			fieldValidationErrors = append(fieldValidationErrors, fmt.Sprintf("[%d]:%s", i, err.Error()))
		}
	}
	if len(fieldValidationErrors) > 0 {
		result := ""
		for i, fieldValidationError := range fieldValidationErrors {
			if i > 0 {
				result += ", "
			}
			result += fieldValidationError
		}
		return errors.New(result)
	}

	return nil
}
