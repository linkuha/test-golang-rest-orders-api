package entity_test

import (
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/entity"
	mockService "github.com/linkuha/test-golang-rest-orders-api/internal/domain/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUserValidateOK(t *testing.T) {
	cases := []struct {
		name   string
		in     *entity.User
		expErr error
	}{
		{
			name: "ok_without_id",
			in: &entity.User{
				Username: "qwerty",
				Password: "testtest",
			},
		},
		{
			name: "ok_without_hash",
			in: &entity.User{
				ID:       "c401f9dc-1e68-4b44-82d9-3a93b09e3fe7",
				Username: "qwerty",
				Password: "testtest",
			},
		},
		{
			name: "ok_without_pass",
			in: &entity.User{
				ID:           "c401f9dc-1e68-4b44-82d9-3a93b09e3fe7",
				Username:     "qwerty",
				PasswordHash: "$2a$04$7AnWA6wlTtkjLefPxCiP0OgzAlc",
			},
		},
		{
			name: "ok_full",
			in: &entity.User{
				ID:           "c401f9dc-1e68-4b44-82d9-3a93b09e3fe7",
				Username:     "qwerty",
				Password:     "testtest",
				PasswordHash: "$2a$04$7AnWA6wlTtkjLefPxCiP0OgzAlc",
			},
		},
	}

	for _, tCase := range cases {
		err := tCase.in.Validate()
		require.NoError(t, err)
	}
}

func TestValidateError(t *testing.T) {
	cases := []struct {
		name   string
		in     *entity.User
		expErr error
	}{
		{
			name: "empty_id",
			in: &entity.User{
				Username: "qwerty",
				Password: "test",
			},
		},
		{
			name: "bad_id",
			in: &entity.User{
				ID:       "c40fdsads09e3fe7",
				Username: "qwerty",
				Password: "test",
			},
		},
		{
			name: "empty_username",
			in: &entity.User{
				ID:       "c401f9dc-sad-3a93b09e3fe7",
				Password: "test",
			},
		},
		{
			name: "long_username",
			in: &entity.User{
				ID:       "c401f9dc-sad-3a93b09e3fe7",
				Username: "jAF6L00cqFfAHcx4sJy1Md8ytKbgmaiLmzqUYqyq5FBE8e94BIwTeRdICrSsyuXugDrGGp6s4HnYFvxQzQLIGmGS8X1Cxap92hHtpb2k5NMmw7WZ1y6P5hCCbLsD6S2TWX4oFP1zGzor15svQtBY2YD0QDNC2YwOkOBwGwmv0o1y61cC42NinHSuYpGhuqmmjsHhIDoCm2M7QB2ODWi0ws0v9Faos05dsJ6uHwEsJlvoOmdUZbS2oSqIXmp9OKJBfjQraYOFM1bXzr",
				Password: "test",
			},
		},
		{
			name: "short_username",
			in: &entity.User{
				ID:       "c401f9dc-sad-3a93b09e3fe7",
				Username: "qw",
				Password: "test",
			},
		},
		{
			name: "bad_username",
			in: &entity.User{
				ID:       "c401f9dc-sad-3a93b09e3fe7",
				Username: "qwerty!~",
				Password: "test",
			},
		},
		{
			name: "empty_pass_and_hash",
			in: &entity.User{
				ID:       "c401f9dc-sad-3a93b09e3fe7",
				Username: "qwerty",
			},
		},
	}

	for _, tCase := range cases {
		err := tCase.in.Validate()
		require.Error(t, err)
	}
}

func TestBeforeCreate(t *testing.T) {
	u := entity.TestUser(t)
	assert.NoError(t, u.BeforeCreate(mockService.NewPasswordEncryptor()))
	assert.NotEmpty(t, u.PasswordHash)
}
