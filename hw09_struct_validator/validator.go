package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var (
	ErrFormat                    = errors.New("invalid validation tag")
	ErrUnsupportedValidationRule = errors.New("invalid validation rule")
	ErrUnsupportedType           = errors.New("invalid validation type")
)

var int64Type, strType = reflect.TypeOf((int64)(0)), reflect.TypeOf("")

type ValidationStep struct {
	Op    string
	Param string
}

func (vs ValidationStep) Int() (int64, error) {
	return strconv.ParseInt(vs.Param, 0, 64)
}

func (vs ValidationStep) SliceOf(t reflect.Type) ([]interface{}, error) {
	set := strings.Split(vs.Param, ",")
	res := make([]interface{}, 0, 5)
	var elem interface{}
	var err error
	for _, strI := range set {
		switch t {
		case int64Type:
			elem, err = strconv.ParseInt(strI, 0, 64)
		default:
			elem, err = strI, nil
		}

		if err != nil {
			return nil, err
		}
		res = append(res, elem)
	}
	return res, nil
}

func parseStep(str string) (*ValidationStep, error) {
	if parts := strings.Split(str, ":"); len(parts) == 2 {
		return &ValidationStep{
			Op:    parts[0],
			Param: parts[1],
		}, nil
	}

	return nil, ErrFormat
}

func parseTag(str string) ([]ValidationStep, error) {
	if str == "" {
		return []ValidationStep{}, nil
	}
	stepsStr := strings.Split(str, "|")
	steps := make([]ValidationStep, 0, len(stepsStr))
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
	Step  ValidationStep
	Val   interface{}
}

func (v ValidationError) Error() string {
	var message string
	switch v.Step.Op {
	case "min":
		message = fmt.Sprintf("min %s expected but got %d", v.Step.Param, v.Val)
	case "max":
		message = fmt.Sprintf("max %s expected but got %d", v.Step.Param, v.Val)
	case "len":
		message = fmt.Sprintf("length for '%s' mismatched, %s expected", v.Val, v.Step.Param)
	case "regexp":
		message = fmt.Sprintf("'%s' doesn't match to regexp '%s'", v.Val, v.Step.Param)
	case "in":
		message = fmt.Sprintf("%v expected to be in {%s} but actually doesn't", v.Val, v.Step.Param)
	default:
		message = "unknown operation"
	}

	return fmt.Sprintf("Field '%s': validator: '%s', message: %s", v.Field, v.Step.Op, message)
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	builder := strings.Builder{}
	for _, ve := range v {
		builder.WriteString(ve.Error())
		builder.WriteString("\n")
	}
	return builder.String()
}

func validateValue(field string, value reflect.Value, steps []ValidationStep) (ValidationErrors, error) {
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

func validateSlice(field string, value reflect.Value, steps []ValidationStep) (ValidationErrors, error) {
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

func callValidator(f string, v interface{}, steps []ValidationStep, vSet validatorSet) (ValidationErrors, error) {
	ev := make(ValidationErrors, 0, 3)

	for _, step := range steps {
		if validator := vSet[step.Op]; validator != nil {
			if valErr, e := validator(f, v, step); e == nil {
				if valErr != nil {
					ev = append(ev, *valErr)
				}
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

		if steps, e := parseTag(f.Tag.Get("validate")); e == nil {
			if len(steps) == 0 {
				continue
			}
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
