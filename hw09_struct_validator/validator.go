package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrFormat                    = errors.New("invalid validation tag")
	ErrUnsupportedValidationRule = errors.New("invalid validation rule")
)

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
	return "TODO// To show validation errors here"
}

func validateInt(field string, v reflect.Value, steps []validationStep) (ev ValidationErrors, e error) {
	ev = make(ValidationErrors, 0, 3)
	val := v.Int()

	for _, step := range steps {
		switch step.Op {
		case "min":
			min, err := strconv.ParseInt(step.Param, 10, 64)
			if err != nil {
				return nil, err
			}
			if val < min {
				ev = append(ev, ValidationError{
					Field: field,
					Err:   fmt.Errorf("validation error for field %s: min %d expected but got %d", field, min, val),
				})
			}
		case "max":
			max, err := strconv.ParseInt(step.Param, 10, 64)
			if err != nil {
				return nil, err
			}
			if val > max {
				ev = append(ev, ValidationError{
					Field: field,
					Err:   fmt.Errorf("validation error for field %s: max %d expected but got %d", field, max, val),
				})
			}
		default:
			return nil, ErrUnsupportedValidationRule
		}
	}

	return ev, nil
}

func validateString(field string, v reflect.Value, steps []validationStep) (ev ValidationErrors, e error) {
	ev = make(ValidationErrors, 0, 3)
	val := v.String()

	for _, step := range steps {
		switch step.Op {
		case "len":
			l, err := strconv.ParseInt(step.Param, 10, 64)
			if err != nil {
				return nil, err
			}
			lv := int64(len(val))
			if lv != l {
				ev = append(ev, ValidationError{
					Field: field,
					Err:   fmt.Errorf("validation error for field %s: len %d expected but got %d", field, l, lv),
				})
			}
		case "regexp":
			re, err := regexp.Compile(step.Param)
			if err != nil {
				return nil, err
			}
			if !re.MatchString(val) {
				ev = append(ev, ValidationError{
					Field: field,
					Err: fmt.Errorf(
						"validation error for field %s: expected '%s' match to '%s' but actually doesn't",
						field, val, step.Param,
					),
				})
			}
		default:
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
		validateTagStr := f.Tag.Get("validate")
		if validateTagStr == "" {
			continue
		}

		steps, err := parseTag(validateTagStr)
		if err != nil {
			return err
		}
		fv := val.Field(i)

		switch f.Type.String() {
		case "int":
			valErr, e := validateInt(f.Name, fv, steps)
			if e != nil {
				return err
			}
			ve = append(ve, valErr...)
		case "string":
			valErr, e := validateString(f.Name, fv, steps)
			if e != nil {
				return err
			}
			ve = append(ve, valErr...)
		}
	}
	return ve
}
