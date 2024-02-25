package apis

import (
	"context"
	"errors"
	"fmt"
	"testing"
)

func TestCreateEmail(t *testing.T) {
	var tests = []struct {
		name     string
		email    string
		expected error
	}{
		// the table itself
		{"Empty email", "", fmt.Errorf("email is empty")},
		{"Email with space only", "    ", fmt.Errorf("email is empty")},
		{"Invalid email without @", "helloworld.com", fmt.Errorf("email is invalid")},
		{"Invalid email without domain", "helloworld@.com", fmt.Errorf("email is invalid")},
		{"Valid", "hello@world.com", nil},
	}
	// The execution loop
	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := CreateUserReq{
				Username: "test",
				Password: "test",
				Email:    tt.email,
				Phone:    "+12222222222",
				Gender:   0,
				Avatar:   ""}
			ans := CheckBasicUserInfoParams(ctx, &request)
			if !errors.Is(ans, tt.expected) && ans.Error() != tt.expected.Error() {
				t.Errorf("got %s, want %s", ans, tt.expected)
			}
		})
	}
}

func TestCreateUsername(t *testing.T) {
	var tests = []struct {
		name     string
		username string
		expected error
	}{
		// the table itself
		{"Empty username", "", fmt.Errorf("username is empty")},
		{"Username with space only", "    ", fmt.Errorf("username is empty")},
		{"Valid", "test_username", nil},
	}
	// The execution loop
	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := CreateUserReq{
				Username: tt.username,
				Password: "test",
				Email:    "test@test.com",
				Phone:    "+12222222222",
				Gender:   0,
				Avatar:   ""}
			ans := CheckBasicUserInfoParams(ctx, &request)
			if !errors.Is(ans, tt.expected) && ans.Error() != tt.expected.Error() {
				t.Errorf("got %s, want %s", ans, tt.expected)
			}
		})
	}
}

func TestCreatePassword(t *testing.T) {
	var tests = []struct {
		name     string
		password string
		expected error
	}{
		// the table itself
		{"Empty password", "", fmt.Errorf("password is empty")},
		{"Password with space only", "    ", fmt.Errorf("password is empty")},
		{"Valid", "testtest", nil},
	}
	// The execution loop
	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := CreateUserReq{
				Username: "test",
				Password: tt.password,
				Email:    "test@test.com",
				Phone:    "+12222222222",
				Gender:   0,
				Avatar:   ""}
			ans := CheckBasicUserInfoParams(ctx, &request)
			if !errors.Is(ans, tt.expected) && ans.Error() != tt.expected.Error() {
				t.Errorf("got %s, want %s", ans, tt.expected)
			}
		})
	}
}
