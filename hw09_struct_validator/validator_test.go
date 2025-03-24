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
		name    string `validate:"in:myApp"`
		Cost    int    `validate:"minimum:20"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
		Footer    string `validate:""`
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
			in: User{ID: "12345678", Name: "Alex", Age: 30, Email: "alex@qwerty..com.", Role: "user"},
			expectedErr: ValidationErrors{
				{Field: "ID", Err: ErrLen},
				{Field: "Email", Err: ErrRegexp},
				{Field: "Role", Err: ErrIn},
			},
		},
		{
			in:          App{Version: "1.2.3", name: "game", Cost: 21},
			expectedErr: ValidationErrors{{Field: "Cost", Err: ErrValidator}},
		},
		{
			in:          Token{},
			expectedErr: ErrTag,
		},
		{
			in: Response{Code: 200},
		},
		{
			in:          Response{Code: 532},
			expectedErr: ValidationErrors{{Field: "Code", Err: ErrIn}},
		},
		{
			in:          "qwerty",
			expectedErr: ErrStruct,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			require.Equal(t, tt.expectedErr, err)
			_ = tt
		})
	}
}
