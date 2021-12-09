package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var (
	ErrValidationFormat          = errors.New("invalid validation")
	ErrValidationTag             = fmt.Errorf("%w tag", ErrValidationFormat)
	ErrUnsupportedValidationRule = fmt.Errorf("%w rule", ErrValidationFormat)
	ErrUnsupportedType           = fmt.Errorf("%w type", ErrValidationFormat)

	ErrValidationFailed = errors.New("validation failed")
	ErrValidationMin    = fmt.Errorf("%w for validator Min", ErrValidationFailed)
	ErrValidationMax    = fmt.Errorf("%w for validator Max", ErrValidationFailed)
	ErrValidationIn     = fmt.Errorf("%w for validator In", ErrValidationFailed)
	ErrValidationLen    = fmt.Errorf("%w for validator Len", ErrValidationFailed)
	ErrValidationRegexp = fmt.Errorf("%w for validator Regexp", ErrValidationFailed)
)

type ValidationStep struct {
	Op    string
	Param string
}

func (vs ValidationStep) value(t reflect.Type) (interface{}, error) {
	return converters[t](vs.Param)
}

func (vs ValidationStep) slice(t reflect.Type) ([]interface{}, error) {
	res := make([]interface{}, 0, 5)
	for _, str := range strings.Split(vs.Param, ",") {
		if elem, err := converters[t](str); err == nil {
			res = append(res, elem)
		} else {
			return nil, err
		}
	}
	return res, nil
}

func newStep(str string) *ValidationStep {
	if parts := strings.Split(str, ":"); len(parts) == 2 {
		return &ValidationStep{
			Op:    parts[0],
			Param: parts[1],
		}
	}

	return nil
}

func newTag(str string) ([]ValidationStep, error) {
	if str == "" {
		return []ValidationStep{}, nil
	}
	stepLines := strings.Split(str, "|")
	steps := make([]ValidationStep, 0, len(stepLines))
	for _, v := range stepLines {
		if s := newStep(v); s != nil {
			steps = append(steps, *s)
			continue
		}
		return nil, fmt.Errorf("'%s' is %w", str, ErrValidationTag)
	}

	return steps, nil
}

type ValidationError struct {
	Field string
	Step  ValidationStep
	Val   interface{}
	Err   error
}

func (v ValidationError) Unwrap() error {
	return v.Err
}

func (v ValidationError) Error() string {
	if v.Err == nil {
		return ""
	}

	if !errors.Is(v.Err, ErrValidationFailed) {
		return v.Err.Error()
	}

	var message string
	switch v.Step.Op {
	case "min":
		message = fmt.Sprintf("%v - min %s expected, but got %d", v.Err, v.Step.Param, v.Val)
	case "max":
		message = fmt.Sprintf("%v - max %s expected, but got %d", v.Err, v.Step.Param, v.Val)
	case "len":
		message = fmt.Sprintf("%v - length for '%s' mismatched, %s expected", v.Err, v.Val, v.Step.Param)
	case "regexp":
		message = fmt.Sprintf("%v - '%s' doesn't match to regexp '%s'", v.Err, v.Val, v.Step.Param)
	case "in":
		message = fmt.Sprintf("%v - %v expected to be in {%s}, but actually doesn't", v.Err, v.Val, v.Step.Param)
	default:
		message = fmt.Sprintf("%v - unknown operation", v.Err)
	}

	return fmt.Sprintf("Field '%s': rule: '%s', message: %s", v.Field, v.Step.Op, message)
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

func validateValue(field string, value reflect.Value, steps []ValidationStep) ValidationErrors {
	if value.Type().Kind() == reflect.Slice {
		return validateSlice(field, value, steps)
	}
	for _, t := range types {
		if value.Type().ConvertibleTo(t) {
			return validate(field, value.Convert(t).Interface(), steps, validatorSets[t])
		}
	}
	return ValidationErrors{
		ValidationError{Field: field, Err: fmt.Errorf("%w %s", ErrUnsupportedType, value.Type())},
	}
}

func validateSlice(field string, value reflect.Value, steps []ValidationStep) ValidationErrors {
	ev := make(ValidationErrors, 0, 3)
	for i := 0; i < value.Len(); i++ {
		if valErr := validateValue(field, value.Index(i), steps); valErr != nil {
			ev = append(ev, valErr...)
		}
	}
	return ev
}

func validate(f string, v interface{}, steps []ValidationStep, vSet validatorSet) ValidationErrors {
	ev := make(ValidationErrors, 0, 3)
	for _, s := range steps {
		if vf := vSet[s.Op]; vf == nil {
			ev = append(ev, ValidationError{
				Field: f, Step: s, Val: v,
				Err: fmt.Errorf("%w '%s' for field '%s'", ErrUnsupportedValidationRule, s.Op, f),
			})
			continue
		} else if valErr := vf(f, v, s); valErr != nil {
			ev = append(ev, *valErr)
		}
	}

	return ev
}

func Validate(v interface{}) ValidationErrors {
	val := reflect.ValueOf(v)
	t := val.Type()
	ve := make(ValidationErrors, 0, 10)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if steps, e := newTag(f.Tag.Get("validate")); e == nil {
			if valErr := validateValue(f.Name, val.Field(i), steps); valErr != nil {
				ve = append(ve, valErr...)
			}
		} else {
			ve = append(ve, ValidationError{Field: f.Name, Err: e})
		}
	}
	return ve
}
