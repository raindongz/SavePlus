package apis

import (
	"context"
	"errors"
	"fmt"
	"testing"
)

func TestLoginEmail(t *testing.T) {
	var tests = []struct {
		name     string
		email    string
		expected error
	}{
		// the table itself
		{"Empty email", "", fmt.Errorf("email is empty")},
		{"Email with space only", "    ", fmt.Errorf("email is empty")},
		{"Valid", "hello@world.com", nil},
	}
	// The execution loop
	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := UserLoginReq{
				Password: "test",
				Email:    tt.email,
			}
			ans := CheckLoginUserInfoParams(ctx, &request)
			if !errors.Is(ans, tt.expected) && ans.Error() != tt.expected.Error() {
				t.Errorf("got %s, want %s", ans, tt.expected)
			}
		})
	}
}

func TestLoginPassword(t *testing.T) {
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
			request := UserLoginReq{
				Password: tt.password,
				Email:    "test@test.com",
			}
			ans := CheckLoginUserInfoParams(ctx, &request)
			if !errors.Is(ans, tt.expected) && ans.Error() != tt.expected.Error() {
				t.Errorf("got %s, want %s", ans, tt.expected)
			}
		})
	}
}
