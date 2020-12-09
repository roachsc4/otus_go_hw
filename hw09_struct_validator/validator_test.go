package hw09_struct_validator //nolint:golint,stylecheck

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/roachsc4/otus_go_hw/hw09_struct_validator/validators"
	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
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

	Counter struct {
		Counter int `validate:"max:100"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			App{"1.0.0"},
			nil,
		},
		{
			Token{
				Header:    []byte("header"),
				Payload:   []byte("payload"),
				Signature: []byte("signature"),
			},
			nil,
		},
		{
			User{
				ID:     "544114dd-f76b-48f3-8780-d97bc35b6d25",
				Name:   "Viktor",
				Age:    20,
				Email:  "asdw@mm.ru",
				Role:   "admin",
				Phones: []string{"88005553535"},
				meta:   nil,
			},
			nil,
		},
		{
			Response{
				Code: 200,
				Body: "{}",
			},
			nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			err := Validate(tt.in)
			require.NoError(t, err)
		})
	}
}

func TestValidateWithErrors(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			App{Version: "1"},
			ValidationErrors{
				ValidationError{
					Field: "Version",
					Err: &validators.StringLengthError{
						Value:          "1",
						RequiredLength: 5,
					},
				},
			},
		},
		{
			User{
				ID:     "abcd",
				Name:   "Ivan",
				Age:    16,
				Email:  "wrong_email",
				Role:   "batman",
				Phones: []string{""},
				meta:   nil,
			},
			ValidationErrors{
				ValidationError{
					Field: "ID",
					Err: &validators.StringLengthError{
						Value:          "abcd",
						RequiredLength: 36,
					},
				},
				ValidationError{
					Field: "Age",
					Err: &validators.MinValidatorError{
						Value:    16,
						MinValue: 18,
					},
				},
				ValidationError{
					Field: "Email",
					Err: &validators.RegexpError{
						Value:               "wrong_email",
						RegexpPatternString: `^\w+@\w+\.\w+$`,
					},
				},
				ValidationError{
					Field: "Role",
					Err: &validators.StringInError{
						Value: "batman",
						AllowedValues: map[string]struct{}{
							"admin": {},
							"stuff": {},
						},
					},
				},
				ValidationError{
					Field: "Phones",
					Err: &validators.StringLengthError{
						Value:          "",
						RequiredLength: 11,
					},
				},
			},
		},
		{
			Response{
				Code: 123,
				Body: "{}",
			},
			ValidationErrors{
				ValidationError{
					Field: "Code",
					Err: &validators.IntInError{
						Value: 123,
						AllowedValues: map[int]struct{}{
							200: {},
							404: {},
							500: {},
						},
					},
				},
			},
		},
		{
			Counter{Counter: 150},
			ValidationErrors{
				ValidationError{
					Field: "Counter",
					Err: &validators.MaxValidatorError{
						Value:    150,
						MaxValue: 100,
					},
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			err := Validate(tt.in)
			require.Equal(t, err, tt.expectedErr)
		})
	}
}

func TestValidateUnsupportedValue(t *testing.T) {
	err := Validate(1)
	require.EqualError(t, err, "unsupported value to validate")
}

func TestWrongValidatorDefinition(t *testing.T) {
	type StructWithWrongValidatorDefinition struct {
		Name string `validate:"len123"`
	}
	err := Validate(StructWithWrongValidatorDefinition{"Ivan"})
	require.EqualError(t, err, "wrong validator definition")
}
