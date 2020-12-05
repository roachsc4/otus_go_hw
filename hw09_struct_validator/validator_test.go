package hw09_struct_validator //nolint:golint,stylecheck

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
		expectedErr string
	}{
		{
			App{Version: "1"},
			"Field: Version, error: string length is 1, but length 5 is required\n",
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
			`Field: ID, error: string length is 4, but length 36 is required
Field: Age, error: int 16 must be greater or equal than 18
Field: Email, error: string wrong_email doesn't fit pattern ^\w+@\w+\.\w+$
Field: Role, error: string batman doesn't fit allowed set map[admin:{} stuff:{}]
Field: Phones, error: string length is 0, but length 11 is required
`,
		},
		{
			Response{
				Code: 123,
				Body: "{}",
			},
			"Field: Code, error: int 123 doesn't fit allowed set map[200:{} 404:{} 500:{}]\n",
		},
		{
			Counter{Counter: 150},
			"Field: Counter, error: int 150 must be equal to or lesser than 100\n",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			err := Validate(tt.in)
			require.EqualError(t, err, tt.expectedErr)
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
