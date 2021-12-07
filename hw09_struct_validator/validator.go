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

func parseStep(str string) (validationStep, error) {
	parts := strings.Split(str, ":")
	if (len(parts)) != 2 {
		return validationStep{}, ErrFormat
	}
	op, param := strings.Trim(parts[0], " \t"), strings.Trim(parts[1], " \t")
	op, param = strings.ToLower(op), strings.ToLower(param)

	return validationStep{
		Op:    op,
		Param: param,
	}, nil
}

func parseTag(str string) ([]validationStep, error) {
	stepsStr := strings.Split(str, "|")
	steps := make([]validationStep, 0, len(stepsStr))
	for _, v := range stepsStr {
		step, err := parseStep(v)
		if err != nil {
			return nil, err
		}
		steps = append(steps, step)
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

func validateInt(field string, value reflect.Value, steps []validationStep) (ev ValidationErrors, e error) {
	ev = make(ValidationErrors, 0, 3)
	val := value.Convert(int64Type).Int()

	for _, step := range steps {
		var valErr *ValidationError
		switch step.Op {
		case "min":
			valErr, e = validateIntMin(field, val, step)
		case "max":
			valErr, e = validateIntMax(field, val, step)
		case "in":
			valErr, e = validateIntIn(field, val, step)
		default:
			return nil, ErrUnsupportedValidationRule
		}
		if e != nil {
			return nil, e
		}
		if valErr != nil {
			ev = append(ev, *valErr)
		}
	}

	return ev, nil
}

func validateString(field string, value reflect.Value, steps []validationStep) (ev ValidationErrors, e error) {
	ev = make(ValidationErrors, 0, 3)
	val := value.Convert(strType).String()

	for _, step := range steps {
		var valErr *ValidationError
		switch step.Op {
		case "len":
			valErr, e = validateStrLen(field, val, step)
		case "regexp":
			valErr, e = validateStrRegex(field, val, step)
		case "in":
			valErr, e = validateStrIn(field, val, step)
		default:
			return nil, ErrUnsupportedValidationRule
		}
		if e != nil {
			return nil, e
		}
		if valErr != nil {
			ev = append(ev, *valErr)
		}
	}

	return ev, nil
}

func validateSlice(field string, value reflect.Value, steps []validationStep) (ev ValidationErrors, e error) {
	ev = make(ValidationErrors, 0, 3)
	e = ErrUnsupportedType
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
		validateTagStr := f.Tag.Get("validate")
		if validateTagStr == "" {
			continue
		}

		steps, err := parseTag(validateTagStr)
		if err != nil {
			return err
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
			return err
		}
		ve = append(ve, valErrs...)
	}
	return ve
}
