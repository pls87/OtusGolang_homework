package hw09structvalidator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
)

type validateFunc = func(field string, val interface{}, step ValidationStep) *ValidationError

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

func validateIntMax(field string, val interface{}, step ValidationStep) *ValidationError {
	var e error
	max, err := strconv.ParseInt(step.Param, 0, 64)
	if err != nil {
		e = fmt.Errorf(
			"%w: field %s, param %s isn't suitable for validator %s",
			ErrValidationFormat, field, step.Param, step.Op,
		)
	} else if v := val.(int64); v > max {
		e = ErrValidationNotPassed
	}

	if e != nil {
		return &ValidationError{Field: field, Step: step, Val: val, Err: e}
	}
	return nil
}

func validateIntMin(field string, val interface{}, step ValidationStep) *ValidationError {
	min, err := strconv.ParseInt(step.Param, 0, 64)
	var e error
	if err != nil {
		e = fmt.Errorf(
			"%w: field %s, param %s isn't suitable for validator %s",
			ErrValidationFormat, field, step.Param, step.Op,
		)
	} else if v := val.(int64); v < min {
		e = ErrValidationNotPassed
	}

	if e != nil {
		return &ValidationError{Field: field, Step: step, Val: val, Err: e}
	}
	return nil
}

func validateIntIn(field string, val interface{}, step ValidationStep) *ValidationError {
	return validateIn(field, val, step, int64Type)
}

func validateStrLen(field string, val interface{}, step ValidationStep) *ValidationError {
	l, err := strconv.ParseInt(step.Param, 0, 64)
	var e error
	if err != nil {
		e = fmt.Errorf(
			"%w: field %s, param %s isn't suitable for validator %s",
			ErrValidationFormat, field, step.Param, step.Op,
		)
	} else if lv := int64(len(val.(string))); lv != l {
		e = ErrValidationNotPassed
	}

	if e != nil {
		return &ValidationError{Field: field, Step: step, Val: val, Err: e}
	}
	return nil
}

func validateStrRegexp(field string, val interface{}, step ValidationStep) *ValidationError {
	var e error
	re, err := regexp.Compile(step.Param)
	if err != nil {
		e = fmt.Errorf(
			"%w: field %s, param %s isn't suitable for validator %s",
			ErrValidationFormat, field, step.Param, step.Op,
		)
	} else if !re.MatchString(val.(string)) {
		e = ErrValidationNotPassed
	}

	if e != nil {
		return &ValidationError{Field: field, Step: step, Val: val, Err: e}
	}
	return nil
}

func validateStrIn(field string, val interface{}, step ValidationStep) *ValidationError {
	return validateIn(field, val, step, strType)
}

func validateIn(field string, val interface{}, step ValidationStep, t reflect.Type) *ValidationError {
	var e error
	set, err := step.SliceOf(t)
	if err != nil {
		e = fmt.Errorf(
			"%w: field %s, param %s isn't suitable for validator %s",
			ErrValidationFormat, field, step.Param, step.Op,
		)
	} else {
		var found bool
		for _, i := range set {
			if found = i == val; found {
				break
			}
		}
		if !found {
			e = ErrValidationNotPassed
		}
	}

	if e != nil {
		return &ValidationError{Field: field, Step: step, Val: val, Err: e}
	}
	return nil
}
