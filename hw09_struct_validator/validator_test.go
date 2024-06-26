package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

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

	TestStruct1 struct {
		Numbers []int    `validate:"min:5|max:50"`
		Strings []string `validate:"len:5"`
	}

	TestStruct2 struct {
		Numbers []int    `validate:"in:5,10,15,20"`
		Strings []string `validate:"in:hello,world,test"`
	}

	TestStruct3 struct {
		RegexTest string `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
	}

	TestStruct4 struct {
		Numbers []int    `validate:"inn:5,10,15,20"`
		Strings []string `validate:"lorem:ipsum,test"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name        string
		in          interface{}
		expectedErr error
	}{
		{
			"Valid User",
			User{
				ID:     "12345678-1234-1234-1234-123456789012",
				Name:   "John Doe",
				Age:    30,
				Email:  "john_doe@example.com",
				Role:   "admin",
				Phones: []string{"12345678901"},
				meta:   nil,
			},
			nil,
		},
		{
			"Invalid User ID",
			User{
				ID:     "123",
				Name:   "Vasya",
				Age:    25,
				Email:  "vasya@example.com",
				Role:   "admin",
				Phones: []string{"12345678901"},
				meta:   nil,
			},
			fmt.Errorf("ID: length of '123' is 3, expected 36\n"),
		},
		{
			"Invalid User Age",
			User{
				ID:     "12345678-1234-1234-1234-123456789012",
				Name:   "Vasya",
				Age:    17,
				Email:  "vasya@example.com",
				Role:   "admin",
				Phones: []string{"12345678901"},
				meta:   nil,
			},
			fmt.Errorf("Age: value 17 is less than min 18\n"),
		},
		{
			"Invalid User Email",
			User{
				ID:     "12345678-1234-1234-1234-123456789012",
				Name:   "Vasya",
				Age:    25,
				Email:  "invalid-email",
				Role:   "admin",
				Phones: []string{"12345678901"},
				meta:   nil,
			},
			fmt.Errorf("Email: 'invalid-email' does not match pattern ^\\w+@\\w+\\.\\w+$\n"),
		},
		{
			"Invalid User Role",
			User{
				ID:     "12345678-1234-1234-1234-123456789012",
				Name:   "Vasya",
				Age:    25,
				Email:  "vasya@example.com",
				Role:   "user",
				Phones: []string{"12345678901"},
				meta:   nil,
			},
			fmt.Errorf("Role: 'user' is not in admin,stuff\n"),
		},
		{
			"Invalid User Phone",
			User{
				ID:     "12345678-1234-1234-1234-123456789012",
				Name:   "Vasya",
				Age:    25,
				Email:  "vasya@example.com",
				Role:   "admin",
				Phones: []string{"12345"},
				meta:   nil,
			},
			fmt.Errorf("Phones[0]: length of '12345' is 5, expected 11\n"),
		},
		{
			"Empty User Phones",
			User{
				ID:     "12345678-1234-1234-1234-123456789012",
				Name:   "Vasya",
				Age:    25,
				Email:  "vasya@example.com",
				Role:   "admin",
				Phones: []string{},
				meta:   nil,
			},
			nil,
		},
		{
			"Valid App",
			App{
				Version: "v1.0",
			},
			fmt.Errorf("Version: length of 'v1.0' is 4, expected 5\n"),
		},
		{
			"Invalid App Version",
			App{
				Version: "v1",
			},
			fmt.Errorf("Version: length of 'v1' is 2, expected 5\n"),
		},
		{
			"Valid Token",
			Token{
				Header:    []byte("header"),
				Payload:   []byte("payload"),
				Signature: []byte("signature"),
			},
			nil,
		},
		{
			"Valid Response",
			Response{
				Code: 200,
				Body: "OK",
			},
			nil,
		},
		{
			"Invalid Response Code",
			Response{
				Code: 201,
				Body: "Created",
			},
			fmt.Errorf("Code: value 201 is not in 200,404,500\n"),
		},
		{
			"Empty Response Body",
			Response{
				Code: 200,
				Body: "",
			},
			nil,
		},
		{
			"Valid TestStruct1",
			TestStruct1{
				Numbers: []int{10, 20, 30, 40, 50},
				Strings: []string{"hello", "world", "test1", "goooo", "lands"},
			},
			nil,
		},
		{
			"Invalid TestStruct1 Numbers",
			TestStruct1{
				Numbers: []int{1, 2, 3, 4},
				Strings: []string{"hello", "world", "test", "go", "land"},
			},
			fmt.Errorf("Numbers[0]: value 1 is less than min 5\nNumbers[1]: value 2 is less than min 5\nNumbers[2]: value 3 is less than min 5\nNumbers[3]: value 4 is less than min 5\nStrings[2]: length of 'test' is 4, expected 5\nStrings[3]: length of 'go' is 2, expected 5\nStrings[4]: length of 'land' is 4, expected 5\n"),
		},
		{
			"Invalid TestStruct1 Strings",
			TestStruct1{
				Numbers: []int{10, 20, 30, 40, 50},
				Strings: []string{"hello", "world", "test"},
			},
			fmt.Errorf("Strings[2]: length of 'test' is 4, expected 5\n"),
		},
		{
			"Valid TestStruct2",
			TestStruct2{
				Numbers: []int{5, 10, 15, 20},
				Strings: []string{"hello", "world", "test"},
			},
			nil,
		},
		{
			"Invalid TestStruct2 Numbers",
			TestStruct2{
				Numbers: []int{1, 2, 3, 4},
				Strings: []string{"hello", "world", "test"},
			},
			fmt.Errorf("Numbers[0]: value 1 is not in 5,10,15,20\nNumbers[1]: value 2 is not in 5,10,15,20\nNumbers[2]: value 3 is not in 5,10,15,20\nNumbers[3]: value 4 is not in 5,10,15,20\n"),
		},
		{
			"Invalid TestStruct2 Strings",
			TestStruct2{
				Numbers: []int{5, 10, 15, 20},
				Strings: []string{"hi", "there", "go"},
			},
			fmt.Errorf("Strings[0]: 'hi' is not in hello,world,test\nStrings[1]: 'there' is not in hello,world,test\nStrings[2]: 'go' is not in hello,world,test\n"),
		},
		{
			"Valid TestStruct3",
			TestStruct3{
				RegexTest: "test@example.com",
			},
			nil,
		},
		{
			"Invalid TestStruct3",
			TestStruct3{
				RegexTest: "invalid-email",
			},
			fmt.Errorf("RegexTest: 'invalid-email' does not match pattern ^\\w+@\\w+\\.\\w+$\n"),
		},
		{
			"Invalid TestStruct4",
			TestStruct4{
				Numbers: []int{5, 10, 15, 20},
				Strings: []string{"hello", "world", "test"},
			},
			nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d: %s", i, tt.name), func(t *testing.T) {
			err := Validate(tt.in)

			if tt.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.True(t, areAllValidationError(err.(ValidationErrors)))
				require.Equal(t, tt.expectedErr.Error(), err.Error())
			}
		})
	}
}

func areAllValidationError(err ValidationErrors) bool {
	for _, e := range err {
		if !e.IsValidation() {
			return false
		}
	}
	return true
}
