package hw09structvalidator

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type validateFunc = func(field string, val interface{}, step validationStep) (*ValidationError, error)

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

func validateIntMax(field string, val interface{}, step validationStep) (*ValidationError, error) {
	max, err := strconv.ParseInt(step.Param, 0, 64)
	if err != nil {
		return nil, err
	}
	if v := val.(int64); v > max {
		return &ValidationError{
			Field: field,
			Err:   fmt.Errorf("validation error - max %d expected but got %d", max, v),
		}, nil
	}
	return nil, nil
}

func validateIntMin(field string, val interface{}, step validationStep) (*ValidationError, error) {
	min, err := strconv.ParseInt(step.Param, 0, 64)
	if err != nil {
		return nil, err
	}
	if v := val.(int64); v < min {
		return &ValidationError{
			Field: field,
			Err:   fmt.Errorf("validation error - min %d expected but got %d", min, v),
		}, nil
	}
	return nil, nil
}

func validateIntIn(field string, val interface{}, step validationStep) (*ValidationError, error) {
	set := strings.Split(step.Param, ",")
	v := val.(int64)
	found := false
	for _, strI := range set {
		strI = strings.Trim(strI, " \t")
		i, err := strconv.ParseInt(strI, 0, 64)
		if err != nil {
			return nil, err
		}
		found = i == v
		if found {
			break
		}
	}
	if !found {
		return &ValidationError{
			Field: field,
			Err: fmt.Errorf(
				"validation error - %d expected to be in {%s} but actually doesn't",
				v, step.Param,
			),
		}, nil
	}
	return nil, nil
}

func validateStrLen(field string, val interface{}, step validationStep) (*ValidationError, error) {
	l, err := strconv.ParseInt(step.Param, 0, 64)
	if err != nil {
		return nil, err
	}
	v := val.(string)

	if lv := int64(len(v)); lv != l {
		return &ValidationError{
			Field: field,
			Err:   fmt.Errorf("validation error - len %d expected but got %d", l, lv),
		}, nil
	}

	return nil, nil
}

func validateStrRegexp(field string, val interface{}, step validationStep) (*ValidationError, error) {
	re, err := regexp.Compile(step.Param)
	if err != nil {
		return nil, err
	}
	v := val.(string)
	if !re.MatchString(v) {
		return &ValidationError{
			Field: field,
			Err: fmt.Errorf(
				"validation error - expected '%s' match to '%s' but actually doesn't",
				v, step.Param,
			),
		}, nil
	}
	return nil, nil
}

// nolint:unparam
func validateStrIn(field string, val interface{}, step validationStep) (*ValidationError, error) {
	set := strings.Split(step.Param, ",")
	v := val.(string)
	found := false
	for _, str := range set {
		found = str == v
		if found {
			break
		}
	}
	if !found {
		return &ValidationError{
			Field: field,
			Err: fmt.Errorf(
				"validation error - '%s' expected to be in {%s} but actually doesn't",
				v, step.Param,
			),
		}, nil
	}

	return nil, nil
}
