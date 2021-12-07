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

func validateInt(field string, value reflect.Value, steps []validationStep) (ValidationErrors, error) {
	ev := make(ValidationErrors, 0, 3)
	val := value.Convert(int64Type).Int()

	for _, step := range steps {
		if validator := intValidators[step.Op]; validator != nil {
			if valErr, e := validator(field, val, step); e == nil && valErr != nil {
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

func validateString(field string, value reflect.Value, steps []validationStep) (ValidationErrors, error) {
	ev := make(ValidationErrors, 0, 3)
	val := value.Convert(strType).String()

	for _, step := range steps {
		if validator := stringValidators[step.Op]; validator != nil {
			if valErr, e := validator(field, val, step); e == nil && valErr != nil {
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

func validateSlice(field string, value reflect.Value, steps []validationStep) (ValidationErrors, error) {
	ev := make(ValidationErrors, 0, 3)
	e := ErrUnsupportedType
	for i := 0; i < value.Len(); i++ {
		var valErr ValidationErrors

		if value.Index(i).Type().ConvertibleTo(int64Type) {
			valErr, e = validateInt(field, value.Index(i), steps)
		} else if value.Index(i).Type().ConvertibleTo(strType) {
			valErr, e = validateString(field, value.Index(i), steps)
		}

		if e != nil {
			return nil, e
		}
		ev = append(ev, valErr...)
	}
	return ev, nil
}

func Validate(v interface{}) (e error) {
	val := reflect.ValueOf(v)
	t := val.Type()
	ve := make(ValidationErrors, 0, 10)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		var tagStr string
		var steps []validationStep

		if tagStr = f.Tag.Get("validate"); tagStr == "" {
			continue
		}

		if steps, e = parseTag(tagStr); e != nil {
			return e
		}

		fv := val.Field(i)
		var valErrs ValidationErrors
		switch {
		case f.Type.ConvertibleTo(int64Type):
			valErrs, e = validateInt(f.Name, fv, steps)
		case f.Type.ConvertibleTo(strType):
			valErrs, e = validateString(f.Name, fv, steps)
		case f.Type.Kind() == reflect.Slice:
			valErrs, e = validateSlice(f.Name, fv, steps)
		default:
			e = ErrUnsupportedType
		}
		if e != nil {
			return e
		}
		ve = append(ve, valErrs...)
	}
	return ve
}
