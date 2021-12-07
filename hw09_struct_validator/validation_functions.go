package hw09structvalidator

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type (
	validateIntFunc = func(field string, val int64, step validationStep) (*ValidationError, error)
	validateStrFunc = func(field string, val string, step validationStep) (*ValidationError, error)
)

var intValidators = map[string]validateIntFunc{
	"max": validateIntMax,
	"min": validateIntMin,
	"in":  validateIntIn,
}

var stringValidators = map[string]validateStrFunc{
	"len":    validateStrLen,
	"regexp": validateStrRegexp,
	"in":     validateStrIn,
}

func validateIntMax(field string, val int64, step validationStep) (*ValidationError, error) {
	max, err := strconv.ParseInt(step.Param, 0, 64)
	if err != nil {
		return nil, err
	}
	if val > max {
		return &ValidationError{
			Field: field,
			Err:   fmt.Errorf("validation error - max %d expected but got %d", max, val),
		}, nil
	}
	return nil, nil
}

func validateIntMin(field string, val int64, step validationStep) (*ValidationError, error) {
	min, err := strconv.ParseInt(step.Param, 0, 64)
	if err != nil {
		return nil, err
	}
	if val < min {
		return &ValidationError{
			Field: field,
			Err:   fmt.Errorf("validation error - min %d expected but got %d", min, val),
		}, nil
	}
	return nil, nil
}

func validateIntIn(field string, val int64, step validationStep) (*ValidationError, error) {
	set := strings.Split(step.Param, ",")
	found := false
	for _, strI := range set {
		strI = strings.Trim(strI, " \t")
		i, err := strconv.ParseInt(strI, 0, 64)
		if err != nil {
			return nil, err
		}
		found = i == val
		if found {
			break
		}
	}
	if !found {
		return &ValidationError{
			Field: field,
			Err: fmt.Errorf(
				"validation error - %d expected to be in {%s} but actually doesn't",
				val, step.Param,
			),
		}, nil
	}
	return nil, nil
}

func validateStrLen(field string, val string, step validationStep) (*ValidationError, error) {
	l, err := strconv.ParseInt(step.Param, 0, 64)
	if err != nil {
		return nil, err
	}

	if lv := int64(len(val)); lv != l {
		return &ValidationError{
			Field: field,
			Err:   fmt.Errorf("validation error - len %d expected but got %d", l, lv),
		}, nil
	}

	return nil, nil
}

func validateStrRegexp(field string, val string, step validationStep) (*ValidationError, error) {
	re, err := regexp.Compile(step.Param)
	if err != nil {
		return nil, err
	}
	if !re.MatchString(val) {
		return &ValidationError{
			Field: field,
			Err: fmt.Errorf(
				"validation error - expected '%s' match to '%s' but actually doesn't",
				val, step.Param,
			),
		}, nil
	}
	return nil, nil
}

// nolint:unparam
func validateStrIn(field string, val string, step validationStep) (*ValidationError, error) {
	set := strings.Split(step.Param, ",")
	found := false
	for _, str := range set {
		found = str == val
		if found {
			break
		}
	}
	if !found {
		return &ValidationError{
			Field: field,
			Err: fmt.Errorf(
				"validation error - '%s' expected to be in {%s} but actually doesn't",
				val, step.Param,
			),
		}, nil
	}

	return nil, nil
}
