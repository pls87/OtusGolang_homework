package hw09structvalidator

import (
	"fmt"
	"testing"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID        string `json:"id" validate:"len:36"`
		Name      string
		Age       uint8    `validate:"min:18|max:50"`
		Email     string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role      UserRole `validate:"in:admin,stuff"`
		Phones    []string `validate:"len:11"`
		ChildAges []int    `validate:"max:17"`
		// meta   json.RawMessage
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

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			// Place your code here.
		},
		// ...
		// Place your code here.
	}

	/*fmt.Println(Validate(User{
		ID:        "1-2-3",
		Name:      "Pavel Lysenko",
		Age:       17,
		Phones:    []string{"1234567890", "23456789012"},
		ChildAges: []int{1, 4, 9, 21},
	}))*/

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			// Place your code here.
			_ = tt
		})
	}
}
