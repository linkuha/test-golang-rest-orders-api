package entity

import "testing"

func TestUser(t *testing.T) *User {
	t.Helper()

	return &User{
		Username: "user",
		Password: "password",
	}
}

func TestExistUser(t *testing.T) *User {
	t.Helper()

	return &User{
		ID:           "c401f9dc-1e68-4b44-82d9-3a93b09e3fe7",
		Username:     "user",
		PasswordHash: "testpassword",
	}
}

func TestExistUser2(t *testing.T) *User {
	t.Helper()

	return &User{
		ID:           "c401f9dc-1e68-4b44-82d9-3a93b09e3fe1",
		Username:     "user2",
		PasswordHash: "testpassword2",
	}
}

func TestFollower(t *testing.T) *Follower {
	t.Helper()

	return &Follower{
		UserID:     "c401f9dc-1e68-4b44-82d9-3a93b09e3fe7",
		FollowerID: "c401f9dc-1e68-4b44-82d9-3a93b09e3fe1",
	}
}
