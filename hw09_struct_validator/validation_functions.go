package hw09structvalidator

import (
	"reflect"
	"regexp"
	"strconv"
)

type validateFunc = func(field string, val interface{}, step ValidationStep) (*ValidationError, error)

type validatorSet map[string]validateFunc

var intValidators = validatorSet{
	"max": validateIntMax,
	"min": validateIntMin,
	"in":  validateIntIn,
}

var stringValidators = validatorSet{
	"len":    validateStrLen,
	"regexp": validateStrRegexp,
	"in":     validateStrIn,
}

func validateIntMax(field string, val interface{}, step ValidationStep) (*ValidationError, error) {
	max, err := strconv.ParseInt(step.Param, 0, 64)
	if err != nil {
		return nil, err
	}
	if v := val.(int64); v > max {
		return &ValidationError{Field: field, Step: step, Val: val}, nil
	}
	return nil, nil
}

func validateIntMin(field string, val interface{}, step ValidationStep) (*ValidationError, error) {
	min, err := strconv.ParseInt(step.Param, 0, 64)
	if err != nil {
		return nil, err
	}
	if v := val.(int64); v < min {
		return &ValidationError{Field: field, Step: step, Val: val}, nil
	}
	return nil, nil
}

func validateIntIn(field string, val interface{}, step ValidationStep) (*ValidationError, error) {
	return validateIn(field, val, step, int64Type)
}

func validateStrLen(field string, val interface{}, step ValidationStep) (*ValidationError, error) {
	l, err := strconv.ParseInt(step.Param, 0, 64)
	if err != nil {
		return nil, err
	}

	if lv := int64(len(val.(string))); lv != l {
		return &ValidationError{Field: field, Step: step, Val: val}, nil
	}

	return nil, nil
}

func validateStrRegexp(field string, val interface{}, step ValidationStep) (*ValidationError, error) {
	re, err := regexp.Compile(step.Param)
	if err != nil {
		return nil, err
	}
	if !re.MatchString(val.(string)) {
		return &ValidationError{Field: field, Step: step, Val: val}, nil
	}
	return nil, nil
}

// nolint:unparam
func validateStrIn(field string, val interface{}, step ValidationStep) (*ValidationError, error) {
	return validateIn(field, val, step, strType)
}

func validateIn(field string, val interface{}, step ValidationStep, t reflect.Type) (*ValidationError, error) {
	set, err := step.SliceOf(t)
	if err != nil {
		return nil, err
	}

	var found bool
	for _, i := range set {
		if found = i == val; found {
			break
		}
	}
	if !found {
		return &ValidationError{Field: field, Step: step, Val: val}, nil
	}
	return nil, nil
}
