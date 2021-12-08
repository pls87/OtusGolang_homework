package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var (
	ErrFormat                    = errors.New("invalid validation tag")
	ErrUnsupportedValidationRule = errors.New("invalid validation rule")
	ErrUnsupportedType           = errors.New("invalid validation type")
)

var int64Type, strType = reflect.TypeOf((int64)(0)), reflect.TypeOf("")

type validationStep struct {
	Op    string
	Param string
}

func parseStep(str string) (*validationStep, error) {
	if parts := strings.Split(str, ":"); len(parts) == 2 {
		return &validationStep{
			Op:    parts[0],
			Param: parts[1],
		}, nil
	}

	return nil, ErrFormat
}

func parseTag(str string) ([]validationStep, error) {
	stepsStr := strings.Split(str, "|")
	steps := make([]validationStep, 0, len(stepsStr))
	for _, v := range stepsStr {
		if step, err := parseStep(v); err == nil {
			steps = append(steps, *step)
			continue
		} else {
			return nil, err
		}
	}

	return steps, nil
}

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	builder := strings.Builder{}
	for _, ve := range v {
		builder.WriteString(fmt.Sprintf("Field '%s': ", ve.Field))
		builder.WriteString(ve.Err.Error())
		builder.WriteString("\n")
	}
	return builder.String()
}

func validateValue(field string, value reflect.Value, steps []validationStep) (ValidationErrors, error) {
	switch {
	case value.Type().ConvertibleTo(int64Type):
		return callValidator(field, value.Convert(int64Type).Int(), steps, intValidators)
	case value.Type().ConvertibleTo(strType):
		return callValidator(field, value.Convert(strType).String(), steps, stringValidators)
	case value.Type().Kind() == reflect.Slice:
		return validateSlice(field, value, steps)
	default:
		return nil, ErrUnsupportedType
	}
}

func validateSlice(field string, value reflect.Value, steps []validationStep) (ValidationErrors, error) {
	ev := make(ValidationErrors, 0, 3)
	for i := 0; i < value.Len(); i++ {
		if valErr, e := validateValue(field, value.Index(i), steps); e == nil {
			ev = append(ev, valErr...)
		} else {
			return nil, e
		}
	}
	return ev, nil
}

func callValidator(f string, v interface{}, steps []validationStep, vSet validatorSet) (ValidationErrors, error) {
	ev := make(ValidationErrors, 0, 3)

	for _, step := range steps {
		if validator := vSet[step.Op]; validator != nil {
			if valErr, e := validator(f, v, step); e == nil && valErr != nil {
				ev = append(ev, *valErr)
			} else {
				return nil, e
			}
		} else {
			return nil, ErrUnsupportedValidationRule
		}
	}

	return ev, nil
}

func Validate(v interface{}) error {
	val := reflect.ValueOf(v)
	t := val.Type()
	ve := make(ValidationErrors, 0, 10)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		var tagStr string
		if tagStr = f.Tag.Get("validate"); tagStr == "" {
			continue
		}

		if steps, e := parseTag(tagStr); e == nil {
			if valErr, err := validateValue(f.Name, val.Field(i), steps); e == nil {
				ve = append(ve, valErr...)
			} else {
				return err
			}
		} else {
			return e
		}
	}
	return ve
}
