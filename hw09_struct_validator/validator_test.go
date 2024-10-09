package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
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
			in: User{
				ID:     "123",
				Name:   "name",
				Age:    123,
				Email:  "test",
				Role:   "test",
				Phones: []string{"123"},
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "ID",
					Err:   ErrValidationStringLengthNotEqual,
				},
				ValidationError{
					Field: "Age",
					Err:   ErrValidationIntNotMax,
				},
				ValidationError{
					Field: "Email",
					Err:   ErrValidationRegExpNotMatch,
				},
				ValidationError{
					Field: "Role",
					Err:   ErrValidationNotIncludesString,
				},
				ValidationError{
					Field: "Phones",
					Err:   ErrValidationStringLengthNotEqual,
				},
			},
		},
		{
			in: User{
				ID:     "4e33a976-b726-4377-a8f3-5ac93e190bfd",
				Name:   "name",
				Age:    44,
				Email:  "q@q.ru",
				Role:   "admin",
				Phones: []string{"89991234567"},
			},
			expectedErr: nil,
		},

		{
			in: App{
				Version: "5",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Version",
					Err:   ErrValidationStringLengthNotEqual,
				},
			},
		},
		{
			in: App{
				Version: "12345",
			},
			expectedErr: nil,
		},

		{
			in: Token{
				Header:    nil,
				Payload:   nil,
				Signature: nil,
			},
			expectedErr: nil,
		},
		{
			in: Token{
				Header:    []byte("application: text/json"),
				Payload:   []byte("status: 200"),
				Signature: []byte("signature"),
			},
			expectedErr: nil,
		},

		{
			in: Response{
				Code: 100,
				Body: "",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Code",
					Err:   ErrValidationIntNotIncludes,
				},
			},
		},
		{
			in: Response{
				Code: 200,
				Body: "",
			},
			expectedErr: nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			require.Equal(t, tt.expectedErr, err)
		})
	}
}
