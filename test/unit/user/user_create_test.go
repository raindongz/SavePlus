package user

import (
	"context"
	"errors"
	"fmt"
	apis "github.com/randongz/save_plus/apis"
	"testing"
)

func TestEmail(t *testing.T) {
	var tests = []struct {
		name  string
		email string
		want  error
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
			request := apis.CreateUserReq{
				Username: "test",
				Password: "test",
				Email:    tt.email,
				Phone:    "+12222222222",
				Gender:   0,
				Avatar:   ""}
			ans := apis.CheckBasicUserInfoParams(ctx, &request)
			if !errors.Is(ans, tt.want) && ans.Error() != tt.want.Error() {
				t.Errorf("got %s, want %s", ans, tt.want)
			}
		})
	}
}
