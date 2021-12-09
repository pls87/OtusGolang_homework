package hw09structvalidator

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID           string `json:"id" validate:"len:36"`
		Name         string
		Age          uint8    `validate:"min:18|max:50"`
		Email        string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role         UserRole `validate:"len:5|in:admin,stuff"`
		Phones       []string `validate:"len:11"`
		ChildrenAges []int    `validate:"max:17"`
		// meta         json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func lookupForError(ve ValidationError, vErrs ValidationErrors) bool {
	var equal bool
	for _, v := range vErrs {
		equal = true
		equal = equal && v.Field == ve.Field
		equal = equal && v.Step.Op == ve.Step.Op
		equal = equal && v.Step.Param == ve.Step.Param
		equal = equal && errors.Is(ve.Err, v.Err)

		if equal {
			return true
		}
	}
	return false
}

func checkErrorsMatch(t *testing.T, expected, actual ValidationErrors) {
	t.Helper()

	require.Equalf(t, len(expected), len(actual), "Different number of errors")
	checkErrorsInclusion(t, expected, actual)
}

func checkErrorsInclusion(t *testing.T, vErrs1, vErrs2 ValidationErrors) {
	t.Helper()

	for _, v := range vErrs2 {
		found := lookupForError(v, vErrs1)
		require.Truef(t, found,
			"{Field: %s, Validator: %v, Error: %v} expected but wasn't found",
			v.Field, v.Step, v.Err,
		)
	}
}

func TestValidatePositive(t *testing.T) {
	tests := []struct {
		name        string
		in          interface{}
		expectedErr ValidationErrors
	}{
		{
			name: "User:case without errors in validation rules format",
			in: User{
				ID:           "1",
				Name:         "Pavel Lysenko",
				Age:          34,
				Email:        "pavel@domain.com",
				Role:         "Developer",
				Phones:       []string{"1234567890", "23456789012"},
				ChildrenAges: []int{1, 4, 9, 21},
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "ID",
					Step: ValidationStep{
						Op:    "len",
						Param: "36",
					},
					Val: 1,
					Err: ErrValidationLen,
				}, ValidationError{
					Field: "Role",
					Step: ValidationStep{
						Op:    "len",
						Param: "5",
					},
					Val: "Developer",
					Err: ErrValidationLen,
				}, ValidationError{
					Field: "Role",
					Step: ValidationStep{
						Op:    "in",
						Param: "admin,stuff",
					},
					Val: "Developer",
					Err: ErrValidationIn,
				}, ValidationError{
					Field: "Phones",
					Step: ValidationStep{
						Op:    "len",
						Param: "11",
					},
					Val: "1234567890",
					Err: ErrValidationLen,
				}, ValidationError{
					Field: "ChildrenAges",
					Step: ValidationStep{
						Op:    "max",
						Param: "17",
					},
					Val: 21,
					Err: ErrValidationMax,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt

			validationErrors := Validate(tt.in)
			checkErrorsMatch(t, tt.expectedErr, validationErrors)
		})
	}
}
