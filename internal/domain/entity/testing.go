package entity

import "testing"

func TestUser(t *testing.T) *User {
	t.Helper()

	return &User{
		Username: "user",
		Password: "password",
	}
}
