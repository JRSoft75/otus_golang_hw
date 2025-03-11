package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
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

	Other struct {
		Array []int `validate:"min:10|max:100"`
	}
)

func TestValidateUser(t *testing.T) {
	tests := []struct {
		name        string
		in          interface{}
		expectedErr error
	}{
		{
			name: "valid user",
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "Alice",
				Age:    25,
				Email:  "alice@example.com",
				Role:   "admin",
				Phones: []string{"12345678901"},
			},
			expectedErr: nil,
		},
		{
			name: "invalid ID length",
			in: User{
				ID:     "short_id",
				Name:   "Alice",
				Age:    25,
				Email:  "alice@example.com",
				Role:   "admin",
				Phones: []string{"12345678901"},
			},
			expectedErr: ValidationErrors{
				{
					Field: "ID",
					Err:   fmt.Errorf("длина строки должна быть 36 символов"),
				},
			},
		},
		{
			name: "invalid age",
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "Alice",
				Age:    15,
				Email:  "alice@example.com",
				Role:   "admin",
				Phones: []string{"12345678901"},
			},
			expectedErr: ValidationErrors{
				{
					Field: "Age",
					Err:   fmt.Errorf("число не может быть меньше 18"),
				},
			},
		},
		{
			name: "invalid email format",
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "Alice",
				Age:    25,
				Email:  "invalid_email",
				Role:   "admin",
				Phones: []string{"12345678901"},
			},
			expectedErr: ValidationErrors{
				{
					Field: "Email",
					Err:   fmt.Errorf("строка не соответствует регулярному выражению ^\\w+@\\w+\\.\\w+$"),
				},
			},
		},
		{
			name: "invalid role",
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "Alice",
				Age:    25,
				Email:  "alice@example.com",
				Role:   "guest",
				Phones: []string{"12345678901"},
			},
			expectedErr: ValidationErrors{
				{
					Field: "Role",
					Err:   fmt.Errorf("строка должна входить в множество: admin, stuff"),
				},
			},
		},
		{
			name: "invalid phone length",
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "Alice",
				Age:    25,
				Email:  "alice@example.com",
				Role:   "admin",
				Phones: []string{"1234567890"}, // 10 digits instead of 11
			},
			expectedErr: ValidationErrors{
				{
					Field: "Phones",
					Err:   fmt.Errorf("[0]:длина строки должна быть 11 символов"),
				},
			},
		},
	}

	RunTests(t, tests)
}

func TestValidateApp(t *testing.T) {
	tests := []struct {
		name        string
		in          interface{}
		expectedErr error
	}{
		{
			name: "valid app",
			in: App{
				Version: "1.0.0",
			},
			expectedErr: nil,
		},
		{
			name: "invalid app",
			in: App{
				Version: "1.0.00",
			},
			expectedErr: ValidationErrors{
				{
					Field: "Version",
					Err:   fmt.Errorf("длина строки должна быть 5 символов"),
				},
			},
		},
	}

	RunTests(t, tests)
}

func TestValidateToken(t *testing.T) {
	tests := []struct {
		name        string
		in          interface{}
		expectedErr error
	}{
		{
			name: "valid token",
			in: Token{
				Header:    []byte("header"),
				Payload:   []byte("payload"),
				Signature: []byte("signature"),
			},
			expectedErr: nil, // Ожидаем, что валидация пройдет без ошибок
		},
		{
			name: "empty token",
			in: Token{
				Header:    []byte{},
				Payload:   []byte{},
				Signature: []byte{},
			},
			expectedErr: nil, // Ожидаем, что валидация пройдет без ошибок
		},
	}

	RunTests(t, tests)
}

func TestValidateResponse(t *testing.T) {
	tests := []struct {
		name        string
		in          interface{}
		expectedErr error
	}{
		{
			name: "valid response",
			in: Response{
				Code: 200,
				Body: "Success",
			},
			expectedErr: nil,
		},
		{
			name: "invalid response code",
			in: Response{
				Code: 300,
				Body: "Error",
			},
			expectedErr: ValidationErrors{
				{
					Field: "Code",
					Err:   fmt.Errorf("число должно входить в множество: 200, 404, 500"),
				},
			},
		},
	}

	RunTests(t, tests)
}

func TestValidateOther(t *testing.T) {
	tests := []struct {
		name        string
		in          interface{}
		expectedErr error
	}{
		{
			name: "valid other",
			in: Other{
				Array: []int{10, 20, 30, 40, 50},
			},
			expectedErr: nil, // Ожидаем, что валидация пройдет без ошибок
		},
		{
			name: "invalid other - more 100",
			in: Other{
				Array: []int{10, 20, 30, 40, 500},
			},
			expectedErr: ValidationErrors{
				{
					Field: "Array",
					Err:   fmt.Errorf("[4]:число не может быть больше 100"),
				},
			},
		},
		{
			name: "invalid other - less 10",
			in: Other{
				Array: []int{1, 20, 30, 40, 50},
			},
			expectedErr: ValidationErrors{
				{
					Field: "Array",
					Err:   fmt.Errorf("[0]:число не может быть меньше 10"),
				},
			},
		},
		{
			name: "invalid other - less 10 and more 100",
			in: Other{
				Array: []int{1, 20, 30, 40, 500},
			},
			expectedErr: ValidationErrors{
				{
					Field: "Array",
					Err:   fmt.Errorf("[0]:число не может быть меньше 10"),
				},
				{
					Field: "Array",
					Err:   fmt.Errorf("[4]:число не может быть больше 100"),
				},
			},
		},
	}

	RunTests(t, tests)
}

func RunTests(t *testing.T, tests []struct {
	name        string
	in          interface{}
	expectedErr error
},
) {
	t.Helper() // Указываем, что это вспомогательная функция
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d: %s", i, tt.name), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			if err != nil && tt.expectedErr == nil {
				t.Errorf("expected: no error, got: %v", err)
			}
			if err == nil && tt.expectedErr != nil {
				t.Errorf("expected error: %v, got: none", tt.expectedErr)
			}
			if err != nil && tt.expectedErr != nil {
				if !reflect.DeepEqual(err, tt.expectedErr) {
					t.Errorf("expected error: %v, got: %v", tt.expectedErr, err)
				}
			}
			_ = tt
		})
	}
}
