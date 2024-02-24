package test

import (
	"context"
	"errors"
	"fmt"
	apis "github.com/randongz/save_plus/apis"
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
			request := apis.UserLoginReq{
				Password: "test",
				Email:    tt.email,
			}
			ans := apis.CheckLoginUserInfoParams(ctx, &request)
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
			request := apis.UserLoginReq{
				Password: tt.password,
				Email:    "test@test.com",
			}
			ans := apis.CheckLoginUserInfoParams(ctx, &request)
			if !errors.Is(ans, tt.expected) && ans.Error() != tt.expected.Error() {
				t.Errorf("got %s, want %s", ans, tt.expected)
			}
		})
	}
}
