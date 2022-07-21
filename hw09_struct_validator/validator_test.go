package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
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

	Account struct {
		ID   string `json:"id" validate:"len:36"`
		User User   `validate:"nested"`
	}

	AccountWithInvalidTag struct {
		ID string `json:"id" validate:"len:"`
	}
)

func TestValidateWithoutErrors(t *testing.T) {
	tests := []struct {
		in interface{}
	}{
		{
			in: User{
				ID:     strings.Repeat("1", 36),
				Age:    30,
				Role:   "admin",
				Email:  "alexx@gmail.com",
				Phones: []string{"89011112233", "89011112299"},
				meta:   nil,
			},
		},
		{
			in: Token{
				Header:    []byte("header"),
				Payload:   []byte("payload"),
				Signature: []byte("signature"),
			},
		},
		{
			in: Account{
				ID: strings.Repeat("1", 36),
				User: User{
					ID:     strings.Repeat("1", 36),
					Age:    30,
					Role:   "admin",
					Email:  "alexx@gmail.com",
					Phones: []string{"89011112233", "89011112299"},
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)

			require.NoError(t, err)

			_ = tt
		})
	}

	t.Run("Empty struct", func(t *testing.T) {
		err := Validate(struct{}{})
		require.NoError(t, err)
	})
}

func TestValidateErrors(t *testing.T) {
	tests := []struct {
		name        string
		in          interface{}
		expectedErr error
	}{
		{
			name: "code is not allowed",
			in:   Response{Code: 201},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Code", Err: ErrValueIsNotAllowed},
			},
		},
		{
			name: "length of version is invalid",
			in:   App{Version: "1234"},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Version", Err: ErrInvalidStringLength},
			},
		},
		{
			name: "role is not allowed",
			in: User{
				ID:    strings.Repeat("1", 36),
				Age:   30,
				Role:  "user",
				Email: "alexx@gmail.com",
			},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Role", Err: ErrValueIsNotAllowed},
			},
		},
		{
			name: "age is less than minimum allowed",
			in: User{
				ID:    strings.Repeat("1", 36),
				Age:   1,
				Role:  "admin",
				Email: "alexx@gmail.com",
			},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Age", Err: ErrLessThanMinimum},
			},
		},
		{
			name: "age is greater than maximum allowed",
			in: User{
				ID:    strings.Repeat("1", 36),
				Age:   55,
				Role:  "admin",
				Email: "alexx@gmail.com",
			},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Age", Err: ErrGreaterThanMaximum},
			},
		},
		{
			name: "wrong email format",
			in: User{
				ID:    strings.Repeat("1", 36),
				Age:   30,
				Role:  "admin",
				Email: "alexx",
			},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Email", Err: ErrWrongFormat},
			},
		},
		{
			name: "incorrect phone's number length",
			in: User{
				ID:     strings.Repeat("1", 36),
				Age:    30,
				Role:   "admin",
				Email:  "alexx@gmail.com",
				Phones: []string{"8901"},
			},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Phones", Err: ErrInvalidStringLength},
			},
		},
		{
			name: "some errors",
			in: User{
				ID:     "1",
				Age:    1,
				Role:   "user",
				Email:  "alexx",
				Phones: []string{"8901"},
			},
			expectedErr: ValidationErrors{
				ValidationError{Field: "ID", Err: ErrInvalidStringLength},
				ValidationError{Field: "Age", Err: ErrLessThanMinimum},
				ValidationError{Field: "Email", Err: ErrWrongFormat},
				ValidationError{Field: "Role", Err: ErrValueIsNotAllowed},
				ValidationError{Field: "Phones", Err: ErrInvalidStringLength},
			},
		},
		{
			name: "some errors in the nested struct",
			in: Account{
				ID: strings.Repeat("1", 36),
				User: User{
					ID:    "1",
					Age:   30,
					Role:  "admin",
					Email: "alexx@gmail.com",
				},
			},
			expectedErr: ValidationErrors{
				ValidationError{Field: "ID", Err: ErrInvalidStringLength},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d: %s", i, tt.name), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)

			var validationErrors ValidationErrors

			require.Error(t, err)
			require.True(t, errors.As(err, &validationErrors))
			require.Equal(t, validationErrors, tt.expectedErr)
			require.Equal(t, tt.expectedErr.Error(), err.Error())

			_ = tt
		})
	}

	t.Run("if value is not a struct returns error", func(t *testing.T) {
		err := Validate(1)
		require.Error(t, err, ErrNotStruct)
	})

	t.Run("if validation tag format is invalid returns error", func(t *testing.T) {
		err := Validate(AccountWithInvalidTag{
			ID: "123",
		})
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrWrongValidatorFormat))
	})
}
