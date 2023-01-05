package entity

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Follower struct {
	UserID     string `json:"user_id" binding:"required"`
	FollowerID string `json:"follower_id" binding:"required"`
	//Status     int
}

func (u *Follower) Validate() error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.UserID, validation.Required, is.UUIDv4),
		validation.Field(&u.FollowerID, validation.Required, is.UUIDv4),
		validation.Field(&u.UserID, validation.NotIn(u.FollowerID)),
	)
}
